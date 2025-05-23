package store

import (
	"context"
	"github.com/webitel/cases/auth"
	"github.com/webitel/cases/model/options"

	_go "github.com/webitel/cases/api/cases"
	dberr "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/model"

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
	Open() *dberr.DBError  // Return custom DB error
	Close() *dberr.DBError // Return custom DB error
}

// ------------ Cases Stores ------------ //
type CaseStore interface {
	// Create a new case
	Create(rpc options.Creator, add *_go.Case) (*_go.Case, error)
	// List cases
	List(rpc options.Searcher) (*_go.CaseList, error)
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
	Create(rpc options.Creator, add *_go.InputCaseLink) (*_go.CaseLink, error)
	// List links
	List(rpc options.Searcher) (*_go.CaseLinkList, error)
	// Update link
	Update(req options.Updator, upd *_go.InputCaseLink) (*_go.CaseLink, error)
	// Delete link
	Delete(req options.Deleter) error
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
	List(rpc options.Searcher) (*_go.CaseFileList, error)
	// Delete Case | File association
	Delete(req options.Deleter) error
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
	Create(rpc options.Creator, input *_go.Status) (*_go.Status, error)
	// List status lookup
	List(rpc options.Searcher) (*_go.StatusList, error)
	// Delete status lookup
	Delete(rpc options.Deleter) error
	// Update status lookup
	Update(rpc options.Updator, input *_go.Status) (*_go.Status, error)
}

type StatusConditionStore interface {
	// Create a new status сondition
	Create(ctx options.Creator, input *_go.StatusCondition) (*_go.StatusCondition, error)
	// List status сondition
	List(ctx options.Searcher, statusId int64) (*_go.StatusConditionList, error)
	// Delete status сondition
	Delete(ctx options.Deleter, statusId int64) error
	// Update status сondition
	Update(ctx options.Updator, input *_go.StatusCondition) (*_go.StatusCondition, error)
}

type CloseReasonGroupStore interface {
	// Create a new close reason lookup
	Create(rpc options.Creator, input *_go.CloseReasonGroup) (*_go.CloseReasonGroup, error)
	// List close reason lookup
	List(rpc options.Searcher) (*_go.CloseReasonGroupList, error)
	// Delete close reason lookup
	Delete(rpc options.Deleter) error
	// Update close reason lookup
	Update(rpc options.Updator, input *_go.CloseReasonGroup) (*_go.CloseReasonGroup, error)
}

type CloseReasonStore interface {
	// Create a new reason
	Create(ctx options.Creator, input *_go.CloseReason) (*_go.CloseReason, error)
	// List reasons
	List(ctx options.Searcher, closeReasonId int64) (*_go.CloseReasonList, error)
	// Delete reason
	Delete(ctx options.Deleter) error
	// Update reason
	Update(ctx options.Updator, input *_go.CloseReason) (*_go.CloseReason, error)
}

type SourceStore interface {
	// Create a new source lookup
	Create(rpc options.Creator, add *_go.Source) (*_go.Source, error)
	// List source lookup
	List(rpc options.Searcher) (*_go.SourceList, error)
	// Delete source lookup
	Delete(rpc options.Deleter) error
	// Update source lookup
	Update(rpc options.Updator, lookup *_go.Source) (*_go.Source, error)
}

type PriorityStore interface {
	// Create a new priority lookup
	Create(rpc options.Creator, add *_go.Priority) (*_go.Priority, error)
	// List priority lookup
	List(rpc options.Searcher, notInSla int64, inSla int64) (*_go.PriorityList, error)
	// Delete priority lookup
	Delete(rpc options.Deleter) error
	// Update priority lookup
	Update(rpc options.Updator, lookup *_go.Priority) (*_go.Priority, error)
}

type SLAStore interface {
	// Create a new SLA lookup
	Create(rpc options.Creator, input *_go.SLA) (*_go.SLA, error)
	// List SLA lookup
	List(rpc options.Searcher) (*_go.SLAList, error)
	// Delete SLA lookup
	Delete(rpc options.Deleter) error
	// Update SLA lookup
	Update(rpc options.Updator, input *_go.SLA) (*_go.SLA, error)
}

type SLAConditionStore interface {
	// Create a new SLA сondition
	Create(ctx options.Creator, add *_go.SLACondition, priorities []int64) (*_go.SLACondition, error)
	// List SLA сondition
	List(ctx options.Searcher) (*_go.SLAConditionList, error)
	// Delete SLA сondition
	Delete(ctx options.Deleter) error
	// Update SLA сondition
	Update(ctx options.Updator, lookup *_go.SLACondition) (*_go.SLACondition, error)
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
	Create(rpc options.Creator, add *_go.Service) (*_go.Service, error)
	// List services
	List(rpc options.Searcher) (*_go.ServiceList, error)
	// Delete service
	Delete(rpc options.Deleter) error
	// Update service
	Update(rpc options.Updator, lookup *_go.Service) (*_go.Service, error)
}
