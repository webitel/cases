package app

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/webitel/cases/auth"
	auth_util "github.com/webitel/cases/auth/util"
	"github.com/webitel/cases/internal/api_handler/grpc"
	"github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
	wlogger "github.com/webitel/webitel-go-kit/infra/logger_client"
	watcherkit "github.com/webitel/webitel-go-kit/pkg/watcher"
)

// In search options extract from context user
// Remove from search options fields functions

type CaseLinkApp struct {
	app    *App
	logger *wlogger.ObjectedLogger
}

func (a *CaseLinkApp) CreateCaseLink(creator options.Creator, input *model.CaseLink) (*model.CaseLink, error) {
	caseID := creator.GetParentID()
	if caseID == 0 {
		return nil, errors.InvalidArgument("case id required")
	}
	if input == nil || input.Url == "" {
		return nil, errors.InvalidArgument("url is required for each link")
	}
	accessMode := auth.Edit
	if creator.GetAuthOpts().IsRbacCheckRequired(model.ScopeCases, accessMode) {
		access, err := a.app.Store.Case().CheckRbacAccess(creator, creator.GetAuthOpts(), accessMode, caseID)
		if err != nil {
			return nil, err
		}
		if !access {
			return nil, errors.Forbidden("user doesn't have required (EDIT) access to the case")
		}
	}
	link, err := a.app.Store.CaseLink().Create(creator, input)
	if err != nil {
		return nil, err
	}

	authOpts := creator.GetAuthOpts()
	if input.Author != nil && input.Author.Id != nil {
		authOpts = auth_util.CloneWithUserID(authOpts, int64(*input.Author.Id))
	}

	message, err := wlogger.NewMessage(
		authOpts.GetUserId(),
		authOpts.GetUserIp(),
		wlogger.CreateAction,
		strconv.FormatInt(link.Id, 10),
		input,
	)
	if err == nil && a.logger != nil {
		_, _ = a.logger.SendContext(context.Background(), authOpts.GetDomainId(), message)
	}

	if a.app.watcherManager != nil {
		if notifyErr := a.app.watcherManager.Notify(
			model.BrokerScopeCaseLinks,
			watcherkit.EventTypeCreate,
			NewLinkWatcherData(authOpts, link, link.Id, authOpts.GetDomainId()),
		); notifyErr != nil {
			slog.ErrorContext(context.Background(), fmt.Sprintf("could not notify link create: %s", notifyErr.Error()))
		}
	}

	return link, nil
}

func (a *CaseLinkApp) UpdateCaseLink(updator options.Updator, input *model.CaseLink) (*model.CaseLink, error) {
	linkIDs := updator.GetEtags()
	if len(linkIDs) == 0 {
		return nil, errors.InvalidArgument("link id required")
	}
	caseID := updator.GetParentID()
	if caseID == 0 {
		return nil, errors.InvalidArgument("case id required")
	}
	accessMode := auth.Edit
	if updator.GetAuthOpts().IsRbacCheckRequired(grpc.CaseLinkMetadata.GetParentScopeName(), accessMode) {
		access, err := a.app.Store.Case().CheckRbacAccess(updator, updator.GetAuthOpts(), accessMode, caseID)
		if err != nil {
			return nil, err
		}
		if !access {
			return nil, errors.Forbidden("user doesn't have required (EDIT) access to the case")
		}
	}
	link, err := a.app.Store.CaseLink().Update(updator, input)
	if err != nil {
		slog.ErrorContext(context.Background(), err.Error())
		return nil, err
	}

	authOpts := updator.GetAuthOpts()
	if input.Author != nil && input.Author.Id != nil {
		authOpts = auth_util.CloneWithUserID(authOpts, int64(*input.Author.Id))
	}

	message, err := wlogger.NewMessage(
		authOpts.GetUserId(),
		authOpts.GetUserIp(),
		wlogger.UpdateAction,
		strconv.FormatInt(link.Id, 10),
		input,
	)
	if err == nil && a.logger != nil {
		_, _ = a.logger.SendContext(context.Background(), authOpts.GetDomainId(), message)
	}

	if a.app.watcherManager != nil {
		if notifyErr := a.app.watcherManager.Notify(
			model.BrokerScopeCaseLinks,
			watcherkit.EventTypeUpdate,
			NewLinkWatcherData(authOpts, link, link.Id, authOpts.GetDomainId()),
		); notifyErr != nil {
			slog.ErrorContext(context.Background(), fmt.Sprintf("could not notify link update: %s", notifyErr.Error()))
		}
	}

	return link, nil
}

