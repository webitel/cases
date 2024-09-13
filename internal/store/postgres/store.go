package postgres

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/model"
	otelpgx "github.com/webitel/webitel-go-kit/tracing/pgx"
)

type Store struct {
	appealStore store.AppealStore

	//-----Status + StatusCondition
	statusStore          store.StatusStore
	statusConditionStore store.StatusConditionStore

	//-----CloseReason + Reason
	closeReasonStore store.CloseReasonStore
	reasonStore      store.ReasonStore

	priorityStore store.PriorityStore

	//-----SLA + SLACondition
	slaStore          store.SLAStore
	slaConditionStore store.SLAConditionStore

	//-----Catalog + Service
	catalogStore        store.CatalogStore
	serviceStore        store.ServiceStore
	accessControllStore store.AccessControlStore
	config              *model.DatabaseConfig
	conn                *pgxpool.Pool
}

func New(config *model.DatabaseConfig) *Store {
	return &Store{config: config}
}

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

func (s *Store) Appeal() store.AppealStore {
	if s.appealStore == nil {
		st, err := NewAppealStore(s)
		if err != nil {
			return nil
		}
		s.appealStore = st
	}
	return s.appealStore
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

func (s *Store) Reason() store.ReasonStore {
	if s.reasonStore == nil {
		st, err := NewReasonStore(s)
		if err != nil {
			return nil
		}
		s.reasonStore = st
	}
	return s.reasonStore
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

func (s *Store) AccessControl() store.AccessControlStore {
	if s.accessControllStore == nil {
		st, err := NewAccessControlStore(s)
		if err != nil {
			return nil
		}
		s.accessControllStore = st
	}
	return s.accessControllStore
}

func (s *Store) Database() (*pgxpool.Pool, model.AppError) {
	if s.conn == nil {
		return nil, model.NewInternalError("cases.store.database.check.bad_arguments", "database connection is not opened")
	}
	return s.conn, nil
}

func (s *Store) Open() model.AppError {
	config, err := pgxpool.ParseConfig(s.config.Url)
	if err != nil {
		return model.NewInternalError("cases.store.open.parse_config.fail", err.Error())
	}

	// Attach the OpenTelemetry tracer for pgx
	config.ConnConfig.Tracer = otelpgx.NewTracer(otelpgx.WithTrimSQLInSpanName())

	conn, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return model.NewInternalError("cases.store.open.connect.fail", err.Error())
	}
	s.conn = conn
	slog.Debug("cases.store.connection_opened", slog.String("message", "postgres: connection opened"))
	return nil
}

func (s *Store) Close() model.AppError {
	if s.conn != nil {
		s.conn.Close()
		slog.Debug("cases.store.connection_closed", slog.String("message", "postgres: connection closed"))
		s.conn = nil
	}
	return nil
}
