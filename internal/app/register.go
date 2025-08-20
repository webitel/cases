package app

import (
	"log"
	"log/slog"

	"google.golang.org/grpc"

	watcherkit "github.com/webitel/webitel-go-kit/pkg/watcher"

	"github.com/webitel/cases/api/cases"
	grpchandler "github.com/webitel/cases/internal/api_handler/grpc"
	"github.com/webitel/cases/internal/model"
)

// serviceRegistration holds information for initializing and registering a gRPC service.
type serviceRegistration struct {
	init     func(*App) (any, error)                    // Initialization function for *App
	register func(grpcServer *grpc.Server, service any) // Registration function for gRPC server
	name     string                                     // Service name for logging
}

// RegisterServices initializes and registers all necessary gRPC services.
func RegisterServices(grpcServer *grpc.Server, appInstance *App) {
	services := []serviceRegistration{
		{
			init: func(a *App) (any, error) { return NewCaseService(a) },
			register: func(s *grpc.Server, svc any) {
				cases.RegisterCasesServer(s, svc.(cases.CasesServer))
			},
			name: "Cases",
		},
		{
			init: func(a *App) (any, error) {
				// Initialize watchers first
				if a.config.TriggerWatcher.Enabled {
					watcher := watcherkit.NewDefaultWatcher()

					// Add logger observer if enabled
					if a.config.LoggerWatcher.Enabled {
						obs, err := NewLoggerObserver(a.wtelLogger, caseCommentsObjScope, defaultLogTimeout)
						if err != nil {
							return nil, err
						}
						watcher.Attach(watcherkit.EventTypeCreate, obs)
						watcher.Attach(watcherkit.EventTypeUpdate, obs)
						watcher.Attach(watcherkit.EventTypeDelete, obs)
					}

					// Add FTS observer if enabled
					if a.config.FtsWatcher.Enabled {
						ftsObserver, err := NewFullTextSearchObserver(a.ftsClient, caseCommentsObjScope, formCommentsFtsModel)
						if err != nil {
							return nil, err
						}
						watcher.Attach(watcherkit.EventTypeCreate, ftsObserver)
						watcher.Attach(watcherkit.EventTypeUpdate, ftsObserver)
						watcher.Attach(watcherkit.EventTypeDelete, ftsObserver)
					}

					// Add trigger observer
					mq, err := NewTriggerObserver(a.rabbitPublisher, a.config.TriggerWatcher, formCaseCommentTriggerModel, slog.With(
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

					// Register the watcher
					a.watcherManager.AddWatcher(caseCommentsObjScope, watcher)
				}

				// Then create gRPC service using App directly
				return grpchandler.NewCaseCommentService(a)
			},
			register: func(s *grpc.Server, svc any) {
				cases.RegisterCaseCommentsServer(s, svc.(cases.CaseCommentsServer))
			},
			name: "CaseComments",
		},
		{
			init: func(a *App) (any, error) {
				if a.config.TriggerWatcher.Enabled {
					watcher := watcherkit.NewDefaultWatcher()
					mq, err := NewTriggerObserver(a.rabbitPublisher, a.config.TriggerWatcher, formCaseLinkTriggerModel, slog.With(
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

					if a.caseResolutionTimer != nil {
						a.caseResolutionTimer.Start()
					}
					a.watcherManager.AddWatcher(model.BrokerScopeCaseLinks, watcher)
				}
				return grpchandler.NewCaseLinkService(a), nil
			},
			register: func(s *grpc.Server, svc any) {
				cases.RegisterCaseLinksServer(s, svc.(cases.CaseLinksServer))
			},
			name: "CaseLinks",
		},
		{
			init: func(a *App) (any, error) { return grpchandler.NewCaseTimelineService(a) },
			register: func(s *grpc.Server, svc any) {
				cases.RegisterCaseTimelineServer(s, svc.(cases.CaseTimelineServer))
			},
			name: "CaseTimeline",
		},
		{
			init: func(a *App) (any, error) { return grpchandler.NewCaseCommunicationService(a) },
			register: func(s *grpc.Server, svc any) {
				cases.RegisterCaseCommunicationsServer(s, svc.(cases.CaseCommunicationsServer))
			},
			name: "CaseCommunications",
		},
		{
			init: func(a *App) (any, error) { return NewRelatedCaseService(a) },
			register: func(s *grpc.Server, svc any) {
				cases.RegisterRelatedCasesServer(s, svc.(cases.RelatedCasesServer))
			},
			name: "RelatedCases",
		},
		{
			init: func(a *App) (any, error) {
				return grpchandler.NewCaseFileService(a)
			},
			register: func(s *grpc.Server, svc any) {
				cases.RegisterCaseFilesServer(s, svc.(cases.CaseFilesServer))
			},
			name: "CaseFiles",
		},
		{
			init: func(a *App) (any, error) { return grpchandler.NewSourceService(a) },
			register: func(s *grpc.Server, svc any) {
				cases.RegisterSourcesServer(s, svc.(cases.SourcesServer))
			},
			name: "Sources",
		},
		{
			init: func(a *App) (any, error) { return grpchandler.NewStatusService(a) },
			register: func(s *grpc.Server, svc any) {
				cases.RegisterStatusesServer(s, svc.(cases.StatusesServer))
			},
			name: "Statuses",
		},
		{
			init: func(a *App) (any, error) { return grpchandler.NewStatusConditionService(a) },
			register: func(s *grpc.Server, svc any) {
				cases.RegisterStatusConditionsServer(s, svc.(cases.StatusConditionsServer))
			},
			name: "StatusConditions",
		},
		{
			init: func(a *App) (any, error) {
				return grpchandler.NewCloseReasonService(a)
			},
			register: func(s *grpc.Server, svc any) {
				cases.RegisterCloseReasonsServer(s, svc.(cases.CloseReasonsServer))
			},
			name: "CloseReasons",
		},
		{
			init: func(a *App) (any, error) {
				return grpchandler.NewCloseReasonGroupsService(a)
			},
			register: func(s *grpc.Server, svc any) {
				cases.RegisterCloseReasonGroupsServer(s, svc.(cases.CloseReasonGroupsServer))
			},
			name: "CloseReasonGroups",
		},
		{
			init: func(a *App) (any, error) {
				return grpchandler.NewPriorityService(a)
			},
			register: func(s *grpc.Server, svc any) {
				cases.RegisterPrioritiesServer(s, svc.(cases.PrioritiesServer))
			},
			name: "Priorities",
		},
		{
			init: func(a *App) (any, error) {
				return grpchandler.NewSLAService(a)
			},
			register: func(s *grpc.Server, svc any) {
				cases.RegisterSLAsServer(s, svc.(cases.SLAsServer))
			},
			name: "SLAs",
		},
		{
			init: func(a *App) (any, error) { return grpchandler.NewSLAConditionService(a) },
			register: func(s *grpc.Server, svc any) {
				cases.RegisterSLAConditionsServer(s, svc.(cases.SLAConditionsServer))
			},
			name: "SLAConditions",
		},
		{
			init: func(a *App) (any, error) { return NewCatalogService(a) },
			register: func(s *grpc.Server, svc any) {
				cases.RegisterCatalogsServer(s, svc.(cases.CatalogsServer))
			},
			name: "Catalogs",
		},
		{
			init: func(a *App) (any, error) { return grpchandler.NewServiceService(a) },
			register: func(s *grpc.Server, svc any) {
				cases.RegisterServicesServer(s, svc.(cases.ServicesServer))
			},
			name: "Services",
		},
	}

	// Initialize and register each service
	for _, service := range services {
		svc, err := service.init(appInstance)
		if err != nil {
			log.Printf("Error initializing %s service: %v", service.name, err)

			continue
		}
		service.register(grpcServer, svc)
		log.Printf("%s service registered successfully", service.name)
	}
}
