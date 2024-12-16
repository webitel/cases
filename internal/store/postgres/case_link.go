package postgres

import (
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	_go "github.com/webitel/cases/api/cases"
	dberr "github.com/webitel/cases/internal/error"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/internal/store/scanner"
	"github.com/webitel/cases/model"
	"github.com/webitel/cases/util"
	"net/url"
)

const (
	caseLinkCreatedByAlias = "cb"
	caseLinkUpdatedByAlias = "ub"
	caseLinkAuthorAlias    = "au"
)

type CaseLinkStore struct {
	storage   store.Store
	mainTable string
}

var (
	CaseLinkFields = []string{
		"created_by", "created_at", "updated_by", "updated_at", "id", "ver", "author", "name", "url",
	}
)

// Create implements store.CaseLinkStore.
func (l *CaseLinkStore) Create(rpc *model.CreateOptions, add *_go.InputCaseLink) (*_go.CaseLink, error) {
	if rpc == nil {
		return nil, dberr.NewDBError("postgres.case_link.create.check_args.opts", "create options required")
	}
	dbErr := ValidateLinkCreate(rpc.ParentID, add)
	if dbErr != nil {
		return nil, dbErr
	}
	if len(rpc.Fields) == 0 {
		rpc.Fields = CaseLinkFields
	}
	base, plan, dbErr := buildCreateLinkQuery(rpc, add)
	if dbErr != nil {
		return nil, dbErr
	}
	db, dbErr := l.storage.Database()
	if dbErr != nil {
		return nil, dbErr
	}
	query, args, goErr := base.ToSql()
	if goErr != nil {
		return nil, dberr.NewDBError("postgres.case_link.create.convert_to_sql.error", goErr.Error())
	}
	row := db.QueryRow(rpc.Context, store.CompactSQL(query), args...)
	res, goErr := l.scanLink(row, plan)
	if goErr != nil {
		return nil, dberr.NewDBError("postgres.case_link.create.scan.error", goErr.Error())
	}

	return res, nil
}

// Delete implements store.CaseLinkStore.
func (l *CaseLinkStore) Delete(opts *model.DeleteOptions) error {
	if opts == nil {
		return dberr.NewDBError("postgres.case_link.delete.check_args.opts", "delete options required")
	}
	if opts.ID == 0 {
		return dberr.NewDBError("postgres.case_link.delete.check_args.id", "id required")
	}
	base := squirrel.
		Delete(l.mainTable).
		Where("id = ?", opts.ID).
		PlaceholderFormat(squirrel.Dollar)
	query, args, err := base.ToSql()
	if err != nil {
		return dberr.NewDBError("postgres.case_link.delete.parse_query.error", err.Error())
	}
	db, dbErr := l.storage.Database()
	if dbErr != nil {
		return dbErr
	}

	res, err := db.Exec(opts.Context, query, args...)
	if err != nil {
		return dberr.NewDBError("postgres.case_link.delete.execute.error", err.Error())
	}
	if affected := res.RowsAffected(); affected == 0 || affected > 1 {
		return dberr.NewDBError("postgres.case_link.delete.final_check.rows", "wrong filters for deleting")
	}
	return nil

}

// List implements store.CaseLinkStore.
func (l *CaseLinkStore) List(opts *model.SearchOptions) (*_go.CaseLinkList, error) {
	// validate
	if opts == nil {
		return nil, dberr.NewDBError("postgres.case_link.list.check_args.opts", "search options required")
	}
	if opts.ParentId == 0 && len(opts.IDs) == 0 {
		return nil, dberr.NewDBError("postgres.case_link.list.check_args.parent_id", "case id required")
	}
	db, dbErr := l.storage.Database()
	if dbErr != nil {
		return nil, dbErr
	}

	// build query
	base := squirrel.
		Select().
		From(l.mainTable).
		PlaceholderFormat(squirrel.Dollar)
	if opts.ParentId != 0 {
		base = base.Where(fmt.Sprintf("%s = ?", store.Ident(l.mainTable, "case_id")), opts.ParentId)
	}
	if len(opts.IDs) != 0 {
		base = base.Where(fmt.Sprintf("%s = any(?)", store.Ident(l.mainTable, "id")), opts.IDs)
	}
	base = store.ApplyPaging(opts.GetPage(), opts.GetSize(), base)
	base = store.ApplyDefaultSorting(opts, base)
	base, plan, dbErr := buildLinkSelectColumnsAndPlan(base, l.mainTable, opts.Fields)
	if dbErr != nil {
		return nil, dbErr
	}

	// execute
	query, args, err := base.ToSql()
	if err != nil {
		return nil, dberr.NewDBError("postgres.case_link.list.convert_sql.error", err.Error())
	}

	rows, err := db.Query(opts.Context, query, args...)
	if err != nil {
		return nil, dberr.NewDBError("postgres.case_link.list.execute.error", err.Error())
	}
	// result
	links, err := l.scanLinks(rows, plan)
	if err != nil {
		return nil, err
	}
	var res _go.CaseLinkList
	if opts.GetSize() > 0 && len(links) > int(opts.GetSize()) {
		res.Next = true
		links = links[:len(links)-1]
	}
	res.Page = res.GetPage()
	res.Items = links
	return &res, nil
}

