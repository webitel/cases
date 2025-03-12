package postgres

import (
	"errors"
	"fmt"
	"github.com/webitel/cases/internal/store/util"
	"github.com/webitel/cases/model/options"
	"github.com/webitel/cases/model/options/defaults"
	"strconv"

	sq "github.com/Masterminds/squirrel"
	"github.com/webitel/cases/api/cases"
	dberr "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/internal/store/postgres/scanner"
)

type CaseFileStore struct {
	storage   *Store
	mainTable string
}

const (
	// Alias for the storage.files table
	fileAlias               = "cf"
	channel                 = "case"
	fileDefaultSort         = "uploaded_at"
	caseFileAuthorAlias     = "au"
	caseFileNotRemovedAlias = "ra"
	caseFileCreatedByAlias  = "cb"
)

// List implements store.CaseFileStore for listing case files.
func (c *CaseFileStore) List(rpc options.SearchOptions) (*cases.CaseFileList, error) {
	// Connect to the database
	d, dbErr := c.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("store.case_file.list.database_connection_error", dbErr)
	}

	// Build the query and plan builder using BuildListCaseFilesSqlizer
	queryBuilder, plan, err := c.BuildListCaseFilesSqlizer(rpc)
	if err != nil {
		return nil, dberr.NewDBInternalError("store.case_file.list.query_build_error", err)
	}

	// Convert the query to SQL
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, dberr.NewDBInternalError("store.case_file.list.query_to_sql_error", err)
	}

	// Execute the query
	rows, err := d.Query(rpc, query, args...)
	if err != nil {
		return nil, dberr.NewDBInternalError("store.case_file.list.execution_error", err)
	}
	defer rows.Close()

	var fileList []*cases.File
	lCount := 0
	next := false
	fetchAll := rpc.GetSize() == -1

	for rows.Next() {
		if !fetchAll && lCount >= rpc.GetSize() {
			next = true
			break
		}

		// Create a new file object
		file := &cases.File{}
		// Build the scan plan using the planBuilder function

		// Scan row into the file fields using the plan directly
		scanArgs := make([]any, len(plan))
		for i, scanFunc := range plan {
			scanArgs[i] = scanFunc(file)
		}

		// Scan row into the file fields using the plan
		if err := rows.Scan(scanArgs...); err != nil {
			return nil, dberr.NewDBInternalError("store.case_file.list.row_scan_error", err)
		}

		fileList = append(fileList, file)
		lCount++
	}

	return &cases.CaseFileList{
		Page:  int64(rpc.GetPage()),
		Next:  next,
		Items: fileList,
	}, nil
}

func (c *CaseFileStore) BuildListCaseFilesSqlizer(
	rpc options.SearchOptions,
) (sq.Sqlizer, []func(file *cases.File) any, error) {

	parentId, ok := rpc.GetFilter("case_id").(int64)
	if !ok || parentId == 0 {
		return nil, nil, errors.New("case id required")
	}
	// Begin building the base query with alias `cf`
	queryBuilder := sq.Select().
		From("storage.files AS cf").
		Where(
			sq.And{
				sq.Eq{"cf.domain_id": rpc.GetAuthOpts().GetDomainId()},
				sq.Eq{"cf.uuid": strconv.Itoa(int(parentId))},
				sq.Eq{"cf.channel": channel},
				sq.Eq{"cf.removed": nil},
			},
		).
		PlaceholderFormat(sq.Dollar)

	// Build select columns and scan plan using buildFilesSelectColumnsAndPlan
	queryBuilder, plan, err := buildFilesSelectColumnsAndPlan(queryBuilder, fileAlias, rpc.GetFields())
	if err != nil {
		return nil, nil, err
	}

	// Apply additional filters, sorting, and pagination as needed
	if len(rpc.GetIDs()) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"cf.id": rpc.GetIDs()})
	}

	// ----------Apply search by name -----------------
	if rpc.GetSearch() != "" {
		queryBuilder = util.AddSearchTerm(queryBuilder, util.Ident(caseLeft, "name"))
	}

	// -------- Apply sorting ----------
	queryBuilder = util.ApplyDefaultSorting(rpc, queryBuilder, fileDefaultSort)

	// ---------Apply paging based on Search Opts ( page ; size ) -----------------
	queryBuilder = util.ApplyPaging(rpc.GetPage(), rpc.GetSize(), queryBuilder)

	return queryBuilder, plan, nil
}

