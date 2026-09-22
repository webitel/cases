package app

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/webitel/cases/auth"
	"github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
	wlogger "github.com/webitel/webitel-go-kit/infra/logger_client"
)

func (a *App) LinkCaseArticle(creator options.Creator, input *model.CaseArticle) (*model.CaseArticle, error) {
	caseID := creator.GetParentID()
	if caseID == 0 {
		return nil, errors.InvalidArgument("case id required")
	}
	if input == nil || input.ArticleID <= 0 {
		return nil, errors.InvalidArgument("article id required")
	}
	if err := a.checkCaseRbac(creator, creator.GetAuthOpts(), auth.Edit, caseID); err != nil {
		return nil, err
	}
	link, err := a.Store.CaseArticle().Link(creator, input)
	if err != nil {
		return nil, err
	}

	a.logCaseArticle(creator, creator.GetAuthOpts(), caseID, link)
	return link, nil
}

func (a *App) UnlinkCaseArticle(deleter options.Deleter) (*model.CaseArticle, error) {
	caseID := deleter.GetParentID()
	if caseID == 0 {
		return nil, errors.InvalidArgument("case id required")
	}
	if len(deleter.GetIDs()) == 0 {
		return nil, errors.InvalidArgument("article id required")
	}
	if err := a.checkCaseRbac(deleter, deleter.GetAuthOpts(), auth.Edit, caseID); err != nil {
		return nil, err
	}
	link, err := a.Store.CaseArticle().Unlink(deleter)
	if err != nil {
		return nil, err
	}
	// The store leaves the close article untouched.
	if link.Source == model.CaseArticleSourceResolution {
		return nil, errors.Aborted("the close article is changed on the case itself")
	}

	a.logCaseArticle(deleter, deleter.GetAuthOpts(), caseID, link)
	return link, nil
}

// ListCaseArticles lists the links of a case, or the readable cases of an article.
func (a *App) ListCaseArticles(searcher options.Searcher) ([]*model.CaseArticle, error) {
	filters := searcher.GetFilter("case_id")
	if len(filters) == 0 {
		if len(searcher.GetFilter("article_id")) == 0 {
			return nil, errors.InvalidArgument("case id or article id required")
		}
		return a.Store.CaseArticle().List(searcher)
	}
	caseID, err := strconv.ParseInt(filters[0].Value, 10, 64)
	if err != nil {
		return nil, errors.InvalidArgument("invalid case id", errors.WithCause(err))
	}
	if err := a.checkCaseRbac(searcher, searcher.GetAuthOpts(), auth.Read, caseID); err != nil {
		return nil, err
	}
	return a.Store.CaseArticle().List(searcher)
}

// checkCaseRbac checks the access mode of the session to a case.
func (a *App) checkCaseRbac(ctx context.Context, session auth.Auther, mode auth.AccessMode, caseID int64) error {
	if !session.IsRbacCheckRequired(model.ScopeCases, mode) {
		return nil
	}
	access, err := a.Store.Case().CheckRbacAccess(ctx, session, mode, caseID)
	if err != nil {
		return err
	}
	if !access {
		return errors.Forbidden("user doesn't have required access to the case")
	}
	return nil
}

// logCaseArticle writes the link change to the case audit log; a failure is only logged.
func (a *App) logCaseArticle(ctx context.Context, session auth.Auther, caseID int64, link *model.CaseArticle) {
	message, err := wlogger.NewMessage(
		session.GetUserId(),
		session.GetUserIp(),
		wlogger.UpdateAction,
		strconv.FormatInt(caseID, 10),
		link,
	)
	if err == nil {
		_, err = a.wtelLogger.SendContext(context.Background(), session.GetDomainId(), model.ScopeCases, message)
	}
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
	}
}
