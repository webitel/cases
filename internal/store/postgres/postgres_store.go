package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/webitel/cases/internal/store"
	lookup2 "github.com/webitel/cases/internal/store/postgres/lookup"
	"github.com/webitel/cases/model"
	"github.com/webitel/wlog"
)

type Store struct {
	config                 *model.DatabaseConfig
	conn                   *sqlx.DB
	appealLookupStore      store.AppealLookupStore
	statusLookupStore      store.StatusLookupStore
	closeReasonLookupStore store.CloseReasonLookupStore
}

func New(config *model.DatabaseConfig) *Store {
	return &Store{config: config}
}

func (s *Store) AppealLookup() store.AppealLookupStore {
	if s.appealLookupStore == nil {
		log, err := lookup2.NewAppealLookupStore(s)
		if err != nil {
			return nil
		}
		s.appealLookupStore = log
	}
	return s.appealLookupStore
}

func (s *Store) CloseReasonLookup() store.CloseReasonLookupStore {
	if s.closeReasonLookupStore == nil {
		log, err := lookup2.NewCloseReasonLookupStore(s)
		if err != nil {
			return nil
		}
		s.closeReasonLookupStore = log
	}
	return s.closeReasonLookupStore
}
func (s *Store) StatusLookup() store.StatusLookupStore {
	if s.statusLookupStore == nil {
		log, err := lookup2.NewStatusLookupStore(s)
		if err != nil {
			return nil
		}
		s.statusLookupStore = log
	}
	return s.statusLookupStore
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