// Update implements store.CaseLinkStore.
func (l *CaseLinkStore) Update(opts *model.UpdateOptions, upd *_go.InputCaseLink) (*_go.CaseLink, error) {
	if opts == nil {
		return nil, dberr.NewDBError("postgres.case_link.update.check_args.opts", "update options required")
	}
	query, plan, dbErr := buildUpdateLinkQuery(opts, upd)
	if dbErr != nil {
		return nil, dbErr
	}
	db, dbErr := l.storage.Database()
	if dbErr != nil {
		return nil, dbErr
	}
	slct, args, err := query.ToSql()
	if err != nil {
		return nil, dberr.NewDBError("postgres.case_link.update.convert_to_sql.error", err.Error())
	}
	row := db.QueryRow(opts.Context, slct, args...)
	res, err := l.scanLink(row, plan)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, dberr.NewDBNotFoundError("postgres.case_link.update.scan_ver.not_found", "Link not found")
		}
		return nil, err
	}
	return res, nil
}

func NewLinkCaseStore(store store.Store) (store.CaseLinkStore, error) {
	if store == nil {
		return nil, dberr.NewDBError("postgres.new_link_case.check.bad_arguments",
			"error creating link case interface to the comment_case table, main store is nil")
	}
	return &CaseLinkStore{storage: store, mainTable: "cases.case_link"}, nil
}

func buildLinkSelectColumnsAndPlan(base squirrel.SelectBuilder, left string, fields []string) (squirrel.SelectBuilder, []func(link *_go.CaseLink) any, *dberr.DBError) {
	var (
		plan           []func(link *_go.CaseLink) any
		createdByAlias string
		joinCreatedBy  = func() {
			if createdByAlias != "" {
				return
			}
			createdByAlias = caseLinkCreatedByAlias
			base = base.LeftJoin(fmt.Sprintf("directory.wbt_user %s ON %[1]s.id = %s.created_by", createdByAlias, left))
			return
		}
		updatedByAlias string
		joinUpdatedBy  = func() {
			if updatedByAlias != "" {
				return
			}
			updatedByAlias = caseLinkUpdatedByAlias
			base = base.LeftJoin(fmt.Sprintf("directory.wbt_user %s ON %[1]s.id = %s.updated_by", updatedByAlias, left))
			return
		}
		authorAlias string
		joinAuthor  = func() {
			if authorAlias != "" {
				return
			}
			joinCreatedBy()
			authorAlias = caseLinkAuthorAlias
			base = base.LeftJoin(fmt.Sprintf("contacts.contact %s ON %[1]s.id = %s.contact_id", authorAlias, createdByAlias))
			return
		}
	)
	if len(fields) == 0 {
		fields = CaseLinkFields
	}

	for _, field := range fields {
		switch field {
		case "id":
			base = base.Column(store.Ident(left, "id"))
			plan = append(plan, func(link *_go.CaseLink) any {
				return &link.Id
			})
		case "ver":
			base = base.Column(store.Ident(left, "ver"))
			plan = append(plan, func(link *_go.CaseLink) any {
				return &link.Ver
			})
		case "created_by":
			joinCreatedBy()
			base = base.Column(fmt.Sprintf("ROW(%[1]s.id, %[1]s.name)::text created_by", createdByAlias))
			plan = append(plan, func(link *_go.CaseLink) any {
				return scanner.ScanRowLookup(&link.CreatedBy)
			})
		case "created_at":
			base = base.Column(store.Ident(left, "created_at"))
			plan = append(plan, func(link *_go.CaseLink) any {
				return scanner.ScanTimestamp(&link.CreatedAt)
			})
		case "updated_by":
			joinUpdatedBy()
			base = base.Column(fmt.Sprintf("ROW(%[1]s.id, %[1]s.name)::text updated_by", updatedByAlias))
			plan = append(plan, func(link *_go.CaseLink) any {
				return scanner.ScanRowLookup(&link.UpdatedBy)
			})
		case "updated_at":
			base = base.Column(store.Ident(left, "updated_at"))
			plan = append(plan, func(link *_go.CaseLink) any {
				return scanner.ScanTimestamp(&link.UpdatedAt)
			})
		case "name":
			base = base.Column(store.Ident(left, "name"))
			plan = append(plan, func(link *_go.CaseLink) any {
				return scanner.ScanText(&link.Name)
			})
		case "url":
			base = base.Column(store.Ident(left, "url"))
			plan = append(plan, func(link *_go.CaseLink) any {
				return &link.Url
			})
		case "author":
			joinAuthor()
			base = base.Column(fmt.Sprintf(`ROW(%[1]s.id, %[1]s.common_name)::text author`, authorAlias))
			plan = append(plan, func(link *_go.CaseLink) any {
				return scanner.ScanRowLookup(&link.Author)
			})
		default:
			return base, nil, dberr.NewDBError("postgres.case_link.build_link_select.cycle_fields.unknown", fmt.Sprintf("%s field is unknown", field))
		}
	}
	if len(plan) == 0 {
		return base, nil, dberr.NewDBError("postgres.case_link.build_link_select.final_check.unknown", "no resulting columns")
	}
	return base, plan, nil
}

