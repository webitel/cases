package postgres

import (
	"fmt"
	"strconv"

	sq "github.com/Masterminds/squirrel"
	"github.com/webitel/cases/api/cases"
	dberr "github.com/webitel/cases/internal/error"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/internal/store/scanner"
	"github.com/webitel/cases/model"
	util "github.com/webitel/cases/util"
)

type CaseFileStore struct {
	storage store.Store
}

const (
	// Alias for the storage.files table
	fileAlias              = "cf"
	channel                = "case"
	fileDefaultSort        = "created_at"
	caseFileAuthorAlias    = "au"
	caseFileCreatedByAlias = "cb"
)

// List implements store.CaseFileStore for listing case files.
func (c *CaseFileStore) List(rpc *model.SearchOptions) (*cases.CaseFileList, error) {
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
	rows, err := d.Query(rpc.Context, query, args...)
	if err != nil {
		return nil, dberr.NewDBInternalError("store.case_file.list.execution_error", err)
	}
	defer rows.Close()

	var fileList []*cases.File
	lCount := 0
	next := false
	fetchAll := rpc.GetSize() == -1

	for rows.Next() {
		if !fetchAll && lCount >= int(rpc.GetSize()) {
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
		Page:  int64(rpc.Page),
		Next:  next,
		Items: fileList,
	}, nil
}

func (c *CaseFileStore) BuildListCaseFilesSqlizer(
	rpc *model.SearchOptions,
) (sq.Sqlizer, []func(file *cases.File) any, error) {
	// Begin building the base query with alias `cf`
	queryBuilder := sq.Select().
		From("storage.files AS cf").
		Where(
			sq.And{
				sq.Eq{"cf.domain_id": rpc.Auth.GetDomainId()},
				sq.Eq{"cf.uuid": strconv.Itoa(int(rpc.ParentId))},
				sq.Eq{"cf.channel": channel},
			},
		).
		PlaceholderFormat(sq.Dollar)

	// Ensure necessary fields are included
	rpc.Fields = util.EnsureIdField(rpc.Fields)

	// Build select columns and scan plan using buildFilesSelectColumnsAndPlan
	queryBuilder, plan, err := buildFilesSelectColumnsAndPlan(queryBuilder, fileAlias, rpc.Fields)
	if err != nil {
		return nil, nil, err
	}

	// Apply additional filters, sorting, and pagination as needed
	if len(rpc.IDs) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"cf.id": rpc.IDs})
	}

	// ----------Apply search by name -----------------
	if rpc.Search != "" {
		queryBuilder = store.AddSearchTerm(queryBuilder, store.Ident(caseLeft, "name"))
	}

	// -------- Apply sorting ----------
	queryBuilder = store.ApplyDefaultSorting(rpc, queryBuilder, fileDefaultSort)

	// ---------Apply paging based on Search Opts ( page ; size ) -----------------
	queryBuilder = store.ApplyPaging(rpc.GetPage(), rpc.GetSize(), queryBuilder)

	return queryBuilder, plan, nil
}

// Helper function to build the select columns and scan plan based on the fields requested.
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
			base = base.LeftJoin(fmt.Sprintf("directory.wbt_user %s ON %[1]s.id = %s.created_by", caseFileCreatedByAlias, left))
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
			base = base.Column(store.Ident(left, "id"))
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
			base = base.Column(store.Ident(left, "uploaded_at"))
			plan = append(plan, func(file *cases.File) any {
				return scanner.ScanTimestamp(&file.CreatedAt)
			})
		case "size":
			base = base.Column(store.Ident(left, "size"))
			plan = append(plan, func(file *cases.File) any {
				return &file.Size
			})
		case "mime":
			base = base.Column(store.Ident(left, "mime_type"))
			plan = append(plan, func(file *cases.File) any {
				return &file.Mime
			})
		case "name":
			base = base.Column(store.Ident(left, "name"))
			plan = append(plan, func(file *cases.File) any {
				return &file.Name
			})
		case "url":
			base = base.Column(store.Ident(left, "url"))
			plan = append(plan, func(file *cases.File) any {
				return &file.Url
			})
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

func buildFilesSelectAsSubquery(opts *model.SearchOptions, caseAlias string) (sq.SelectBuilder, []func(file *cases.File) any, int, *dberr.DBError) {
	alias := "files"
	if caseAlias == alias {
		alias = "sub_" + alias
	}
	base := sq.
		Select().
		From("storage.files "+alias).
		Where(fmt.Sprintf("%s = %s::text", store.Ident(alias, "uuid"), store.Ident(caseAlias, "id"))).
		Where(fmt.Sprintf("%s = ?", store.Ident(alias, "channel")), channel)
	base = store.ApplyPaging(opts.GetPage(), opts.GetSize(), base)

	base, scanPlan, dbErr := buildFilesSelectColumnsAndPlan(base, alias, opts.Fields)
	if dbErr != nil {
		return base, nil, 0, dbErr
	}

	return base, scanPlan, 0, nil
}

// NewCaseFileStore initializes a new CaseFileStore.
func NewCaseFileStore(store store.Store) (store.CaseFileStore, error) {
	if store == nil {
		return nil, dberr.NewDBError("postgres.new_case_file.check.bad_arguments", "error creating case file interface, main store is nil")
	}
	return &CaseFileStore{storage: store}, nil
}