// Delete implements store.CaseFileStore.
func (c *CaseFileStore) Delete(rpc options.DeleteOptions) error {
	if rpc == nil {
		return dberr.NewDBError("postgres.case_file.delete.check_args.opts", "delete options required")
	}
	if len(rpc.GetIDs()) == 0 {
		return dberr.NewDBError("postgres.case_file.delete.check_args.id", "id required")
	}
	if rpc.GetParentID() == 0 {
		return dberr.NewDBError("postgres.case_file.delete.check_args.id", "case id required")
	}

	// convert int64 to varchar (datatype in DB)
	uuid := strconv.Itoa(int(rpc.GetParentID()))
	base := sq.
		Update(c.mainTable).
		Set("removed", true).
		Where("id = ANY(?)", rpc.GetIDs()).
		Where(sq.Eq{"domain_id": rpc.GetAuthOpts().GetDomainId()}).
		Where(sq.Eq{"uuid": uuid}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := base.ToSql()
	query = util.CompactSQL(query)

	if err != nil {
		return dberr.NewDBError("postgres.case_file.delete.parse_query.error", err.Error())
	}
	db, dbErr := c.storage.Database()
	if dbErr != nil {
		return dbErr
	}

	res, err := db.Exec(rpc, query, args...)
	if err != nil {
		return dberr.NewDBError("postgres.case_file.delete.execute.error", err.Error())
	}
	if affected := res.RowsAffected(); affected == 0 || affected > 1 {
		return dberr.NewDBNoRowsError("postgres.case_file.delete.final_check.rows")
	}
	return nil
}

func buildFilesSelectColumnsAndPlan(
	base sq.SelectBuilder,
	left string,
	fields []string,
) (sq.SelectBuilder, []func(file *cases.File) any, *dberr.DBError) {
	var (
		plan           []func(file *cases.File) any
		createdByAlias string
		joinCreatedBy  = func() {
			if createdByAlias != "" {
				return
			}
			createdByAlias = caseFileCreatedByAlias
			base = base.LeftJoin(fmt.Sprintf("directory.wbt_user %s ON %[1]s.id = %s.uploaded_by", caseFileCreatedByAlias, left))
		}
		authorAlias string
		joinAuthor  = func() {
			if authorAlias != "" {
				return
			}
			joinCreatedBy()
			authorAlias = caseFileAuthorAlias
			base = base.LeftJoin(fmt.Sprintf("contacts.contact %s ON %[1]s.id = %s.contact_id", authorAlias, createdByAlias))
		}
	)

	for _, field := range fields {
		switch field {
		case "id":
			base = base.Column(util.Ident(left, "id"))
			plan = append(plan, func(file *cases.File) any {
				return &file.Id
			})
		case "created_by":
			joinCreatedBy()
			base = base.Column(fmt.Sprintf("ROW(%[1]s.id, %[1]s.name)::text created_by", caseFileCreatedByAlias))
			plan = append(plan, func(file *cases.File) any {
				return scanner.ScanRowLookup(&file.CreatedBy)
			})
		case "created_at":
			base = base.Column(util.Ident(left, "uploaded_at"))
			plan = append(plan, func(file *cases.File) any {
				return scanner.ScanTimestamp(&file.CreatedAt)
			})
		case "size":
			base = base.Column(util.Ident(left, "size"))
			plan = append(plan, func(file *cases.File) any {
				return &file.Size
			})
		case "mime":
			base = base.Column(util.Ident(left, "mime_type"))
			plan = append(plan, func(file *cases.File) any {
				return &file.Mime
			})
		case "name":
			base = base.Column(util.Ident(left, "view_name"))
			plan = append(plan, func(file *cases.File) any {
				return &file.Name
			})
		// case "url":
		//	base = base.Column(store.Ident(left, "url"))
		//	plan = append(plan, func(file *cases.File) any {
		//		return &file.Url
		//	})
		case "author":
			joinAuthor()
			base = base.Column(fmt.Sprintf(`ROW(%[1]s.id, %[1]s.common_name)::text author`, caseFileAuthorAlias))
			plan = append(plan, func(file *cases.File) any {
				return scanner.ScanRowLookup(&file.Author)
			})
		default:
			return base, nil, dberr.NewDBError("postgres.case_file.build_file_select.cycle_fields.unknown", fmt.Sprintf("%s field is unknown", field))
		}
	}

	if len(plan) == 0 {
		return base, nil, dberr.NewDBError("postgres.case_file.build_file_select.final_check.unknown", "no resulting columns")
	}

	return base, plan, nil
}

func buildFilesSelectAsSubquery(fields []string, caseAlias string) (sq.SelectBuilder, []func(file *cases.File) any, int, *dberr.DBError) {
	alias := "files"
	if caseAlias == alias {
		alias = "sub_" + alias
	}
	base := sq.
		Select().
		From("storage.files " + alias).
		Where(fmt.Sprintf("%s = %s::text", util.Ident(alias, "uuid"), util.Ident(caseAlias, "id"))).
		Where(fmt.Sprintf("%s = '%s'", util.Ident(alias, "channel"), channel))
	base = util.ApplyPaging(1, defaults.DefaultSearchSize, base)

	base, scanPlan, dbErr := buildFilesSelectColumnsAndPlan(base, alias, fields)
	if dbErr != nil {
		return base, nil, 0, dbErr
	}

	return base, scanPlan, 0, nil
}

// NewCaseFileStore initializes a new CaseFileStore.
func NewCaseFileStore(store *Store) (store.CaseFileStore, error) {
	if store == nil {
		return nil, dberr.NewDBError("postgres.new_case_file.check.bad_arguments", "error creating case file interface, main store is nil")
	}
	return &CaseFileStore{storage: store, mainTable: "storage.files"}, nil
}