func buildCreateLinkQuery(rpc *model.CreateOptions, add *_go.InputCaseLink) (squirrel.Sqlizer, []func(link *_go.CaseLink) any, *dberr.DBError) {
	insertAlias := "i"
	insert := squirrel.
		Insert("cases.case_link").
		Columns("created_by", "updated_by", "name", "url", "case_id", "dc").
		Values(rpc.Session.GetUserId(), rpc.Session.GetUserId(), add.GetName(), add.GetUrl(), rpc.ParentID, rpc.Session.GetDomainId()).
		Suffix("RETURNING *")
	// select
	query, args, _ := store.FormAsCTE(insert, insertAlias)
	base := squirrel.Select().From(insertAlias).Prefix(query, args...).PlaceholderFormat(squirrel.Dollar).Where("i.created_by = ?", 10454)
	// build plan from columns
	return buildLinkSelectColumnsAndPlan(base, insertAlias, rpc.Fields)
}

func buildUpdateLinkQuery(opts *model.UpdateOptions, add *_go.InputCaseLink) (squirrel.Sqlizer, []func(link *_go.CaseLink) any, *dberr.DBError) {
	if len(opts.Etags) == 0 {
		return nil, nil, dberr.NewDBError("postgres.case_link.update.etag.empty", "link etag required")
	}
	if len(opts.Mask) == 0 {
		return nil, nil, dberr.NewDBError("postgres.case_link.update.mask.empty", "link update mask required")
	}
	tid := opts.Etags[0]
	// insert
	update := squirrel.
		Update("cases.case_link").
		Set("updated_by", opts.Session.GetUserId()).
		Set("updated_at", opts.Time).
		Set("ver", squirrel.Expr("ver+1")).
		Where("id = ?", tid.GetOid()).
		Where("ver = ?", tid.GetVer()).
		Suffix("RETURNING *").
		PlaceholderFormat(squirrel.Dollar)
	for _, field := range opts.Mask {
		switch field {
		case "url":
			_, err := url.Parse(add.Url)
			if err != nil {
				return nil, nil, dberr.NewDBError("postgres.case_link.build_update_query.url.validate.error", err.Error())
			}
			update = update.Set("url", add.Url)
		case "name":
			update = update.Set("name", add.Name)
		}
	}
	prefixAlias := "upd"
	prefix, args, err := store.FormAsCTE(update, prefixAlias)
	if err != nil {
		return nil, nil, dberr.NewDBError("postgres.case_link.build_update_query.form_cte.error", err.Error())
	}
	slct := squirrel.Select().Prefix(prefix, args...).From(prefixAlias)
	// build plan from columns
	return buildLinkSelectColumnsAndPlan(slct, prefixAlias, opts.Fields)
}

