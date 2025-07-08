package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/webitel/cases/auth"
	"github.com/webitel/cases/internal/api_handler/grpc"
	"github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
	watcherkit "github.com/webitel/webitel-go-kit/pkg/watcher"
	"google.golang.org/grpc/codes"
)

const caseCommentsObjScope = model.ScopeCaseComments

// ListCaseComments lists case comments with filters and pagination.
func (s *App) ListCaseComments(searcher options.Searcher) ([]*model.CaseComment, error) {
	comments, err := s.Store.CaseComment().List(searcher)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

// UpdateCaseComment updates a case comment in the store.
func (s *App) UpdateCaseComment(updator options.Updator, input *model.CaseComment) (*model.CaseComment, error) {
	if input.Text == "" {
		return nil, errors.InvalidArgument("Text is required")
	}

	updatedComment, err := s.Store.CaseComment().Update(updator, input)
	if err != nil {
		return nil, err
	}

	if notifyErr := s.watcherManager.Notify(
		caseCommentsObjScope,
		watcherkit.EventTypeUpdate,
		NewCaseCommentWatcherData(updator.GetAuthOpts(), updatedComment, updatedComment.Id, updatedComment.CaseId, updatedComment.RoleIds),
	); notifyErr != nil {
		slog.ErrorContext(context.Background(), fmt.Sprintf("could not notify comment update: %s", notifyErr.Error()))
	}

	return updatedComment, nil
}

// DeleteCaseComment deletes a case comment from the store.
func (s *App) DeleteCaseComment(deleter options.Deleter) (*model.CaseComment, error) {
	deletedComment, err := s.Store.CaseComment().Delete(deleter)
	if err != nil {
		return nil, err
	}

	if notifyErr := s.watcherManager.Notify(
		caseCommentsObjScope,
		watcherkit.EventTypeDelete,
		NewCaseCommentWatcherData(deleter.GetAuthOpts(), deletedComment, deletedComment.Id, deletedComment.CaseId, deletedComment.RoleIds),
	); notifyErr != nil {
		slog.ErrorContext(context.Background(), fmt.Sprintf("could not notify comment delete: %s", notifyErr.Error()))
	}

	return deletedComment, nil
}

// PublishCaseComment creates a new case comment.
func (s *App) PublishCaseComment(creator options.Creator, input *model.CaseComment) (*model.CaseComment, error) {
	if input.Text == "" {
		return nil, errors.InvalidArgument("text is required")
	}

	accessMode := auth.Read
	if !creator.GetAuthOpts().CheckObacAccess(grpc.CaseCommentMetadata.GetParentScopeName(), accessMode) {
		return nil, errors.New("user doesn't have required (EDIT) access to the case", errors.WithCode(codes.PermissionDenied))
	}
	if creator.GetAuthOpts().IsRbacCheckRequired(grpc.CaseCommentMetadata.GetParentScopeName(), accessMode) {
		access, err := s.Store.Case().CheckRbacAccess(context.Background(), creator.GetAuthOpts(), accessMode, creator.GetParentID())
		if err != nil {
			return nil, err
		}
		if !access {
			return nil, errors.New("user doesn't have required (EDIT) access to the case", errors.WithCode(codes.PermissionDenied))
		}
	}

	comment, err := s.Store.CaseComment().Publish(creator, input)
	if err != nil {
		return nil, err
	}

	if notifyErr := s.watcherManager.Notify(
		caseCommentsObjScope,
		watcherkit.EventTypeCreate,
		NewCaseCommentWatcherData(creator.GetAuthOpts(), comment, comment.Id, comment.CaseId, comment.RoleIds),
	); notifyErr != nil {
		slog.ErrorContext(context.Background(), fmt.Sprintf("could not notify comment create: %s", notifyErr.Error()))
	}

	return comment, nil
}

func formCommentsFtsModel(comment *model.CaseComment, params map[string]any) (*model.FtsCaseComment, error) {
	roles, ok := params["role_ids"].([]int64)
	if !ok {
		return nil, fmt.Errorf("role ids required for FTS model")
	}
	caseId, ok := params["case_id"].(int64)
	if !ok {
		return nil, fmt.Errorf("case id required for FTS model")
	}

	return &model.FtsCaseComment{
		ParentId:  caseId,
		Comment:   comment.Text,
		RoleIds:   roles,
		CreatedAt: comment.CreatedAt.Unix() * 1000, // Convert to milliseconds
	}, nil
}

type CaseCommentWatcherData struct {
	comment *model.CaseComment
	Args    map[string]any
}

func NewCaseCommentWatcherData(session auth.Auther, comment *model.CaseComment, id, caseId int64, roleIds []int64) *CaseCommentWatcherData {
	return &CaseCommentWatcherData{comment: comment, Args: map[string]any{"session": session, "obj": comment, "case_id": caseId, "role_ids": roleIds, "id": id}}
}

func (wd *CaseCommentWatcherData) Marshal() ([]byte, error) {
	return json.Marshal(wd.comment)
}

func (wd *CaseCommentWatcherData) GetArgs() map[string]any {
	return wd.Args
}
