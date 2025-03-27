// Code generated by protoc-gen-go-webitel. DO NOT EDIT.
// versions:
// - protoc-gen-go-webitel v1.0.0
// - protoc                (unknown)

package cases

// WebitelServicesInfo is the list of services defined in proto files.
type WebitelServicesInfo map[string]WebitelServices

type WebitelServices struct {
	ObjClass           string
	AdditionalLicenses []string
	WebitelMethods     map[string]WebitelMethod
}

// WebitelMethod is the list of methods defined in this service.
type WebitelMethod struct {
	HttpBindings []*HttpBinding
	Access       int
	Input        string
	Output       string
}

type HttpBinding struct {
	Path   string
	Method string
}

var WebitelAPI = WebitelServicesInfo{
	"Services": WebitelServices{
		ObjClass:           "dictionaries",
		AdditionalLicenses: []string{},
		WebitelMethods: map[string]WebitelMethod{
			"ListServices": WebitelMethod{
				Access: 1,
				Input:  "ListServiceRequest",
				Output: "ServiceList",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/services",
						Method: "GET",
					},
				},
			},
			"CreateService": WebitelMethod{
				Access: 2,
				Input:  "CreateServiceRequest",
				Output: "Service",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/services",
						Method: "POST",
					},
				},
			},
			"UpdateService": WebitelMethod{
				Access: 2,
				Input:  "UpdateServiceRequest",
				Output: "Service",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/services/{id}",
						Method: "PUT",
					},
					{
						Path:   "/cases/services/{id}",
						Method: "PATCH",
					},
				},
			},
			"DeleteService": WebitelMethod{
				Access: 2,
				Input:  "DeleteServiceRequest",
				Output: "ServiceList",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/services/{id}",
						Method: "DELETE",
					},
				},
			},
			"LocateService": WebitelMethod{
				Access: 1,
				Input:  "LocateServiceRequest",
				Output: "LocateServiceResponse",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/services/{id}",
						Method: "GET",
					},
				},
			},
		},
	},
	"CaseComments": WebitelServices{
		ObjClass:           "case_comments",
		AdditionalLicenses: []string{},
		WebitelMethods: map[string]WebitelMethod{
			"LocateComment": WebitelMethod{
				Access: 1,
				Input:  "LocateCommentRequest",
				Output: "CaseComment",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/comments/{etag}",
						Method: "GET",
					},
				},
			},
			"UpdateComment": WebitelMethod{
				Access: 2,
				Input:  "UpdateCommentRequest",
				Output: "CaseComment",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/comments/{input.etag}",
						Method: "PUT",
					},
					{
						Path:   "/cases/comments/{input.etag}",
						Method: "PATCH",
					},
				},
			},
			"DeleteComment": WebitelMethod{
				Access: 3,
				Input:  "DeleteCommentRequest",
				Output: "CaseComment",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/comments/{etag}",
						Method: "DELETE",
					},
				},
			},
			"ListComments": WebitelMethod{
				Access: 1,
				Input:  "ListCommentsRequest",
				Output: "CaseCommentList",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/{case_etag}/comments",
						Method: "GET",
					},
				},
			},
			"PublishComment": WebitelMethod{
				Access: 0,
				Input:  "PublishCommentRequest",
				Output: "CaseComment",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/{case_etag}/comments",
						Method: "POST",
					},
				},
			},
		},
	},
	"RelatedCases": WebitelServices{
		ObjClass:           "cases",
		AdditionalLicenses: []string{},
		WebitelMethods: map[string]WebitelMethod{
			"LocateRelatedCase": WebitelMethod{
				Access: 1,
				Input:  "LocateRelatedCaseRequest",
				Output: "RelatedCase",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/{primary_case_etag}/related/{etag}",
						Method: "GET",
					},
				},
			},
			"CreateRelatedCase": WebitelMethod{
				Access: 2,
				Input:  "CreateRelatedCaseRequest",
				Output: "RelatedCase",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/{primary_case_etag}/related",
						Method: "POST",
					},
				},
			},
			"UpdateRelatedCase": WebitelMethod{
				Access: 2,
				Input:  "UpdateRelatedCaseRequest",
				Output: "RelatedCase",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/{input.primary_case.id}/related/{etag}",
						Method: "PUT",
					},
					{
						Path:   "/cases/{input.primary_case.id}/related/{etag}",
						Method: "PATCH",
					},
				},
			},
			"DeleteRelatedCase": WebitelMethod{
				Access: 2,
				Input:  "DeleteRelatedCaseRequest",
				Output: "RelatedCase",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/{primary_case_etag}/related/{etag}",
						Method: "DELETE",
					},
				},
			},
			"ListRelatedCases": WebitelMethod{
				Access: 1,
				Input:  "ListRelatedCasesRequest",
				Output: "RelatedCaseList",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/{primary_case_etag}/related",
						Method: "GET",
					},
				},
			},
		},
	},
	"CaseFiles": WebitelServices{
		ObjClass:           "cases",
		AdditionalLicenses: []string{},
		WebitelMethods: map[string]WebitelMethod{
			"ListFiles": WebitelMethod{
				Access: 1,
				Input:  "ListFilesRequest",
				Output: "CaseFileList",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/{case_etag}/files",
						Method: "GET",
					},
				},
			},
			"DeleteFile": WebitelMethod{
				Access: 3,
				Input:  "DeleteFileRequest",
				Output: "File",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/{case_etag}/files/{id}",
						Method: "DELETE",
					},
				},
			},
		},
	},
	"CaseLinks": WebitelServices{
		ObjClass:           "cases",
		AdditionalLicenses: []string{},
		WebitelMethods: map[string]WebitelMethod{
			"LocateLink": WebitelMethod{
				Access: 1,
				Input:  "LocateLinkRequest",
				Output: "CaseLink",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/{case_etag}/links/{etag}",
						Method: "GET",
					},
				},
			},
			"CreateLink": WebitelMethod{
				Access: 2,
				Input:  "CreateLinkRequest",
				Output: "CaseLink",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/{case_etag}/links",
						Method: "POST",
					},
				},
			},
			"UpdateLink": WebitelMethod{
				Access: 2,
				Input:  "UpdateLinkRequest",
				Output: "CaseLink",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/{case_etag}/links/{input.etag}",
						Method: "PUT",
					},
					{
						Path:   "/cases/{case_etag}/links/{input.etag}",
						Method: "PATCH",
					},
				},
			},
			"DeleteLink": WebitelMethod{
				Access: 2,
				Input:  "DeleteLinkRequest",
				Output: "CaseLink",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/{case_etag}/links/{etag}",
						Method: "DELETE",
					},
				},
			},
			"ListLinks": WebitelMethod{
				Access: 1,
				Input:  "ListLinksRequest",
				Output: "CaseLinkList",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/{case_etag}/links",
						Method: "GET",
					},
				},
			},
		},
	},
	"Priorities": WebitelServices{
		ObjClass:           "dictionaries",
		AdditionalLicenses: []string{},
		WebitelMethods: map[string]WebitelMethod{
			"ListPriorities": WebitelMethod{
				Access: 1,
				Input:  "ListPriorityRequest",
				Output: "PriorityList",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/priorities",
						Method: "GET",
					},
				},
			},
			"CreatePriority": WebitelMethod{
				Access: 0,
				Input:  "CreatePriorityRequest",
				Output: "Priority",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/priorities",
						Method: "POST",
					},
				},
			},
			"UpdatePriority": WebitelMethod{
				Access: 2,
				Input:  "UpdatePriorityRequest",
				Output: "Priority",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/priorities/{id}",
						Method: "PUT",
					},
					{
						Path:   "/cases/priorities/{id}",
						Method: "PATCH",
					},
				},
			},
			"DeletePriority": WebitelMethod{
				Access: 3,
				Input:  "DeletePriorityRequest",
				Output: "Priority",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/priorities/{id}",
						Method: "DELETE",
					},
				},
			},
			"LocatePriority": WebitelMethod{
				Access: 1,
				Input:  "LocatePriorityRequest",
				Output: "LocatePriorityResponse",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/priorities/{id}",
						Method: "GET",
					},
				},
			},
		},
	},
	"StatusConditions": WebitelServices{
		ObjClass:           "dictionaries",
		AdditionalLicenses: []string{},
		WebitelMethods: map[string]WebitelMethod{
			"ListStatusConditions": WebitelMethod{
				Access: 1,
				Input:  "ListStatusConditionRequest",
				Output: "StatusConditionList",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/statuses/{status_id}/status",
						Method: "GET",
					},
				},
			},
			"CreateStatusCondition": WebitelMethod{
				Access: 2,
				Input:  "CreateStatusConditionRequest",
				Output: "StatusCondition",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/statuses/{status_id}/status",
						Method: "POST",
					},
				},
			},
			"UpdateStatusCondition": WebitelMethod{
				Access: 2,
				Input:  "UpdateStatusConditionRequest",
				Output: "StatusCondition",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/statuses/{status_id}/status/{id}",
						Method: "PUT",
					},
					{
						Path:   "/statuses/{status_id}/status/{id}",
						Method: "PATCH",
					},
				},
			},
			"DeleteStatusCondition": WebitelMethod{
				Access: 2,
				Input:  "DeleteStatusConditionRequest",
				Output: "StatusCondition",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/statuses/{status_id}/status/{id}",
						Method: "DELETE",
					},
				},
			},
			"LocateStatusCondition": WebitelMethod{
				Access: 1,
				Input:  "LocateStatusConditionRequest",
				Output: "LocateStatusConditionResponse",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/statuses/{status_id}/status/{id}",
						Method: "GET",
					},
				},
			},
		},
	},
	"Sources": WebitelServices{
		ObjClass:           "dictionaries",
		AdditionalLicenses: []string{},
		WebitelMethods: map[string]WebitelMethod{
			"ListSources": WebitelMethod{
				Access: 1,
				Input:  "ListSourceRequest",
				Output: "SourceList",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/sources",
						Method: "GET",
					},
				},
			},
			"CreateSource": WebitelMethod{
				Access: 0,
				Input:  "CreateSourceRequest",
				Output: "Source",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/sources",
						Method: "POST",
					},
				},
			},
			"UpdateSource": WebitelMethod{
				Access: 2,
				Input:  "UpdateSourceRequest",
				Output: "Source",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/sources/{id}",
						Method: "PUT",
					},
					{
						Path:   "/cases/sources/{id}",
						Method: "PATCH",
					},
				},
			},
			"DeleteSource": WebitelMethod{
				Access: 3,
				Input:  "DeleteSourceRequest",
				Output: "Source",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/sources/{id}",
						Method: "DELETE",
					},
				},
			},
			"LocateSource": WebitelMethod{
				Access: 1,
				Input:  "LocateSourceRequest",
				Output: "LocateSourceResponse",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/sources/{id}",
						Method: "GET",
					},
				},
			},
		},
	},
	"SLAConditions": WebitelServices{
		ObjClass:           "dictionaries",
		AdditionalLicenses: []string{},
		WebitelMethods: map[string]WebitelMethod{
			"ListSLAConditions": WebitelMethod{
				Access: 1,
				Input:  "ListSLAConditionRequest",
				Output: "SLAConditionList",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/slas/{sla_id}/sla_conditions",
						Method: "GET",
					},
				},
			},
			"CreateSLACondition": WebitelMethod{
				Access: 2,
				Input:  "CreateSLAConditionRequest",
				Output: "SLACondition",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/slas/{sla_id}/sla_condition",
						Method: "POST",
					},
				},
			},
			"UpdateSLACondition": WebitelMethod{
				Access: 2,
				Input:  "UpdateSLAConditionRequest",
				Output: "SLACondition",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/slas/{sla_id}/sla_condition/{id}",
						Method: "PUT",
					},
					{
						Path:   "/slas/{sla_id}/sla_condition/{id}",
						Method: "PATCH",
					},
				},
			},
			"DeleteSLACondition": WebitelMethod{
				Access: 2,
				Input:  "DeleteSLAConditionRequest",
				Output: "SLACondition",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/slas/{sla_id}/sla_condition/{id}",
						Method: "DELETE",
					},
				},
			},
			"LocateSLACondition": WebitelMethod{
				Access: 1,
				Input:  "LocateSLAConditionRequest",
				Output: "LocateSLAConditionResponse",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/slas/{sla_id}/sla_condition/{id}",
						Method: "GET",
					},
				},
			},
		},
	},
	"Cases": WebitelServices{
		ObjClass:           "cases",
		AdditionalLicenses: []string{},
		WebitelMethods: map[string]WebitelMethod{
			"SearchCases": WebitelMethod{
				Access: 1,
				Input:  "SearchCasesRequest",
				Output: "CaseList",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases",
						Method: "GET",
					},
					{
						Path:   "/contacts/{contact_id}/cases",
						Method: "GET",
					},
				},
			},
			"LocateCase": WebitelMethod{
				Access: 1,
				Input:  "LocateCaseRequest",
				Output: "Case",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/{etag}",
						Method: "GET",
					},
				},
			},
			"CreateCase": WebitelMethod{
				Access: 0,
				Input:  "CreateCaseRequest",
				Output: "Case",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases",
						Method: "POST",
					},
				},
			},
			"UpdateCase": WebitelMethod{
				Access: 2,
				Input:  "UpdateCaseRequest",
				Output: "Case",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/{input.etag}",
						Method: "PUT",
					},
					{
						Path:   "/cases/{input.etag}",
						Method: "PATCH",
					},
				},
			},
			"DeleteCase": WebitelMethod{
				Access: 3,
				Input:  "DeleteCaseRequest",
				Output: "Case",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/{etag}",
						Method: "DELETE",
					},
				},
			},
		},
	},
	"CaseCommunications": WebitelServices{
		ObjClass:           "cases",
		AdditionalLicenses: []string{},
		WebitelMethods: map[string]WebitelMethod{
			"LinkCommunication": WebitelMethod{
				Access: 1,
				Input:  "LinkCommunicationRequest",
				Output: "LinkCommunicationResponse",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/{case_etag}/communication",
						Method: "POST",
					},
				},
			},
			"UnlinkCommunication": WebitelMethod{
				Access: 2,
				Input:  "UnlinkCommunicationRequest",
				Output: "UnlinkCommunicationResponse",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/{case_etag}/communication/{id}",
						Method: "DELETE",
					},
				},
			},
			"ListCommunications": WebitelMethod{
				Access: 1,
				Input:  "ListCommunicationsRequest",
				Output: "ListCommunicationsResponse",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/{case_etag}/communication",
						Method: "GET",
					},
				},
			},
		},
	},
	"CaseTimeline": WebitelServices{
		ObjClass:           "cases",
		AdditionalLicenses: []string{},
		WebitelMethods: map[string]WebitelMethod{
			"GetTimeline": WebitelMethod{
				Access: 1,
				Input:  "GetTimelineRequest",
				Output: "GetTimelineResponse",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/{case_id}/timeline",
						Method: "GET",
					},
				},
			},
			"GetTimelineCounter": WebitelMethod{
				Access: 1,
				Input:  "GetTimelineCounterRequest",
				Output: "GetTimelineCounterResponse",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/{case_id}/timeline/counter",
						Method: "GET",
					},
				},
			},
		},
	},
	"Catalogs": WebitelServices{
		ObjClass:           "dictionaries",
		AdditionalLicenses: []string{},
		WebitelMethods: map[string]WebitelMethod{
			"ListCatalogs": WebitelMethod{
				Access: 1,
				Input:  "ListCatalogRequest",
				Output: "CatalogList",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/catalogs",
						Method: "GET",
					},
				},
			},
			"CreateCatalog": WebitelMethod{
				Access: 0,
				Input:  "CreateCatalogRequest",
				Output: "Catalog",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/catalogs",
						Method: "POST",
					},
				},
			},
			"UpdateCatalog": WebitelMethod{
				Access: 2,
				Input:  "UpdateCatalogRequest",
				Output: "Catalog",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/catalogs/{id}",
						Method: "PUT",
					},
					{
						Path:   "/cases/catalogs/{id}",
						Method: "PATCH",
					},
				},
			},
			"DeleteCatalog": WebitelMethod{
				Access: 3,
				Input:  "DeleteCatalogRequest",
				Output: "CatalogList",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/catalogs/{id}",
						Method: "DELETE",
					},
				},
			},
			"LocateCatalog": WebitelMethod{
				Access: 1,
				Input:  "LocateCatalogRequest",
				Output: "LocateCatalogResponse",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/catalogs/{id}",
						Method: "GET",
					},
				},
			},
		},
	},
	"CloseReasons": WebitelServices{
		ObjClass:           "dictionaries",
		AdditionalLicenses: []string{},
		WebitelMethods: map[string]WebitelMethod{
			"ListCloseReasons": WebitelMethod{
				Access: 1,
				Input:  "ListCloseReasonRequest",
				Output: "CloseReasonList",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/close_reason_groups/{close_reason_group_id}/close_reasons",
						Method: "GET",
					},
				},
			},
			"CreateCloseReason": WebitelMethod{
				Access: 2,
				Input:  "CreateCloseReasonRequest",
				Output: "CloseReason",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/close_reason_groups/{close_reason_group_id}/close_reasons",
						Method: "POST",
					},
				},
			},
			"UpdateCloseReason": WebitelMethod{
				Access: 2,
				Input:  "UpdateCloseReasonRequest",
				Output: "CloseReason",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/close_reason_groups/{close_reason_group_id}/close_reasons/{id}",
						Method: "PUT",
					},
					{
						Path:   "/close_reason_groups/{close_reason_group_id}/close_reasons/{id}",
						Method: "PATCH",
					},
				},
			},
			"DeleteCloseReason": WebitelMethod{
				Access: 2,
				Input:  "DeleteCloseReasonRequest",
				Output: "CloseReason",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/close_reason_groups/{close_reason_group_id}/close_reasons/{id}",
						Method: "DELETE",
					},
				},
			},
			"LocateCloseReason": WebitelMethod{
				Access: 1,
				Input:  "LocateCloseReasonRequest",
				Output: "LocateCloseReasonResponse",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/close_reason_groups/{close_reason_group_id}/close_reasons/{id}",
						Method: "GET",
					},
				},
			},
		},
	},
	"CloseReasonGroups": WebitelServices{
		ObjClass:           "dictionaries",
		AdditionalLicenses: []string{},
		WebitelMethods: map[string]WebitelMethod{
			"ListCloseReasonGroups": WebitelMethod{
				Access: 1,
				Input:  "ListCloseReasonGroupsRequest",
				Output: "CloseReasonGroupList",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/close_reason_groups",
						Method: "GET",
					},
				},
			},
			"CreateCloseReasonGroup": WebitelMethod{
				Access: 0,
				Input:  "CreateCloseReasonGroupRequest",
				Output: "CloseReasonGroup",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/close_reason_groups",
						Method: "POST",
					},
				},
			},
			"UpdateCloseReasonGroup": WebitelMethod{
				Access: 2,
				Input:  "UpdateCloseReasonGroupRequest",
				Output: "CloseReasonGroup",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/close_reason_groups/{id}",
						Method: "PUT",
					},
					{
						Path:   "/cases/close_reason_groups/{id}",
						Method: "PATCH",
					},
				},
			},
			"DeleteCloseReasonGroup": WebitelMethod{
				Access: 3,
				Input:  "DeleteCloseReasonGroupRequest",
				Output: "CloseReasonGroup",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/close_reason_groups/{id}",
						Method: "DELETE",
					},
				},
			},
			"LocateCloseReasonGroup": WebitelMethod{
				Access: 1,
				Input:  "LocateCloseReasonGroupRequest",
				Output: "LocateCloseReasonGroupResponse",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/close_reason_groups/{id}",
						Method: "GET",
					},
				},
			},
		},
	},
	"SLAs": WebitelServices{
		ObjClass:           "dictionaries",
		AdditionalLicenses: []string{},
		WebitelMethods: map[string]WebitelMethod{
			"ListSLAs": WebitelMethod{
				Access: 1,
				Input:  "ListSLARequest",
				Output: "SLAList",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/slas",
						Method: "GET",
					},
				},
			},
			"CreateSLA": WebitelMethod{
				Access: 0,
				Input:  "CreateSLARequest",
				Output: "SLA",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/slas",
						Method: "POST",
					},
				},
			},
			"UpdateSLA": WebitelMethod{
				Access: 2,
				Input:  "UpdateSLARequest",
				Output: "SLA",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/slas/{id}",
						Method: "PUT",
					},
					{
						Path:   "/cases/slas/{id}",
						Method: "PATCH",
					},
				},
			},
			"DeleteSLA": WebitelMethod{
				Access: 3,
				Input:  "DeleteSLARequest",
				Output: "SLA",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/slas/{id}",
						Method: "DELETE",
					},
				},
			},
			"LocateSLA": WebitelMethod{
				Access: 1,
				Input:  "LocateSLARequest",
				Output: "LocateSLAResponse",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/slas/{id}",
						Method: "GET",
					},
				},
			},
		},
	},
	"Statuses": WebitelServices{
		ObjClass:           "dictionaries",
		AdditionalLicenses: []string{},
		WebitelMethods: map[string]WebitelMethod{
			"ListStatuses": WebitelMethod{
				Access: 1,
				Input:  "ListStatusRequest",
				Output: "StatusList",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/statuses",
						Method: "GET",
					},
				},
			},
			"CreateStatus": WebitelMethod{
				Access: 0,
				Input:  "CreateStatusRequest",
				Output: "Status",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/statuses",
						Method: "POST",
					},
				},
			},
			"UpdateStatus": WebitelMethod{
				Access: 2,
				Input:  "UpdateStatusRequest",
				Output: "Status",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/statuses/{id}",
						Method: "PUT",
					},
					{
						Path:   "/cases/statuses/{id}",
						Method: "PATCH",
					},
				},
			},
			"DeleteStatus": WebitelMethod{
				Access: 3,
				Input:  "DeleteStatusRequest",
				Output: "Status",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/statuses/{id}",
						Method: "DELETE",
					},
				},
			},
			"LocateStatus": WebitelMethod{
				Access: 1,
				Input:  "LocateStatusRequest",
				Output: "LocateStatusResponse",
				HttpBindings: []*HttpBinding{
					{
						Path:   "/cases/statuses/{id}",
						Method: "GET",
					},
				},
			},
		},
	},
}
