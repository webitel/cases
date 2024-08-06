package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/model"
	"github.com/webitel/wlog"
)

type Store struct {
	config               *model.DatabaseConfig
	conn                 *pgxpool.Pool
	appealStore          store.AppealStore
	statusConditionStore store.StatusConditionStore
	closeReasonStore     store.CloseReasonStore
	statusStore          store.StatusStore
}

func New(config *model.DatabaseConfig) *Store {
	return &Store{config: config}
}

func (s *Store) Appeal() store.AppealStore {
	if s.appealStore == nil {
		log, err := NewAppealStore(s)
		if err != nil {
			return nil
		}
		s.appealStore = log
	}
	return s.appealStore
}

func (s *Store) CloseReason() store.CloseReasonStore {
	if s.closeReasonStore == nil {
		log, err := NewCloseReasonStore(s)
		if err != nil {
			return nil
		}
		s.closeReasonStore = log
	}
	return s.closeReasonStore
}

func (s *Store) Status() store.StatusStore {
	if s.statusStore == nil {
		log, err := NewStatusStore(s)
		if err != nil {
			return nil
		}
		s.statusStore = log
	}
	return s.statusStore
}

func (s *Store) StatusCondition() store.StatusConditionStore {
	if s.statusConditionStore == nil {
		log, err := NewStatusConditionStore(s)
		if err != nil {
			return nil
		}
		s.statusConditionStore = log
	}
	return s.statusConditionStore
}

func (s *Store) Database() (*pgxpool.Pool, model.AppError) {
	if s.conn == nil {
		return nil, model.NewInternalError("postgres.store.database.check.bad_arguments", "database connection is not opened")
	}
	return s.conn, nil
}

func (s *Store) Open() model.AppError {
	config, err := pgxpool.ParseConfig(s.config.Url)
	if err != nil {
		return model.NewInternalError("postgres.store.open.parse_config.fail", err.Error())
	}

	conn, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return model.NewInternalError("postgres.store.open.connect.fail", err.Error())
	}
	s.conn = conn
	wlog.Debug("postgres: connection opened")
	return nil
}

func (s *Store) Close() model.AppError {
	if s.conn != nil {
		s.conn.Close()
		wlog.Debug("postgres: connection closed")
		s.conn = nil
	}
	return nil
}
