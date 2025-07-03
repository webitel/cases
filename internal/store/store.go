package store

import (
	"context"

	"github.com/webitel/cases/auth"
	"github.com/webitel/cases/internal/model/options"

	_go "github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/internal/model"

	custom "github.com/webitel/custom/store"
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

	// ------------ Custom Store ------------ //
	Custom() custom.Catalog

	// ------------ Database Management ------------ //
	Open() error  // Return custom DB error
	Close() error // Return custom DB error
}

// ------------ Cases Stores ------------ //
type CaseStore interface {
	// Create a new case
	Create(rpc options.Creator, add *_go.Case) (*_go.Case, error)
	// List cases
	List(rpc options.Searcher, queryTarget *model.CaseQueryTarget) (*_go.CaseList, error)
	// Update case
	Update(req options.Updator, upd *_go.Case) (*_go.Case, error)
	// Delete case
	Delete(req options.Deleter) error
	// Check case by current auth options
	CheckRbacAccess(ctx context.Context, auth auth.Auther, access auth.AccessMode, caseId int64) (bool, error)
	SetOverdueCases(so options.Searcher) ([]*_go.Case, bool, error)
}

// RelatedCases attribute attached to the case (n:1)
type CaseLinkStore interface {
	// Create link
	Create(rpc options.Creator, add *model.CaseLink) (*model.CaseLink, error)
	// List links
	List(rpc options.Searcher) ([]*model.CaseLink, error)
	// Update link
	Update(req options.Updator, upd *model.CaseLink) (*model.CaseLink, error)
	// Delete link
	Delete(req options.Deleter) (*model.CaseLink, error)
}

// Comments attribute attached to the case (n:1)
type CaseCommentStore interface {
	// Create comment
	Publish(rpc options.Creator, add *_go.CaseComment) (*_go.CaseComment, error)
	// List comments
	List(rpc options.Searcher) (*_go.CaseCommentList, error)
	// Update comment
	Update(req options.Updator, upd *_go.CaseComment) (*_go.CaseComment, error)
	// Delete comment
	Delete(req options.Deleter) error
}

// Case timeline
type CaseTimelineStore interface {
	Get(rpc options.Searcher) (*_go.GetTimelineResponse, error)
	GetCounter(rpc options.Searcher) ([]*model.TimelineCounter, error)
}

// Case connected communications (chats, calls etc.)
type CaseCommunicationStore interface {
	Link(options.Creator, []*_go.InputCaseCommunication) ([]*_go.CaseCommunication, error)
	Unlink(options.Deleter) (int64, error)
	List(opts options.Searcher) (*_go.ListCommunicationsResponse, error)
}

type CaseFileStore interface {
	// List files
	List(rpc options.Searcher) ([]*model.CaseFile, error)
	// Delete Case | File association
	Delete(req options.Deleter) (*model.CaseFile, error)
}

type RelatedCaseStore interface {
	// Create relation
	Create(rpc options.Creator, relation *_go.RelationType, userID int64) (*_go.RelatedCase, error)
	// List related cases
	List(rpc options.Searcher) (*_go.RelatedCaseList, error)
	// Update relation
	Update(req options.Updator, upd *_go.InputRelatedCase, userID int64) (*_go.RelatedCase, error)
	// Delete relation
	Delete(req options.Deleter) error
}

// ------------Access Control------------//
type AccessControlStore interface {
	// Check if user has Rbac access
	RbacAccess(ctx context.Context, domainId, id int64, groups []int, access uint8, table string) (bool, error)
}

// ------------ Dictionary Stores ------------ //
type StatusStore interface {
	// Create a new status lookup
	Create(rpc options.Creator, input *model.Status) (*model.Status, error)
	// List status lookup
	List(rpc options.Searcher) ([]*model.Status, error)
	// Delete status lookup
	Delete(rpc options.Deleter) (*model.Status, error)
	// Update status lookup
	Update(rpc options.Updator, input *model.Status) (*model.Status, error)
}

