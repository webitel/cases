package postgres

import (
	"context"

	_go "buf.build/gen/go/webitel/cases/protocolbuffers/go"
	dberr "github.com/webitel/cases/internal/error"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/model"
)

type CommentCase struct {
	storage store.Store
}

// Create implements store.CommentCaseStore.
func (c *CommentCase) Create(ctx context.Context, rpc *model.CreateOptions, add *_go.CaseComment) (*_go.CaseComment, error) {
	panic("unimplemented")
}

// Delete implements store.CommentCaseStore.
func (c *CommentCase) Delete(ctx context.Context, req *model.DeleteOptions) (*_go.CaseComment, error) {
	panic("unimplemented")
}

// List implements store.CommentCaseStore.
func (c *CommentCase) List(ctx context.Context, rpc *model.SearchOptions) (*_go.CaseCommentList, error) {
	panic("unimplemented")
}

// Merge implements store.CommentCaseStore.
func (c *CommentCase) Merge(ctx context.Context, req *model.CreateOptions) (*_go.CaseCommentList, error) {
	panic("unimplemented")
}

// Update implements store.CommentCaseStore.
func (c *CommentCase) Update(ctx context.Context, req *model.UpdateOptions) (*_go.CaseComment, error) {
	panic("unimplemented")
}

func NewCommentCaseStore(store store.Store) (store.CommentCaseStore, error) {
	if store == nil {
		return nil, dberr.NewDBError("postgres.new_comment_case.check.bad_arguments",
			"error creating comment case interface to the comment_case table, main store is nil")
	}
	return &CommentCase{storage: store}, nil
}
