// Package db owns the Postgres connection pool and the one pattern the whole
// data layer depends on: running tenant queries inside a transaction that has
// the RLS context set.
package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DB wraps a pgx connection pool. backend-core connects as the runtime role
// (physio_app), which is subject to Row-Level Security, so a query only ever
// sees a tenant's rows when the tenant context has been set for the surrounding
// transaction — see WithTenant.
type DB struct {
	pool *pgxpool.Pool
}

// New opens the pool from dsn and verifies the connection with a ping. The
// caller owns the returned *DB and must Close it on shutdown.
func New(ctx context.Context, dsn string) (*DB, error) {
	if dsn == "" {
		return nil, errors.New("database DSN is empty (set DATABASE_URL)")
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("connect to postgres: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return &DB{pool: pool}, nil
}

// Close releases every connection in the pool.
func (d *DB) Close() {
	d.pool.Close()
}

// Pool exposes the raw pool for the few queries that intentionally run WITHOUT a
// tenant context: the SECURITY DEFINER auth bootstrap functions (login, refresh,
// invite lookups), which bypass RLS by design because they run before a clinic
// is known. Everything that touches tenant data must go through WithTenant
// instead. Note that this is fail-safe either way: a plain tenant query on the
// pool with no context set matches zero rows rather than leaking another clinic.
func (d *DB) Pool() *pgxpool.Pool {
	return d.pool
}

// WithTenant runs fn inside a single transaction whose RLS context is pinned to
// clinicID, then commits (or rolls back if fn returns an error).
//
// set_config's third argument is true (is_local), so the setting is scoped to
// THIS transaction only. That is essential with a pooled connection: a plain
// SET would persist on the connection and leak this request's tenant into the
// next request that reuses it. When the transaction ends, the setting is gone.
//
// An empty clinicID sets the context to NULL, which makes current_clinic_id()
// return NULL and every policy match no rows — fail closed, never open.
func (d *DB) WithTenant(ctx context.Context, clinicID string, fn func(pgx.Tx) error) error {
	tx, err := d.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	// Safe to call after Commit: pgx makes a second completion a no-op, so a
	// committed transaction is not rolled back here.
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, "SELECT set_config('app.current_clinic_id', $1, true)", clinicID); err != nil {
		return fmt.Errorf("set tenant context: %w", err)
	}

	if err := fn(tx); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}
