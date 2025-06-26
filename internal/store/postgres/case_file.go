package postgres

import (
	"fmt"
	"strconv"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
	"github.com/webitel/cases/internal/model/options/defaults"
	"github.com/webitel/cases/internal/store/util"

	sq "github.com/Masterminds/squirrel"
	"github.com/webitel/cases/internal/store"
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
func (c *CaseFileStore) List(rpc options.Searcher) ([]*model.CaseFile, error) {
	// Connect to the database
	d, dbErr := c.storage.Database()
	if dbErr != nil {
		return nil, dbErr
	}

	// Build the query and plan builder using BuildListCaseFilesSqlizer
	queryBuilder, err := c.BuildListCaseFilesSqlizer(rpc)
	if err != nil {
		return nil, err
	}

	// Convert the query to SQL
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, ParseError(err)
	}
	var items []*model.CaseFile
	if err := pgxscan.Select(rpc, d, &items, query, args...); err != nil {
		return nil, ParseError(err)
	}
	return items, nil
}

func (c *CaseFileStore) BuildListCaseFilesSqlizer(
	rpc options.Searcher,
) (sq.SelectBuilder, error) {
	fields := rpc.GetFields()
	if len(fields) == 0 {
		fields = []string{"id", "name", "size", "mime", "created_at", "created_by", "author"}
	}

	parentId, ok := rpc.GetFilter("case_id").(int64)
	if !ok || parentId == 0 {
		return sq.SelectBuilder{}, errors.InvalidArgument("case id required")
	}
	// Begin building the base query with alias `cf`
	queryBuilder, err := buildCaseFileSelectColumns(
		sq.Select().From("storage.files AS cf"),
		fields,
		fileAlias,
	)
	if err != nil {
		return sq.SelectBuilder{}, err
	}

	queryBuilder = queryBuilder.Where(
		sq.And{
			sq.Eq{"cf.domain_id": rpc.GetAuthOpts().GetDomainId()},
			sq.Eq{"cf.uuid": strconv.Itoa(int(parentId))},
			sq.Eq{"cf.channel": channel},
			sq.Eq{"cf.removed": nil},
		},
	).
		PlaceholderFormat(sq.Dollar)

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

	return queryBuilder, nil
}

// Delete implements store.CaseFileStore.
func (c *CaseFileStore) Delete(rpc options.Deleter) (*model.CaseFile, error) {
	if rpc == nil {
		return nil, errors.New("delete options required")
	}
	if len(rpc.GetIDs()) == 0 {
		return nil, errors.New("id required")
	}
	if rpc.GetParentID() == 0 {
		return nil, errors.New("case id required")
	}

	fields := []string{"id", "name", "size", "mime", "created_at", "created_by", "author"}

	ids := rpc.GetIDs()
	if len(ids) == 1 {
		// Optionally, you can convert a single value to a slice
		ids = []int64{ids[0]}
	}

	// convert int64 to varchar (datatype in DB)
	uuid := strconv.Itoa(int(rpc.GetParentID()))
	updateBuilder := sq.
		Update(c.mainTable).
		Set("removed", true).
		Where("id = ANY(?)", ids).
		Where(sq.Eq{"domain_id": rpc.GetAuthOpts().GetDomainId()}).
		Where(sq.Eq{"uuid": uuid}).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar)

	updateSQL, updateArgs, err := updateBuilder.ToSql()
	updateSQL = util.CompactSQL(updateSQL)

	if err != nil {
		return nil, ParseError(err)
	}

	cteSQL := "WITH deleted AS (" + updateSQL + ")"
	selectBuilder, err := buildCaseFileSelectColumns(
		sq.Select().Prefix(cteSQL).From("deleted cf"),
		fields,
		fileAlias,
	)
	if err != nil {
		return nil, ParseError(err)
	}
	selectBuilder = selectBuilder.PlaceholderFormat(sq.Dollar)

	query, _, err := selectBuilder.ToSql()
	if err != nil {
		return nil, ParseError(err)
	}
	query = util.CompactSQL(query)

	db, dbErr := c.storage.Database()
	if dbErr != nil {
		return nil, dbErr
	}

	var result model.CaseFile
	err = pgxscan.Get(rpc, db, &result, query, updateArgs...)
	if err != nil {
		if pgxscan.NotFound(err) {
			return nil, errors.New("case file not found")
		}
		return nil, ParseError(err)
	}

	return &result, nil
}

func buildCaseFileSelectColumns(
	base sq.SelectBuilder,
	fields []string,
	tableAlias string,
) (sq.SelectBuilder, error) {
	var (
		createdByAlias string
		joinCreatedBy  = func(alias string) string {
			if createdByAlias != "" {
				return createdByAlias
			}
			base = base.LeftJoin(fmt.Sprintf("directory.wbt_user %s ON %s.uploaded_by = %s.id", alias, tableAlias, alias))
			createdByAlias = alias
			return alias
		}
		authorAlias string
		joinAuthor  = func(alias string) string {
			if authorAlias != "" {
				return authorAlias
			}
			joinCreatedBy(caseFileCreatedByAlias)
			authorAlias = alias
			base = base.LeftJoin(fmt.Sprintf("contacts.contact %s ON %s.contact_id = %s.id", alias, caseFileCreatedByAlias, alias))
			return alias
		}
	)
	base = base.Column(fmt.Sprintf("%s.id", tableAlias))
	for _, field := range fields {
		switch field {
		case "id":
			// already set
		case "name":
			base = base.Column(fmt.Sprintf("%s.view_name AS name", tableAlias))
		case "size":
			base = base.Column(fmt.Sprintf("%s.size", tableAlias))
		case "mime":
			base = base.Column(fmt.Sprintf("%s.mime_type AS mime", tableAlias))
		case "created_at":
			base = base.Column(fmt.Sprintf("%s.uploaded_at AS created_at", tableAlias))
			/* 		case "url":
			base = base.Column(fmt.Sprintf("%s.url", tableAlias)) */
		case "created_by":
			cb := caseFileCreatedByAlias
			joinCreatedBy(cb)
			base = base.Column(fmt.Sprintf("%s.id AS created_by_id", cb))
			base = base.Column(fmt.Sprintf("%s.name AS created_by_name", cb))
		case "author":
			au := caseFileAuthorAlias
			joinAuthor(au)
			base = base.Column(fmt.Sprintf("%s.id AS contact_id", au))
			base = base.Column(fmt.Sprintf("%s.common_name AS contact_name", au))
		default:
			return base, errors.New("unknown field: "+field)
		}
	}
	return base, nil
}

func buildFilesSelectAsSubquery(fields []string, caseAlias string) (sq.SelectBuilder, error) {
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

	base, dbErr := buildCaseFileSelectColumns(base, fields, alias)
	if dbErr != nil {
		return base, dbErr
	}

	return base, nil
}

// NewCaseFileStore initializes a new CaseFileStore.
func NewCaseFileStore(store *Store) (store.CaseFileStore, error) {
	if store == nil {
		return nil, errors.New("error creating case file interface, main store is nil")
	}
	return &CaseFileStore{storage: store, mainTable: "storage.files"}, nil
}
