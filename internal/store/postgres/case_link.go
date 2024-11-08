package postgres

import (
	"github.com/jackc/pgx/v5"
	_go "github.com/webitel/cases/api/cases"
	dberr "github.com/webitel/cases/internal/error"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/model"
)

type LinkScan func(link *_go.CaseLink) store.ScanFunc

var linkSelectFieldsMap = map[string]string{
	"id":         "case_link.id",
	"ver":        "case_link.ver",
	"dc":         "case_link.dc",
	"created_at": "case_link.created_at",
	"created_by": "case_link.created_by",
	"updated_at": "case_link.updated_at",
	"updated_by": "case_link.updated_by",
	"name":       "case_link.name",
	"url":        "case_link.url",
	"case_id":    "case_link.case_id",
}

func getLinkSelectFields() []string {
	var fields []string
	for selectField, _ := range linkSelectFieldsMap {
		fields = append(fields, selectField)
	}
	return fields
}

type CaseLinkStore struct {
	storage   store.Store
	mainTable string
}

// Create implements store.LinkCaseStore.
func (l *CaseLinkStore) Create(rpc *model.CreateOptions, caseId int64, add *_go.InputCaseLink) (*_go.CaseLink, error) {
	if rpc == nil {
		return nil, dberr.NewDBError("postgres.case_link.create.check_args.opts", "create options required")
	}
	//db, err := l.storage.Database()
	//if err != nil {
	//	return nil, err
	//}
	//var res _go.CaseLink
	//row := db.QueryRow(rpc.Context,
	//	`WITH i AS (INSERT INTO cases.case_link (created_by, updated_by, name, url, case_id, dc) VALUES (?, ?,
	//                                                                                            ?,
	//                                                                                            ?,
	//                                                                                            ?, ?) RETURNING *)
	//		SELECT i.id,
	//			   i.dc,
	//			   (SELECT ROW (id, name) FROM directory.wbt_user WHERE id = i.created_by),
	//			   i.created_at,
	//			   (SELECT ROW (id, name) FROM directory.wbt_user WHERE id = i.updated_by),
	//			   i.updated_at,
	//			   i.name,
	//			   i.url,
	//			   (SELECT ROW (ct.id, ct.common_name)
	//				FROM contacts.contact ct
	//				WHERE id = (SELECT contact_id FROM directory.wbt_user WHERE id = i.created_by)) author
	//		FROM i`,
	//	rpc.Session.GetUserId(), rpc.Session.GetUserId(), add.GetName(), add.GetUrl(), caseId, rpc.Session.GetDomainId(),
	//)
	//plan := []LinkScan{
	//	func(link *_go.CaseLink) store.ScanFunc {
	//		return func(src any) error {
	//
	//		}
	//	},
	//}
	return nil, nil
}

// Delete implements store.LinkCaseStore.
func (l *CaseLinkStore) Delete(req *model.DeleteOptions) (*_go.CaseLink, error) {
	panic("unimplemented")
}

// List implements store.LinkCaseStore.
func (l *CaseLinkStore) List(rpc *model.SearchOptions) (*_go.CaseLinkList, error) {
	panic("unimplemented")
}

// Merge implements store.LinkCaseStore.
func (l *CaseLinkStore) Merge(req *model.CreateOptions) (*_go.CaseLinkList, error) {
	panic("unimplemented")
}

// Update implements store.LinkCaseStore.
func (l *CaseLinkStore) Update(req *model.UpdateOptions, upd *_go.InputCaseLink) (*_go.CaseLink, error) {
	panic("unimplemented")
}

func NewLinkCaseStore(store store.Store) (store.LinkCaseStore, error) {
	if store == nil {
		return nil, dberr.NewDBError("postgres.new_link_case.check.bad_arguments",
			"error creating link case interface to the comment_case table, main store is nil")
	}
	return &CaseLinkStore{storage: store, mainTable: "cases.case_link"}, nil
}

func LinksScan(rows pgx.Rows) ([]*_go.CaseLink, error) {
	//var result []*_go.CaseLink
	//for rows.Next() {
	//	err := rows.Err()
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//}
	return nil, nil
}