func (l *CaseLinkStore) scanLinks(rows pgx.Rows, plan []func(link *_go.CaseLink) any) ([]*_go.CaseLink, error) {
	var res []*_go.CaseLink

	for rows.Next() {
		link, err := l.scanLink(rows, plan)
		if err != nil {
			return nil, err
		}
		res = append(res, link)
	}
	return res, nil
}

func (l *CaseLinkStore) scanLink(row pgx.Row, plan []func(link *_go.CaseLink) any) (*_go.CaseLink, error) {
	var link _go.CaseLink
	var scanPlan []any
	for _, scan := range plan {
		scanPlan = append(scanPlan, scan(&link))
	}
	err := row.Scan(scanPlan...)
	if err != nil {
		return nil, err
	}

	return &link, nil
}

func buildLinkSelectAsSubquery(opts *model.SearchOptions, caseAlias string) (updatedBase squirrel.SelectBuilder, scanPlan []func(link *_go.CaseLink) any, filtersApplied int, dbErr *dberr.DBError) {
	alias := "links"
	if caseAlias == alias {
		alias = "sub_" + alias
	}
	base := squirrel.
		Select().
		From("cases.case_link " + alias).
		Where(fmt.Sprintf("%s = %s", store.Ident(alias, "case_id"), store.Ident(caseAlias, "id")))
	if opts.Search != "" {
		base = store.AddSearchTerm(base, opts.Search, store.Ident(alias, "name"), store.Ident(alias, "url"))
	}

	base, plan, dbErr := buildLinkSelectColumnsAndPlan(base, alias, opts.Fields)
	if dbErr != nil {
		return base, nil, 0, dbErr
	}
	base, applied, dbErr := applyCaseLinkFilters(opts, base, alias)
	if dbErr != nil {
		return base, nil, 0, dbErr
	}
	base = store.ApplyPaging(opts.GetPage(), opts.GetSize(), base)

	return base, plan, applied, nil
}

func applyCaseLinkFilters(opts *model.SearchOptions, base squirrel.SelectBuilder, alias string) (updatedBase squirrel.SelectBuilder, filtersApplied int, err *dberr.DBError) {
	if opts == nil || len(opts.Filter) == 0 {
		return base, 0, nil
	}

	for column, value := range opts.Filter {
		if !util.ContainsStringIgnoreCase(opts.Fields, column) {
			continue
		}
		switch column {
		case "created_by":
			switch v := value.(type) {
			case int64, int, int32, *int64, *int, *int32:
				base = base.Where(fmt.Sprintf("%s = ?", store.Ident(caseLinkCreatedByAlias, "id")), v)
			case string, *string:
				// apply search
				//base = store.AddSearchTerm(base, )
			}
		case "author":
			switch v := value.(type) {
			case int64, int, int32, *int64, *int, *int32:
				//
				base = base.Where(fmt.Sprintf("%s = ?", store.Ident(caseLinkUpdatedByAlias, "id")), v)
			case string, *string:
				// apply search
				//base = store.AddSearchTerm(base, )
			}
		case "updated_by":
			switch v := value.(type) {
			case int64, int, int32, *int64, *int, *int32:
				//
				base = base.Where(fmt.Sprintf("%s = ?", store.Ident(caseLinkAuthorAlias, "id")), v)
			case string, *string:
				// apply search
				//base = store.AddSearchTerm(base, )
			}
		}
		filtersApplied++
	}
	return base, filtersApplied, nil
}

func ValidateLinkCreate(caseId int64, input *_go.InputCaseLink) *dberr.DBError {
	if caseId <= 0 {
		return dberr.NewDBError("postgres.case_link.validate_create.check_args.case_id", "case id required")
	}
	if input == nil || input.Url == "" {
		return dberr.NewDBError("postgres.case_link.validate_create.check_args.input", "input for link required")
	}
	_, err := url.Parse(input.Url)
	if err != nil {
		return dberr.NewDBError("postgres.case_link.validate_create.validate_url.error", err.Error())
	}
	return nil
}
