package app

import (
	"github.com/webitel/cases/internal/model"
	watcherkit "github.com/webitel/webitel-go-kit/pkg/watcher"
	"log"
	"log/slog"

	grpchandler "github.com/webitel/cases/internal/api_handler/grpc"

	cases "github.com/webitel/cases/api/cases"
	"google.golang.org/grpc"
)

// serviceRegistration holds information for initializing and registering a gRPC service.
type serviceRegistration struct {
	init     func(*App) (interface{}, error)                    // Initialization function for *App
	register func(grpcServer *grpc.Server, service interface{}) // Registration function for gRPC server
	name     string                                             // Service name for logging
}

// RegisterServices initializes and registers all necessary gRPC services.
func RegisterServices(grpcServer *grpc.Server, appInstance *App) {
	services := []serviceRegistration{
		{
			init: func(a *App) (interface{}, error) { return NewCaseService(a) },
			register: func(s *grpc.Server, svc interface{}) {
				cases.RegisterCasesServer(s, svc.(cases.CasesServer))
			},
			name: "Cases",
		},
		{
			init: func(a *App) (interface{}, error) {
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
			register: func(s *grpc.Server, svc interface{}) {
				cases.RegisterCaseCommentsServer(s, svc.(cases.CaseCommentsServer))
			},
			name: "CaseComments",
		},
		{
			init: func(a *App) (interface{}, error) {
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
			register: func(s *grpc.Server, svc interface{}) {
				cases.RegisterCaseLinksServer(s, svc.(cases.CaseLinksServer))
			},
			name: "CaseLinks",
		},
		{
			init: func(a *App) (interface{}, error) { return grpchandler.NewCaseTimelineService(a) },
			register: func(s *grpc.Server, svc interface{}) {
				cases.RegisterCaseTimelineServer(s, svc.(cases.CaseTimelineServer))
			},
			name: "CaseTimeline",
		},
		{
			init: func(a *App) (interface{}, error) { return NewCaseCommunicationService(a) },
			register: func(s *grpc.Server, svc interface{}) {
				cases.RegisterCaseCommunicationsServer(s, svc.(cases.CaseCommunicationsServer))
			},
			name: "CaseCommunications",
		},
		{
			init: func(a *App) (interface{}, error) { return NewRelatedCaseService(a) },
			register: func(s *grpc.Server, svc interface{}) {
				cases.RegisterRelatedCasesServer(s, svc.(cases.RelatedCasesServer))
			},
			name: "RelatedCases",
		},
		{
			init: func(a *App) (interface{}, error) {
				return grpchandler.NewCaseFileService(a)
			},
			register: func(s *grpc.Server, svc interface{}) {
				cases.RegisterCaseFilesServer(s, svc.(cases.CaseFilesServer))
			},
			name: "CaseFiles",
		},
		{
			init: func(a *App) (interface{}, error) { return grpchandler.NewSourceService(a) },
			register: func(s *grpc.Server, svc interface{}) {
				cases.RegisterSourcesServer(s, svc.(cases.SourcesServer))
			},
			name: "Sources",
		},
		{
			init: func(a *App) (interface{}, error) { return grpchandler.NewStatusService(a) },
			register: func(s *grpc.Server, svc interface{}) {
				cases.RegisterStatusesServer(s, svc.(cases.StatusesServer))
			},
			name: "Statuses",
		},
		{
			init: func(a *App) (interface{}, error) { return grpchandler.NewStatusConditionService(a) },
			register: func(s *grpc.Server, svc interface{}) {
				cases.RegisterStatusConditionsServer(s, svc.(cases.StatusConditionsServer))
			},
			name: "StatusConditions",
		},
		{
			init: func(a *App) (interface{}, error) {
				return grpchandler.NewCloseReasonService(a)
			},
			register: func(s *grpc.Server, svc interface{}) {
				cases.RegisterCloseReasonsServer(s, svc.(cases.CloseReasonsServer))
			},
			name: "CloseReasons",
		},
		{
			init: func(a *App) (interface{}, error) {
				return grpchandler.NewCloseReasonGroupsService(a)
			},
			register: func(s *grpc.Server, svc interface{}) {
				cases.RegisterCloseReasonGroupsServer(s, svc.(cases.CloseReasonGroupsServer))
			},
			name: "CloseReasonGroups",
		},
		{
			init: func(a *App) (interface{}, error) {
				return grpchandler.NewPriorityService(a)
			},
			register: func(s *grpc.Server, svc interface{}) {
				cases.RegisterPrioritiesServer(s, svc.(cases.PrioritiesServer))
			},
			name: "Priorities",
		},
		{
			init: func(a *App) (interface{}, error) {
				return grpchandler.NewSLAService(a)
			},
			register: func(s *grpc.Server, svc interface{}) {
				cases.RegisterSLAsServer(s, svc.(cases.SLAsServer))
			},
			name: "SLAs",
		},
		{
			init: func(a *App) (interface{}, error) { return grpchandler.NewSLAConditionService(a) },
			register: func(s *grpc.Server, svc interface{}) {
				cases.RegisterSLAConditionsServer(s, svc.(cases.SLAConditionsServer))
			},
			name: "SLAConditions",
		},
		{
			init: func(a *App) (interface{}, error) { return NewCatalogService(a) },
			register: func(s *grpc.Server, svc interface{}) {
				cases.RegisterCatalogsServer(s, svc.(cases.CatalogsServer))
			},
			name: "Catalogs",
		},
		{
			init: func(a *App) (interface{}, error) { return grpchandler.NewServiceService(a) },
			register: func(s *grpc.Server, svc interface{}) {
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
