package store

import (
	"context"
	"github.com/webitel/cases/auth"

	"github.com/jackc/pgx/v5/pgxpool"
	_go "github.com/webitel/cases/api/cases"
	dberr "github.com/webitel/cases/internal/error"
	"github.com/webitel/cases/model"
)

// Store is an interface that defines all the methods and properties that a store should implement in Cases service

type Store interface {
	// ------------ Cases Stores ------------ //
	Case() CaseStore
	CaseComment() CaseCommentStore
	CaseLink() CaseLinkStore
	CaseFile() CaseFileStore
	CaseTimeline() CaseTimelineStore
	CaseCommunication() CaseCommunicationStore
	RelatedCase() RelatedCaseStore

	// ------------ Dictionary Stores ------------ //
	Source() SourceStore
	Priority() PriorityStore

	// ------------ Closure reasons Stores ------------ //
	CloseReasonGroup() CloseReasonGroupStore
	CloseReason() CloseReasonStore

	// ------------ Status ------------ //
	Status() StatusStore
	StatusCondition() StatusConditionStore

	// ------------ SLA Stores ------------ //
	SLA() SLAStore
	SLACondition() SLAConditionStore

	// ------------ Catalog and Service Stores ------------ //
	Catalog() CatalogStore
	Service() ServiceStore

	// ------------ Access Control ------------ //
	AccessControl() AccessControlStore

	// ------------ Database Management ------------ //
	Database() (*pgxpool.Pool, *dberr.DBError) // Return custom DB error
	Open() *dberr.DBError                      // Return custom DB error
	Close() *dberr.DBError                     // Return custom DB error
}

// ------------ Cases Stores ------------ //
type CaseStore interface {
	// Create a new case
	Create(rpc *model.CreateOptions, add *_go.Case) (*_go.Case, error)
	// List cases
	List(rpc *model.SearchOptions) (*_go.CaseList, error)
	// Update case
	Update(req *model.UpdateOptions, upd *_go.Case) (*_go.Case, error)
	// Delete case
	Delete(req *model.DeleteOptions) error
	// Check case by current auth options
	CheckRbacAccess(ctx context.Context, auth auth.Auther, access auth.AccessMode, caseId int64) (bool, error)
}

// RelatedCases attribute attached to the case (n:1)
type CaseLinkStore interface {
	// Create link
	Create(rpc *model.CreateOptions, add *_go.InputCaseLink) (*_go.CaseLink, error)
	// List links
	List(rpc *model.SearchOptions) (*_go.CaseLinkList, error)
	// Update link
	Update(req *model.UpdateOptions, upd *_go.InputCaseLink) (*_go.CaseLink, error)
	// Delete link
	Delete(req *model.DeleteOptions) error
}

// Comments attribute attached to the case (n:1)
type CaseCommentStore interface {
	// Create comment
	Publish(rpc *model.CreateOptions, add *_go.CaseComment) (*_go.CaseComment, error)
	// List comments
	List(rpc *model.SearchOptions) (*_go.CaseCommentList, error)
	// Update comment
	Update(req *model.UpdateOptions, upd *_go.CaseComment) (*_go.CaseComment, error)
	// Delete comment
	Delete(req *model.DeleteOptions) error
}

// Case timeline
type CaseTimelineStore interface {
	Get(rpc *model.SearchOptions) (*_go.GetTimelineResponse, error)
	GetCounter(rpc *model.SearchOptions) ([]*model.TimelineCounter, error)
}

// Case connected communications (chats, calls etc.)
type CaseCommunicationStore interface {
	Link(*model.CreateOptions, []*_go.InputCaseCommunication) ([]*_go.CaseCommunication, error)
	Unlink(*model.DeleteOptions) (int64, error)
	List(opts *model.SearchOptions) (*_go.ListCommunicationsResponse, error)
}

type CaseFileStore interface {
	// List files
	List(rpc *model.SearchOptions) (*_go.CaseFileList, error)
}

type RelatedCaseStore interface {
	// Create relation
	Create(rpc *model.CreateOptions, relation *_go.RelationType) (*_go.RelatedCase, error)
	// List related cases
	List(rpc *model.SearchOptions) (*_go.RelatedCaseList, error)
	// Update relation
	Update(req *model.UpdateOptions, upd *_go.InputRelatedCase) (*_go.RelatedCase, error)
	// Delete relation
	Delete(req *model.DeleteOptions) error
}

// ------------Access Control------------//
type AccessControlStore interface {
	// Check if user has Rbac access
	RbacAccess(ctx context.Context, domainId, id int64, groups []int, access uint8, table string) (bool, error)
}

