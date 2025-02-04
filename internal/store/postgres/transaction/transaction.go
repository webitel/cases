package transaction

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Transaction defines the interface for handling transactions.
type Transaction interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
}

// TxManager provides a concrete implementation of the Transaction interface.
type TxManager struct {
	tx pgx.Tx
}

// NewTxManager creates a new transaction manager from an existing pgx.Tx.
func NewTxManager(tx pgx.Tx) *TxManager {
	return &TxManager{tx: tx}
}

// Commit commits the current transaction.
func (m *TxManager) Commit(ctx context.Context) error {
	if err := m.tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

// Rollback rolls back the current transaction.
func (m *TxManager) Rollback(ctx context.Context) error {
	if err := m.tx.Rollback(ctx); err != nil {
		return err
	}
	return nil
}

// QueryRow executes a query that is expected to return at most one row.
func (m *TxManager) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	// Execute the query and return a single row result
	return m.tx.QueryRow(ctx, sql, args...)
}

// Query executes a query that returns multiple rows.
func (m *TxManager) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	// Execute the query and return the rows result set
	rows, err := m.tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// Exec executes a query that doesn't return rows and returns a CommandTag.
func (m *TxManager) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	// Execute the command and return the CommandTag and any error encountered
	commandTag, err := m.tx.Exec(ctx, sql, args...)
	if err != nil {
		return pgconn.CommandTag{}, err
	}
	return commandTag, nil
}
