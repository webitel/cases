package postgres

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	conf "github.com/webitel/cases/config"
	dberr "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/store"
	otelpgx "github.com/webitel/webitel-go-kit/tracing/pgx"

	custom "github.com/webitel/custom/store"
)

// Store is the struct implementing the Store interface.
type Store struct {
	//------------cases stores ------------ ----//
	caseStore              store.CaseStore
	linkCaseStore          store.CaseLinkStore
	caseCommentStore       store.CaseCommentStore
	caseFileStore          store.CaseFileStore
	caseTimelineStore      store.CaseTimelineStore
	caseCommunicationStore store.CaseCommunicationStore
	relatedCaseStore       store.RelatedCaseStore
	//----------dictionary stores ------------ //
	sourceStore           store.SourceStore
	statusStore           store.StatusStore
	statusConditionStore  store.StatusConditionStore
	closeReasonGroupStore store.CloseReasonGroupStore
	closeReasonStore      store.CloseReasonStore
	priorityStore         store.PriorityStore
	slaStore              store.SLAStore
	slaConditionStore     store.SLAConditionStore
	catalogStore          store.CatalogStore
	serviceStore          store.ServiceStore
	config                *conf.DatabaseConfig
	conn                  *pgxpool.Pool

	// region: [custom] fields ..
	customStore custom.Catalog
	// endregion: [custom] fields ..
}

// New creates a new Store instance.
func New(config *conf.DatabaseConfig) *Store {
	return &Store{config: config}
}

// -------------Cases Stores ------------ //

func (s *Store) Case() store.CaseStore {
	if s.caseStore == nil {
		cs, err := NewCaseStore(s)
		if err != nil {
			return nil
		}
		s.caseStore = cs
	}
	return s.caseStore
}

func (s *Store) CaseLink() store.CaseLinkStore {
	if s.linkCaseStore == nil {
		linkCase, err := NewCaseLinkStore(s)
		if err != nil {
			return nil
		}
		s.linkCaseStore = linkCase
	}
	return s.linkCaseStore
}

func (s *Store) CaseComment() store.CaseCommentStore {
	if s.caseCommentStore == nil {
		caseComment, err := NewCaseCommentStore(s)
		if err != nil {
			return nil
		}
		s.caseCommentStore = caseComment
	}
	return s.caseCommentStore
}

func (s *Store) CaseFile() store.CaseFileStore {
	if s.caseFileStore == nil {
		caseFile, err := NewCaseFileStore(s)
		if err != nil {
			return nil
		}
		s.caseFileStore = caseFile
	}
	return s.caseFileStore
}

func (s *Store) CaseTimeline() store.CaseTimelineStore {
	if s.caseTimelineStore == nil {
		caseTimeline, err := NewCaseTimelineStore(s)
		if err != nil {
			return nil
		}
		s.caseTimelineStore = caseTimeline
	}
	return s.caseTimelineStore
}

func (s *Store) CaseCommunication() store.CaseCommunicationStore {
	if s.caseCommunicationStore == nil {
		caseCommunication, err := NewCaseCommunicationStore(s)
		if err != nil {
			return nil
		}
		s.caseCommunicationStore = caseCommunication
	}
	return s.caseCommunicationStore
}

func (s *Store) RelatedCase() store.RelatedCaseStore {
	if s.relatedCaseStore == nil {
		relatedCase, err := NewRelatedCaseStore(s)
		if err != nil {
			return nil
		}
		s.relatedCaseStore = relatedCase
	}
	return s.relatedCaseStore
}

// -------------Dictionary Stores ------------ //
func (s *Store) Status() store.StatusStore {
	if s.statusStore == nil {
		st, err := NewStatusStore(s)
		if err != nil {
			return nil
		}
		s.statusStore = st
	}
	return s.statusStore
}

func (s *Store) StatusCondition() store.StatusConditionStore {
	if s.statusConditionStore == nil {
		st, err := NewStatusConditionStore(s)
		if err != nil {
			return nil
		}
		s.statusConditionStore = st
	}
	return s.statusConditionStore
}

func (s *Store) Source() store.SourceStore {
	if s.sourceStore == nil {
		st, err := NewSourceStore(s)
		if err != nil {
			return nil
		}
		s.sourceStore = st
	}
	return s.sourceStore
}

func (s *Store) CloseReasonGroup() store.CloseReasonGroupStore {
	if s.closeReasonGroupStore == nil {
		st, err := NewCloseReasonGroupStore(s)
		if err != nil {
			return nil
		}
		s.closeReasonGroupStore = st
	}
	return s.closeReasonGroupStore
}

func (s *Store) CloseReason() store.CloseReasonStore {
	if s.closeReasonStore == nil {
		st, err := NewCloseReasonStore(s)
		if err != nil {
			return nil
		}
		s.closeReasonStore = st
	}
	return s.closeReasonStore
}

func (s *Store) Priority() store.PriorityStore {
	if s.priorityStore == nil {
		st, err := NewPriorityStore(s)
		if err != nil {
			return nil
		}
		s.priorityStore = st
	}
	return s.priorityStore
}

func (s *Store) SLA() store.SLAStore {
	if s.slaStore == nil {
		st, err := NewSLAStore(s)
		if err != nil {
			return nil
		}
		s.slaStore = st
	}
	return s.slaStore
}

func (s *Store) SLACondition() store.SLAConditionStore {
	if s.slaConditionStore == nil {
		sc, err := NewSLAConditionStore(s)
		if err != nil {
			return nil
		}
		s.slaConditionStore = sc
	}
	return s.slaConditionStore
}

func (s *Store) Catalog() store.CatalogStore {
	if s.catalogStore == nil {
		catalog, err := NewCatalogStore(s)
		if err != nil {
			return nil
		}
		s.catalogStore = catalog
	}
	return s.catalogStore
}

func (s *Store) Service() store.ServiceStore {
	if s.serviceStore == nil {
		service, err := NewServiceStore(s)
		if err != nil {
			return nil
		}
		s.serviceStore = service
	}
	return s.serviceStore
}

// Database returns the database connection or a custom error if it is not opened.
func (s *Store) Database() (*pgxpool.Pool, *dberr.DBError) { // Return custom DB error
	if s.conn == nil {
		return nil, dberr.NewDBError("store.database.check.bad_arguments", "database connection is not opened")
	}
	return s.conn, nil
}

// Open establishes a connection to the database and returns a custom error if it fails.
func (s *Store) Open() *dberr.DBError {
	config, err := pgxpool.ParseConfig(s.config.Url)
	if err != nil {
		return dberr.NewDBError("store.open.parse_config.fail", err.Error())
	}

	// Attach the OpenTelemetry tracer for pgx
	config.ConnConfig.Tracer = otelpgx.NewTracer(otelpgx.WithTrimSQLInSpanName())

	conn, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return dberr.NewDBError("store.open.connect.fail", err.Error())
	}
	s.conn = conn
	slog.Debug("cases.store.connection_opened", slog.String("message", "postgres: connection opened"))
	return nil
}

// Close closes the database connection and returns a custom error if it fails.
func (s *Store) Close() *dberr.DBError {
	if s.conn != nil {
		s.conn.Close()
		slog.Debug("cases.store.connection_closed", slog.String("message", "postgres: connection closed"))
		s.conn = nil
	}
	return nil
}
