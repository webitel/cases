package postgres

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	_go "github.com/webitel/cases/api/cases"
	dberr "github.com/webitel/cases/internal/error"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/internal/store/scanner"
	"github.com/webitel/cases/model"
	"net/url"
)

type CaseLinkStore struct {
	storage   store.Store
	mainTable string
}

var CaseLinkFields = []string{
	"created_by", "created_at", "updated_by", "updated_at", "id", "ver", "author", "name", "url", "etag",
}

type LinkScan func(link *_go.CaseLink) any

var ident = func(left, right string) string {
	return fmt.Sprintf("%s.%s", left, right)
}

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
	base, plan, dbErr := buildCreateLinkQuery(rpc, rpc.ParentID, add)
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
	res, goErr := scanLink(row, plan)
	if goErr != nil {
		return nil, dberr.NewDBError("postgres.case_link.create.scan.error", goErr.Error())
	}

	return res, nil
}

// Delete implements store.CaseLinkStore.
func (l *CaseLinkStore) Delete(req *model.DeleteOptions) (*_go.CaseLink, error) {
	panic("unimplemented")
}

// List implements store.CaseLinkStore.
func (l *CaseLinkStore) List(rpc *model.SearchOptions) (*_go.CaseLinkList, error) {
	panic("unimplemented")
}

// Merge implements store.CaseLinkStore.
func (l *CaseLinkStore) Merge(req *model.CreateOptions) (*_go.CaseLinkList, error) {
	panic("unimplemented")
}

// Update implements store.CaseLinkStore.
func (l *CaseLinkStore) Update(req *model.UpdateOptions, upd *_go.InputCaseLink) (*_go.CaseLink, error) {
	panic("unimplemented")
}

func NewLinkCaseStore(store store.Store) (store.CaseLinkStore, error) {
	if store == nil {
		return nil, dberr.NewDBError("postgres.new_link_case.check.bad_arguments",
			"error creating link case interface to the comment_case table, main store is nil")
	}
	return &CaseLinkStore{storage: store, mainTable: "cases.case_link"}, nil
}

func buildLinkSelectColumnsAndPlan(base squirrel.SelectBuilder, left string, fields []string) (squirrel.SelectBuilder, []LinkScan, *dberr.DBError) {
	var plan []LinkScan

	for _, field := range fields {
		switch field {
		case "id":
			base = base.Column(ident(left, "id"))
			plan = append(plan, func(link *_go.CaseLink) any {
				return &link.Id
			})
		case "ver":
			base = base.Column(ident(left, "ver"))
			plan = append(plan, func(link *_go.CaseLink) any {
				return &link.Ver
			})
		case "created_by":
			base = base.Column(fmt.Sprintf("(SELECT ROW(id, name)::text FROM directory.wbt_user WHERE id = %s.created_by) created_by", left))
			plan = append(plan, func(link *_go.CaseLink) any {
				return scanner.ScanRowLookup(&link.CreatedBy)
			})
		case "created_at":
			base = base.Column(ident(left, "created_at"))
			plan = append(plan, func(link *_go.CaseLink) any {
				return scanner.ScanTimestamp(&link.CreatedAt)
			})
		case "updated_by":
			base = base.Column("(SELECT ROW(id, name)::text FROM directory.wbt_user WHERE id = i.updated_by) updated_by")
			plan = append(plan, func(link *_go.CaseLink) any {
				return scanner.ScanRowLookup(&link.UpdatedBy)
			})
		case "updated_at":
			base = base.Column(ident(left, "updated_at"))
			plan = append(plan, func(link *_go.CaseLink) any {
				return scanner.ScanTimestamp(&link.UpdatedAt)
			})
		case "name":
			base = base.Column(ident(left, "name"))
			plan = append(plan, func(link *_go.CaseLink) any {
				return scanner.ScanText(&link.Name)
			})
		case "url":
			base = base.Column(ident(left, "url"))
			plan = append(plan, func(link *_go.CaseLink) any {
				return &link.Url
			})
		case "author":
			base = base.Column(fmt.Sprintf(`(SELECT ROW(ct.id, ct.common_name)::text
					FROM contacts.contact ct
					WHERE id = (SELECT contact_id FROM directory.wbt_user WHERE id = %s.created_by)) author`, left))
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

func buildCreateLinkQuery(rpc *model.CreateOptions, caseId int64, add *_go.InputCaseLink) (squirrel.Sqlizer, []LinkScan, *dberr.DBError) {
	// insert
	base := squirrel.SelectBuilder{}.Prefix(`WITH i AS (INSERT INTO cases.case_link (created_by, updated_by, name, url, case_id, dc) VALUES (?, ?,
	                                                                                           NULLIF(?, ''),
	                                                                                           ?,
	                                                                                           ?, ?) RETURNING *)`,
		rpc.Session.GetUserId(), rpc.Session.GetUserId(), add.GetName(), add.GetUrl(), caseId, rpc.Session.GetDomainId())

	// select
	base = base.From("i").PlaceholderFormat(squirrel.Dollar)
	// build plan from columns
	base, plan, err := buildLinkSelectColumnsAndPlan(base, "i", rpc.Fields)
	if err != nil {
		return nil, nil, err
	}

	return base, plan, nil
}

func scanLinks(rows pgx.Rows, plan []LinkScan) ([]*_go.CaseLink, error) {
	var res []*_go.CaseLink

	for rows.Next() {
		link, err := scanLink(rows, plan)
		if err != nil {
			return nil, err
		}
		res = append(res, link)
	}
	return res, nil
}

func scanLink(row pgx.Row, plan []LinkScan) (*_go.CaseLink, error) {
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
