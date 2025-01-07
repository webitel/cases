package app

import (
	"log"

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
			init: func(a *App) (interface{}, error) { return NewCaseCommentService(a) },
			register: func(s *grpc.Server, svc interface{}) {
				cases.RegisterCaseCommentsServer(s, svc.(cases.CaseCommentsServer))
			},
			name: "CaseComments",
		},
		{
			init: func(a *App) (interface{}, error) { return NewCaseLinkService(a) },
			register: func(s *grpc.Server, svc interface{}) {
				cases.RegisterCaseLinksServer(s, svc.(cases.CaseLinksServer))
			},
			name: "CaseLinks",
		},
		{
			init: func(a *App) (interface{}, error) { return NewCaseTimelineService(a) },
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
			init: func(a *App) (interface{}, error) { return NewCaseFileService(a) },
			register: func(s *grpc.Server, svc interface{}) {
				cases.RegisterCaseFilesServer(s, svc.(cases.CaseFilesServer))
			},
			name: "CaseFiles",
		},
		{
			init: func(a *App) (interface{}, error) { return NewSourceService(a) },
			register: func(s *grpc.Server, svc interface{}) {
				cases.RegisterSourcesServer(s, svc.(cases.SourcesServer))
			},
			name: "Sources",
		},
		{
			init: func(a *App) (interface{}, error) { return NewStatusService(a) },
			register: func(s *grpc.Server, svc interface{}) {
				cases.RegisterStatusesServer(s, svc.(cases.StatusesServer))
			},
			name: "Statuses",
		},
		{
			init: func(a *App) (interface{}, error) { return NewStatusConditionService(a) },
			register: func(s *grpc.Server, svc interface{}) {
				cases.RegisterStatusConditionsServer(s, svc.(cases.StatusConditionsServer))
			},
			name: "StatusConditions",
		},
		{
			init: func(a *App) (interface{}, error) { return NewCloseReasonService(a) },
			register: func(s *grpc.Server, svc interface{}) {
				cases.RegisterCloseReasonsServer(s, svc.(cases.CloseReasonsServer))
			},
			name: "CloseReasons",
		},
		{
			init: func(a *App) (interface{}, error) { return NewCloseReasonGroupsService(a) },
			register: func(s *grpc.Server, svc interface{}) {
				cases.RegisterCloseReasonGroupsServer(s, svc.(cases.CloseReasonGroupsServer))
			},
			name: "CloseReasonGroups",
		},
		{
			init: func(a *App) (interface{}, error) { return NewPriorityService(a) },
			register: func(s *grpc.Server, svc interface{}) {
				cases.RegisterPrioritiesServer(s, svc.(cases.PrioritiesServer))
			},
			name: "Priorities",
		},
		{
			init: func(a *App) (interface{}, error) { return NewSLAService(a) },
			register: func(s *grpc.Server, svc interface{}) {
				cases.RegisterSLAsServer(s, svc.(cases.SLAsServer))
			},
			name: "SLAs",
		},
		{
			init: func(a *App) (interface{}, error) { return NewSLAConditionService(a) },
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
			init: func(a *App) (interface{}, error) { return NewServiceService(a) },
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
