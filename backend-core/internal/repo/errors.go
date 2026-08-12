// Package repo holds the persistence layer. Each entity has its own repository
// struct (UserRepo, SessionRepo, ...) taking a *db.DB, so they can be tested and
// mocked independently. This file holds only what all of them share.
package repo

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

// Shared sentinel errors, so handlers can branch on intent (map to 404 / 409)
// without importing pgx or knowing SQL details.
var (
	// ErrNotFound is returned when a query matches no rows. Under RLS a row that
	// belongs to another tenant is indistinguishable from one that does not
	// exist — both surface as ErrNotFound, which is exactly what we want.
	ErrNotFound = errors.New("not found")

	// ErrConflict is returned when a write violates a unique constraint, such as
	// registering an email that is already taken.
	ErrConflict = errors.New("conflict")
)

// isUniqueViolation reports whether err is a Postgres unique-constraint
// violation (SQLSTATE 23505). Repositories use it to translate a duplicate-key
// database error into ErrConflict.
func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
