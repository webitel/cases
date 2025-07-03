package postgres

import (
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/store"
	"regexp"
	"sync"
)

func ParseError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return errors.Wrap(store.ErrNoRows, errors.WithCause(err))
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.UniqueViolation:
			return errors.Wrap(store.ErrUniqueViolation, errors.WithCause(err),
				errors.WithValue("entity", findColumn(pgErr.Detail)+" = "+findValue(pgErr.Detail)),
			)
		case pgerrcode.ForeignKeyViolation:
			msg := "value is still referenced by the parent table"
			if findForeignKeyTable(pgErr.Detail) != "" {
				msg = "value isn't present in the parent table"
			}

			return errors.Wrap(store.ErrForeignKeyViolation, errors.WithCause(err), errors.AppendMessage(msg),
				errors.WithValue("value", findColumn(pgErr.Detail)+" = "+findValue(pgErr.Detail)),
				errors.WithValue("foreign_table", findForeignKeyTable(pgErr.Detail)),
			)
		case pgerrcode.CheckViolation:
			return errors.Wrap(store.ErrCheckViolation, errors.WithCause(err),
				errors.AppendMessage(checkViolationErrorRegistry[pgErr.ConstraintName]),
				errors.WithValue("constraint", pgErr.ConstraintName),
			)
		case pgerrcode.NotNullViolation:
			return errors.Wrap(store.ErrNotNullViolation, errors.WithCause(err),
				errors.WithValue("column", pgErr.TableName+"."+pgErr.ColumnName),
			)
		}
	}

	return errors.Wrap(store.ErrInternal, errors.WithCause(err))
}

var checkViolationErrorRegistry = map[string]string{}
var constraintMu sync.RWMutex

// RegisterConstraint register custom database check constraint (like "CHECK
// balance > 0").
// Postgres doesn't define a very useful message for constraint
// failures (new row for relation "accounts" violates check constraint), so you
// can define your own.
//   - name - should be the name of the constraint in the database.
//   - message - your own custom error message
//
// Panics if you attempt to register two constraints with the same name.
func RegisterConstraint(name, message string) {
	constraintMu.Lock()
	defer constraintMu.Unlock()
	if _, dup := checkViolationErrorRegistry[name]; dup {
		panic("register constraint called twice for name " + name)
	}

	checkViolationErrorRegistry[name] = message
}

var columnFinder = regexp.MustCompile(`Key \((.+)\)=`)

// findColumn finds the column in the given pq Detail error string. If the
// column does not exist, the empty string is returned.
// Detail can look like this:
//
//	Key (id)=(3c7d2b4a-3fc8-4782-a518-4ce9efef51e7) already exists.
func findColumn(detail string) string {
	results := columnFinder.FindStringSubmatch(detail)
	if len(results) < 2 {
		return ""
	} else {
		return results[1]
	}
}

var valueFinder = regexp.MustCompile(`Key \(.+\)=\((.+)\)`)

// findColumn finds the column in the given pq Detail error string.
// If the column does not exist, the empty string is returned.
// Detail can look like this:
//
//	Key (id)=(3c7d2b4a-3fc8-4782-a518-4ce9efef51e7) already exists.
func findValue(detail string) string {
	results := valueFinder.FindStringSubmatch(detail)
	if len(results) < 2 {
		return ""
	}

	return results[1]
}

var foreignKeyFinder = regexp.MustCompile(`not present in table "(.+)"`)

// findForeignKeyTable finds the referenced table in the given pq Detail error
// string. If we can't find the table, we return the empty string.
// Detail can look like this:
//
//	Key (account_id)=(91f47e99-d616-4d8c-9c02-cbd13bceac60) is not present in table "accounts"
func findForeignKeyTable(detail string) string {
	results := foreignKeyFinder.FindStringSubmatch(detail)
	if len(results) < 2 {
		return ""
	}
	return results[1]
}