func (a *CaseLinkApp) DeleteCaseLink(deleter options.Deleter) (*model.CaseLink, error) {
	linkIDs := deleter.GetIDs()
	if len(linkIDs) == 0 {
		return nil, errors.InvalidArgument("link id required")
	}
	caseID := deleter.GetParentID()
	if caseID == 0 {
		return nil, errors.InvalidArgument("case id required")
	}
	accessMode := auth.Edit
	if deleter.GetAuthOpts().IsRbacCheckRequired(grpc.CaseLinkMetadata.GetParentScopeName(), accessMode) {
		access, err := a.app.Store.Case().CheckRbacAccess(deleter, deleter.GetAuthOpts(), accessMode, caseID)
		if err != nil {
			return nil, err
		}
		if !access {
			return nil, errors.Forbidden("user doesn't have required (Edit) access to the case")
		}
	}
	link, err := a.app.Store.CaseLink().Delete(deleter)
	if err != nil {
		return nil, err
	}

	authOpts := deleter.GetAuthOpts()
	message, err := wlogger.NewMessage(
		authOpts.GetUserId(),
		authOpts.GetUserIp(),
		wlogger.UpdateAction,
		strconv.FormatInt(linkIDs[0], 10),
		link,
	)
	if err == nil && a.logger != nil {
		_, _ = a.logger.SendContext(context.Background(), authOpts.GetDomainId(), message)
	}

	if a.app.watcherManager != nil {
		if notifyErr := a.app.watcherManager.Notify(
			model.BrokerScopeCaseLinks,
			watcherkit.EventTypeDelete,
			NewLinkWatcherData(authOpts, link, linkIDs[0], authOpts.GetDomainId()),
		); notifyErr != nil {
			slog.ErrorContext(context.Background(), fmt.Sprintf("could not notify link delete: %s", notifyErr.Error()))
		}
	}

	return link, nil
}

func (a *CaseLinkApp) ListCaseLinks(searcher options.Searcher) ([]*model.CaseLink, error) {
	caseID, ok := searcher.GetFilter("case_id").(int64)
	if !ok || caseID == 0 {
		return nil, errors.InvalidArgument("case id required")
	}
	accessMode := auth.Read
	if searcher.GetAuthOpts().IsRbacCheckRequired(grpc.CaseLinkMetadata.GetParentScopeName(), accessMode) {
		access, err := a.app.Store.Case().CheckRbacAccess(searcher, searcher.GetAuthOpts(), accessMode, caseID)
		if err != nil {
			return nil, err
		}
		if !access {
			return nil, errors.Forbidden("user doesn't have required (READ) access to the case")
		}
	}
	links, err := a.app.Store.CaseLink().List(searcher)
	if err != nil {
		return nil, err
	}
	return links, nil
}

func NewCaseLinkApp(app *App) (*CaseLinkApp, error) {
	if app == nil {
		return nil, errors.New(
			"unable to init service, app is nil",
		)
	}
	logger, err := app.wtelLogger.GetObjectedLogger("cases")
	if err != nil {
		return nil, err
	}
	service := &CaseLinkApp{
		app:    app,
		logger: logger,
	}
	watcher := watcherkit.NewDefaultWatcher()

	if app.config.TriggerWatcher.Enabled {
		mq, err := NewTriggerObserver(app.rabbitPublisher, app.config.TriggerWatcher, formCaseLinkTriggerModel, slog.With(
			slog.Group("context",
				slog.String("scope", "watcher")),
		))

		if err != nil {
			return nil, err
		}
		watcher.Attach(watcherkit.EventTypeCreate, mq)
		watcher.Attach(watcherkit.EventTypeUpdate, mq)
		watcher.Attach(watcherkit.EventTypeDelete, mq)
		watcher.Attach(watcherkit.EventTypeResolutionTime, mq)

		if app.caseResolutionTimer != nil {
			app.caseResolutionTimer.Start()
		}
	}

	app.watcherManager.AddWatcher(model.BrokerScopeCaseLinks, watcher)

	return service, nil
}

type CaseLinkWatcherData struct {
	link *model.CaseLink
	Args map[string]any
}

func (wd *CaseLinkWatcherData) GetArgs() map[string]any {
	return wd.Args
}

func NewLinkWatcherData(session auth.Auther, link *model.CaseLink, linkId int64, dc int64) *CaseLinkWatcherData {
	return &CaseLinkWatcherData{
		link: link,
		Args: map[string]any{
			"session":   session,
			"obj":       link,
			"id":        linkId,
			"domain_id": dc,
		},
	}
}