type StatusConditionStore interface {
	// Create a new status сondition
	Create(ctx options.Creator, input *model.StatusCondition) (*model.StatusCondition, error)
	// List status сondition
	List(ctx options.Searcher) ([]*model.StatusCondition, error)
	// Delete status сondition
	Delete(ctx options.Deleter) (*model.StatusCondition, error)
	// Update status сondition
	Update(ctx options.Updator, input *model.StatusCondition) (*model.StatusCondition, error)
}

type CloseReasonGroupStore interface {
	// Create a new close reason lookup
	Create(rpc options.Creator, input *model.CloseReasonGroup) (*model.CloseReasonGroup, error)
	// List close reason lookup
	List(rpc options.Searcher) ([]*model.CloseReasonGroup, error)
	// Delete close reason lookup
	Delete(rpc options.Deleter) error
	// Update close reason lookup
	Update(rpc options.Updator, input *model.CloseReasonGroup) (*model.CloseReasonGroup, error)
}

type CloseReasonStore interface {
	// Create a new reason
	Create(creator options.Creator, input *model.CloseReason) (*model.CloseReason, error)
	// List reasons
	List(searcher options.Searcher, closeReasonId int64) ([]*model.CloseReason, error)
	// Delete reason
	Delete(deleter options.Deleter) (*model.CloseReason, error)
	// Update reason
	Update(updator options.Updator, input *model.CloseReason) (*model.CloseReason, error)
}

type SourceStore interface {
	// Create a new source lookup
	Create(rpc options.Creator, add *model.Source) (*model.Source, error)
	// List source lookup
	List(rpc options.Searcher) ([]*model.Source, error)
	// Delete source lookup
	Delete(rpc options.Deleter) (*model.Source, error)
	// Update source lookup
	Update(rpc options.Updator, lookup *model.Source) (*model.Source, error)
}

type PriorityStore interface {
	// Create a new priority lookup
	Create(rpc options.Creator, add *model.Priority) (*model.Priority, error)
	// List priority lookup
	List(rpc options.Searcher, notInSla int64, inSla int64) ([]*model.Priority, error)
	// Delete priority lookup
	Delete(rpc options.Deleter) (*model.Priority, error)
	// Update priority lookup
	Update(rpc options.Updator, lookup *model.Priority) (*model.Priority, error)
}

type SLAStore interface {
	// Create a new SLA lookup
	Create(rpc options.Creator, add *model.SLA) (*model.SLA, error)
	// List SLA lookup
	List(rpc options.Searcher) ([]*model.SLA, error)
	// Delete SLA lookup
	Delete(rpc options.Deleter) (*model.SLA, error)
	// Update SLA lookup
	Update(rpc options.Updator, input *model.SLA) (*model.SLA, error)
}

type SLAConditionStore interface {
	// Create a new SLA сondition
	Create(ctx options.Creator, add *model.SLACondition) (*model.SLACondition, error)
	// List SLA сondition
	List(ctx options.Searcher) ([]*model.SLACondition, error)
	// Delete SLA сondition
	Delete(ctx options.Deleter) (*model.SLACondition, error)
	// Update SLA сondition
	Update(ctx options.Updator, lookup *model.SLACondition) (*model.SLACondition, error)
}

// CatalogStore is parent store managing service catalogs.
type CatalogStore interface {
	// Create a new catalog
	Create(rpc options.Creator, add *_go.Catalog) (*_go.Catalog, error)
	// List catalogs
	List(rpc options.Searcher, depth int64, subfields []string, hasSubservices bool) (*_go.CatalogList, error)
	// Delete catalog
	Delete(rpc options.Deleter) error
	// Update catalog
	Update(rpc options.Updator, lookup *_go.Catalog) (*_go.Catalog, error)
}

// Service is child store managing services within catalogs.
type ServiceStore interface {
	// Create a new service
	Create(rpc options.Creator, add *model.Service) (*model.Service, error)
	// List services
	List(rpc options.Searcher) ([]*model.Service, error)
	// Delete service
	Delete(rpc options.Deleter) (*model.Service, error)
	// Update service
	Update(rpc options.Updator, lookup *model.Service) (*model.Service, error)
}
