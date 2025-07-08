package postgres

import (
	"fmt"
	"strconv"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
	"github.com/webitel/cases/internal/model/options/defaults"
	storeutil "github.com/webitel/cases/internal/store/util"

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
	chatChannel             = "chat"
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
		fields = []string{"id", "name", "size", "mime", "created_at", "created_by", "author", "url"}
	}
	caseIDFilters := rpc.GetFilter("case_id")
	if len(caseIDFilters) == 0 {
		return sq.Select(), errors.New("case id required")
	}
	var ctes []*storeutil.CTE

	// Build the CTE for chat files
	conversationsWithFiles := sq.Select("communication_id chat_id").
		From("cases.case_communication case_comms").
		LeftJoin("call_center.cc_communication comm_type ON case_comms.communication_type = comm_type.id").
		Where("comm_type.channel = 'Messaging'").
		Where("case_id = ?", caseIDFilters[0].Value)
	ctes = append(ctes, storeutil.NewCTE("connected_chats", conversationsWithFiles))

	// Build the CTE for chat messages with files
	messageWithFiles := sq.Select(
		"file_id",
		"jsonb_build_object('id', ch.user_id, 'name', COALESCE(usr.name, bot.name, cli.name), 'type', COALESCE(ch.type, 'bot')) as author",
	).
		From("chat.message m").
		LeftJoin("chat.channel ch ON ch.id = m.channel_id").
		LeftJoin("chat.bot bot ON bot.id = ch.user_id").
		LeftJoin("directory.wbt_user usr ON usr.id = ch.user_id").
		LeftJoin("chat.client cli ON cli.id = ch.user_id").
		Where("m.conversation_id::text = ANY (SELECT chat_id FROM connected_chats)").
		Where("file_id IS NOT NULL")
	ctes = append(ctes, storeutil.NewCTE("chat_upload", messageWithFiles))

	directFiles := sq.Select(
		"f.id as file_id",
		"jsonb_build_object('id', usr.id, 'name', COALESCE(usr.name, usr.username), 'type', 'webitel') as author",
	).
		From("storage.files f").
		LeftJoin("directory.wbt_user usr ON f.uploaded_by = usr.id").
		Where("uuid = ?::text", caseIDFilters[0].Value).
		Where("channel = 'case'")
	ctes = append(ctes, storeutil.NewCTE("direct_upload", directFiles))

	// Union the CTEs for chat files and direct files
	fileMessages := sq.Expr("SELECT * FROM chat_upload UNION SELECT * FROM direct_upload")
	ctes = append(ctes, storeutil.NewCTE("union_types", fileMessages))

	ctesQuery, ctesArgs, err := storeutil.FormAsCTEs(ctes)
	if err != nil {
		return conversationsWithFiles, err
	}
	unionAlias := "files"

	// Begin building the base query with alias `cf`
	queryBuilder := sq.Select().
		From("storage.files AS cf").
		InnerJoin(fmt.Sprintf("union_types %s ON %s = %s", unionAlias, storeutil.Ident(unionAlias, "file_id"), storeutil.Ident(fileAlias, "id"))).
		PlaceholderFormat(sq.Dollar).
		Prefix(ctesQuery, ctesArgs...)
	// Build select columns and scan plan using buildFilesSelectColumnsAndPlan
	queryBuilder, err = buildCaseFileSelectColumns(queryBuilder, rpc.GetFields(), fileAlias, unionAlias)
	if err != nil {
		return queryBuilder, err
	}

	// Apply additional filters, sorting, and pagination as needed
	if len(rpc.GetIDs()) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"cf.id": rpc.GetIDs()})
	}

	// ----------Apply search by name -----------------
	if rpc.GetSearch() != "" {
		queryBuilder = storeutil.AddSearchTerm(queryBuilder, storeutil.Ident(caseLeft, "name"))
	}

	// -------- Apply sorting ----------
	queryBuilder = storeutil.ApplyDefaultSorting(rpc, queryBuilder, fileDefaultSort)

	// ---------Apply paging based on Search Opts ( page ; size ) -----------------
	queryBuilder = storeutil.ApplyPaging(rpc.GetPage(), rpc.GetSize(), queryBuilder)

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
	updateSQL = storeutil.CompactSQL(updateSQL)

	if err != nil {
		return nil, ParseError(err)
	}

	cteSQL := "WITH deleted AS (" + updateSQL + ")"
	selectBuilder, err := buildCaseFileSelectColumns(
		sq.Select().Prefix(cteSQL).From("deleted cf"),
		fields,
		fileAlias,
		"",
	)
	if err != nil {
		return nil, ParseError(err)
	}
	selectBuilder = selectBuilder.PlaceholderFormat(sq.Dollar)

	query, _, err := selectBuilder.ToSql()
	if err != nil {
		return nil, ParseError(err)
	}
	query = storeutil.CompactSQL(query)

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
	authorCTEAlias string,
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
		case "created_by":
			if authorCTEAlias != "" {
				base = base.Column(storeutil.Ident(authorCTEAlias, "author created_by"))
			} else {
				cb := caseFileCreatedByAlias
				joinCreatedBy(cb)
				base = base.Column(fmt.Sprintf("jsonb_build_object('id', %s.id, 'name', %[1]s.name, 'type', 'webitel') created_by", cb))
			}
		case "url":
			base = base.Column(fmt.Sprintf("%s.file_url AS url", tableAlias))
		default:
			return base, errors.New("unknown field: " + field)
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
		Where(fmt.Sprintf("%s = %s::text", storeutil.Ident(alias, "uuid"), storeutil.Ident(caseAlias, "id"))).
		Where(fmt.Sprintf("%s = '%s'", storeutil.Ident(alias, "channel"), channel))
	base = storeutil.ApplyPaging(1, defaults.DefaultSearchSize, base)

	base, dbErr := buildCaseFileSelectColumns(base, fields, alias, "")
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
