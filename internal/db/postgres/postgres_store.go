package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/webitel/cases/internal/db"
	lookup2 "github.com/webitel/cases/internal/db/postgres/lookup"
	"github.com/webitel/cases/model"
	"github.com/webitel/wlog"
)

type PostgresStore struct {
	config                 *model.DatabaseConfig
	conn                   *sqlx.DB
	appealLookupStore      db.AppealLookupStore
	statusLookupStore      db.StatusLookupStore
	closeReasonLookupStore db.CloseReasonLookupStore
}

func New(config *model.DatabaseConfig) *PostgresStore {
	return &PostgresStore{config: config}
}

func (s *PostgresStore) Appeal() db.AppealLookupStore {
	if s.appealLookupStore == nil {
		log, err := lookup2.NewAppealLookupStore(s)
		if err != nil {
			return nil
		}
		s.appealLookupStore = log
	}
	return s.appealLookupStore
}

func (s *PostgresStore) CloseReason() db.CloseReasonLookupStore {
	if s.closeReasonLookupStore == nil {
		log, err := lookup2.NewCloseReasonLookupStore(s)
		if err != nil {
			return nil
		}
		s.closeReasonLookupStore = log
	}
	return s.closeReasonLookupStore
}
func (s *PostgresStore) Status() db.StatusLookupStore {
	if s.statusLookupStore == nil {
		log, err := lookup2.NewStatusLookupStore(s)
		if err != nil {
			return nil
		}
		s.statusLookupStore = log
	}
	return s.statusLookupStore
}

func (s *PostgresStore) Database() (*sqlx.DB, model.AppError) {
	if s.conn == nil {
		model.NewInternalError("postgres.db.database.check.bad_arguments", "database connection is not opened")
	}
	return s.conn, nil
}

func (s *PostgresStore) Open() model.AppError {
	db, err := sqlx.Connect("pgx", s.config.Url)
	if err != nil {
		return model.NewInternalError("postgres.db.open.connect.fail", err.Error())
	}
	s.conn = db
	wlog.Debug(fmt.Sprintf("postgres: connection opened"))
	return nil
}

func (s *PostgresStore) Close() model.AppError {
	err := s.conn.Close()
	if err != nil {
		return model.NewInternalError("postgres.db.close.disconnect.fail", err.Error())
	}
	s.conn = nil
	wlog.Debug(fmt.Sprintf("postgres: connection closed"))
	return nil
}
