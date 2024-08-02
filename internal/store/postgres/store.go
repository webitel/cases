package postgres

import (
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/model"
	"github.com/webitel/wlog"
)

type Store struct {
	config               *model.DatabaseConfig
	conn                 *sqlx.DB
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

func (s *Store) Database() (*sqlx.DB, model.AppError) {
	if s.conn == nil {
		model.NewInternalError("postgres.store.database.check.bad_arguments", "database connection is not opened")
	}
	return s.conn, nil
}

func (s *Store) Open() model.AppError {
	db, err := sqlx.Connect("pgx", s.config.Url)
	if err != nil {
		return model.NewInternalError("postgres.store.open.connect.fail", err.Error())
	}
	s.conn = db
	wlog.Debug(fmt.Sprintf("postgres: connection opened"))
	return nil
}

func (s *Store) Close() model.AppError {
	err := s.conn.Close()
	if err != nil {
		return model.NewInternalError("postgres.store.close.disconnect.fail", err.Error())
	}
	s.conn = nil
	wlog.Debug(fmt.Sprintf("postgres: connection closed"))
	return nil
}
