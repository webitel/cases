package postgres

import (
	"context"
	"time"

	_go "buf.build/gen/go/webitel/cases/protocolbuffers/go"
	general "buf.build/gen/go/webitel/general/protocolbuffers/go"

	dberr "github.com/webitel/cases/internal/error"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/model"
	util "github.com/webitel/cases/util"
)

type CaseComment struct {
	storage store.Store
}

// Create implements store.CommentCaseStore.
func (c *CaseComment) Create(
	ctx context.Context,
	rpc *model.CreateOptions,
	add *_go.CaseComment,
) (*_go.CaseComment, error) {
	// Establish database connection
	d, dbErr := c.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.comment.create.database_connection_error", dbErr)
	}

	// Build query and arguments
	query, args, err := c.buildCreateCommentQuery(rpc, add)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.comment.create.query_build_error", err)
	}

	var (
		createdByLookup, updatedByLookup general.Lookup
		createdAt, updatedAt             time.Time
	)

	// Execute query and scan result
	err = d.QueryRow(ctx, query, args...).Scan(
		&add.Id, &createdAt, &add.Text,
		&createdByLookup.Id, &createdByLookup.Name,
		&updatedAt, &updatedByLookup.Id, &updatedByLookup.Name,
	)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.comment.create.execution_error", err)
	}

	// Populate and return the created comment
	return &_go.CaseComment{
		Id:        add.Id,
		Text:      add.Text,
		CreatedAt: util.Timestamp(createdAt),
		UpdatedAt: util.Timestamp(updatedAt),
		CreatedBy: &createdByLookup,
		UpdatedBy: &updatedByLookup,
	}, nil
}

func (c *CaseComment) buildCreateCommentQuery(rpc *model.CreateOptions, comment *_go.CaseComment) (string, []interface{}, error) {
	args := []interface{}{
		rpc.Session.GetDomainId(), // dc
		rpc.Ids[1],                // case_id
		rpc.CurrentTime(),         // created_at (and updated_at)
		rpc.Session.GetUserId(),   // created_by (and updated_by)
		comment.Text,              // comment text
	}

	return createCommentQuery, args, nil
}

// Delete implements store.CommentCaseStore.
func (c *CaseComment) Delete(
	ctx context.Context,
	req *model.DeleteOptions,
) (*_go.CaseComment, error) {
	panic("unimplemented")
}

// List implements store.CommentCaseStore.
func (c *CaseComment) List(
	ctx context.Context,
	rpc *model.SearchOptions,
) (*_go.CaseCommentList, error) {
	panic("unimplemented")
}

// Merge implements store.CommentCaseStore.
func (c *CaseComment) Merge(
	ctx context.Context,
	req *model.CreateOptions,
	comments []*_go.CaseComment,
) (*_go.CaseCommentList, error) {
	panic("unimplemented")
}

// Update implements store.CommentCaseStore.
func (c *CaseComment) Update(
	ctx context.Context,
	req *model.UpdateOptions,
	upd *_go.CaseComment,
) (*_go.CaseComment, error) {
	panic("unimplemented")
}

var createCommentQuery = store.CompactSQL(`WITH ins AS (
    INSERT INTO cases.case_comment (dc, case_id, created_at, created_by, updated_at, updated_by, comment)
        VALUES ($1, $2, $3, $4, $3, $4, NULLIF($5, ''))
        RETURNING id, dc, case_id, created_at, created_by, updated_at, updated_by, comment)
SELECT ins.id,
       ins.created_at,
       ins.comment,
       ins.created_by                     AS created_by_id,
       COALESCE(c.name::text, c.username) AS created_by_name,
       ins.updated_at,
       ins.updated_by                     AS updated_by_id,
       COALESCE(u.name::text, u.username) AS updated_by_name
FROM ins
         LEFT JOIN directory.wbt_user u ON u.id = ins.updated_by
         LEFT JOIN directory.wbt_user c ON c.id = ins.created_by;`)

func NewCaseCommentStore(store store.Store) (store.CaseCommentStore, error) {
	if store == nil {
		return nil, dberr.NewDBError("postgres.new_case_comment.check.bad_arguments",
			"error creating comment case interface to the case_comment table, main store is nil")
	}
	return &CaseComment{storage: store}, nil
}