// ------------ Dictionary Stores ------------ //
type StatusStore interface {
	// Create a new status lookup
	Create(rpc *model.CreateOptions, add *_go.Status) (*_go.Status, error)
	// List status lookup
	List(rpc *model.SearchOptions) (*_go.StatusList, error)
	// Delete status lookup
	Delete(rpc *model.DeleteOptions) error
	// Update status lookup
	Update(rpc *model.UpdateOptions, lookup *_go.Status) (*_go.Status, error)
}

type StatusConditionStore interface {
	// Create a new status сondition
	Create(ctx *model.CreateOptions, add *_go.StatusCondition) (*_go.StatusCondition, error)
	// List status сondition
	List(ctx *model.SearchOptions, statusId int64) (*_go.StatusConditionList, error)
	// Delete status сondition
	Delete(ctx *model.DeleteOptions, statusId int64) error
	// Update status сondition
	Update(ctx *model.UpdateOptions, status *_go.StatusCondition) (*_go.StatusCondition, error)
}

type CloseReasonGroupStore interface {
	// Create a new close reason lookup
	Create(rpc *model.CreateOptions, add *_go.CloseReasonGroup) (*_go.CloseReasonGroup, error)
	// List close reason lookup
	List(rpc *model.SearchOptions) (*_go.CloseReasonGroupList, error)
	// Delete close reason lookup
	Delete(rpc *model.DeleteOptions) error
	// Update close reason lookup
	Update(rpc *model.UpdateOptions, lookup *_go.CloseReasonGroup) (*_go.CloseReasonGroup, error)
}

type CloseReasonStore interface {
	// Create a new reason
	Create(ctx *model.CreateOptions, add *_go.CloseReason) (*_go.CloseReason, error)
	// List reasons
	List(ctx *model.SearchOptions, closeReasonId int64) (*_go.CloseReasonList, error)
	// Delete reason
	Delete(ctx *model.DeleteOptions, closeReasonId int64) error
	// Update reason
	Update(ctx *model.UpdateOptions, lookup *_go.CloseReason) (*_go.CloseReason, error)
}

type SourceStore interface {
	// Create a new source lookup
	Create(rpc *model.CreateOptions, add *_go.Source) (*_go.Source, error)
	// List source lookup
	List(rpc *model.SearchOptions) (*_go.SourceList, error)
	// Delete source lookup
	Delete(rpc *model.DeleteOptions) error
	// Update source lookup
	Update(rpc *model.UpdateOptions, lookup *_go.Source) (*_go.Source, error)
}

type PriorityStore interface {
	// Create a new priority lookup
	Create(rpc *model.CreateOptions, add *_go.Priority) (*_go.Priority, error)
	// List priority lookup
	List(rpc *model.SearchOptions, notInSla int64) (*_go.PriorityList, error)
	// Delete priority lookup
	Delete(rpc *model.DeleteOptions) error
	// Update priority lookup
	Update(rpc *model.UpdateOptions, lookup *_go.Priority) (*_go.Priority, error)
}

type SLAStore interface {
	// Create a new SLA lookup
	Create(rpc *model.CreateOptions, add *_go.SLA) (*_go.SLA, error)
	// List SLA lookup
	List(rpc *model.SearchOptions) (*_go.SLAList, error)
	// Delete SLA lookup
	Delete(rpc *model.DeleteOptions) error
	// Update SLA lookup
	Update(rpc *model.UpdateOptions, lookup *_go.SLA) (*_go.SLA, error)
}

type SLAConditionStore interface {
	// Create a new SLA сondition
	Create(ctx *model.CreateOptions, add *_go.SLACondition, priorities []int64) (*_go.SLACondition, error)
	// List SLA сondition
	List(ctx *model.SearchOptions) (*_go.SLAConditionList, error)
	// Delete SLA сondition
	Delete(ctx *model.DeleteOptions) error
	// Update SLA сondition
	Update(ctx *model.UpdateOptions, lookup *_go.SLACondition) (*_go.SLACondition, error)
}

// CatalogStore is parent store managing service catalogs.
type CatalogStore interface {
	// Create a new catalog
	Create(rpc *model.CreateOptions, add *_go.Catalog) (*_go.Catalog, error)
	// List catalogs
	List(rpc *model.SearchOptions, depth int64, subfields []string, hasSubservices bool) (*_go.CatalogList, error)
	// Delete catalog
	Delete(rpc *model.DeleteOptions) error
	// Update catalog
	Update(rpc *model.UpdateOptions, lookup *_go.Catalog) (*_go.Catalog, error)
}

// Service is child store managing services within catalogs.
type ServiceStore interface {
	// Create a new service
	Create(rpc *model.CreateOptions, add *_go.Service) (*_go.Service, error)
	// List services
	List(rpc *model.SearchOptions) (*_go.ServiceList, error)
	// Delete service
	Delete(rpc *model.DeleteOptions) error
	// Update service
	Update(rpc *model.UpdateOptions, lookup *_go.Service) (*_go.Service, error)
}
