package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/webitel/cases/model"
	"github.com/webitel/cases/storage"
	"github.com/webitel/cases/storage/postgres/lookup"
	"github.com/webitel/wlog"
)

type Store struct {
	config                 *model.DatabaseConfig
	conn                   *sqlx.DB
	appealLookupStore      storage.AppealLookupStore
	statusLookupStore      storage.StatusLookupStore
	closeReasonLookupStore storage.CloseReasonLookupStore
}

func New(config *model.DatabaseConfig) *Store {
	return &Store{config: config}
}

func (s *Store) Appeal() storage.AppealLookupStore {
	if s.appealLookupStore == nil {
		log, err := lookup.NewAppealLookupStore(s)
		if err != nil {
			return nil
		}
		s.appealLookupStore = log
	}
	return s.appealLookupStore
}

func (s *Store) CloseReason() storage.CloseReasonLookupStore {
	if s.closeReasonLookupStore == nil {
		log, err := lookup.NewCloseReasonLookupStore(s)
		if err != nil {
			return nil
		}
		s.closeReasonLookupStore = log
	}
	return s.closeReasonLookupStore
}
func (s *Store) Status() storage.StatusLookupStore {
	if s.statusLookupStore == nil {
		log, err := lookup.NewStatusLookupStore(s)
		if err != nil {
			return nil
		}
		s.statusLookupStore = log
	}
	return s.statusLookupStore
}

func (s *Store) Database() (*sqlx.DB, model.AppError) {
	if s.conn == nil {
		model.NewInternalError("postgres.storage.database.check.bad_arguments", "database connection is not opened")
	}
	return s.conn, nil
}

func (s *Store) Open() model.AppError {
	db, err := sqlx.Connect("pgx", s.config.Url)
	if err != nil {
		return model.NewInternalError("postgres.storage.open.connect.fail", err.Error())
	}
	s.conn = db
	wlog.Debug(fmt.Sprintf("postgres: connection opened"))
	return nil
}

func (s *Store) Close() model.AppError {
	err := s.conn.Close()
	if err != nil {
		return model.NewInternalError("postgres.storage.close.disconnect.fail", err.Error())
	}
	s.conn = nil
	wlog.Debug(fmt.Sprintf("postgres: connection closed"))
	return nil
}
