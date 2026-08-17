// Package auth is backend-core's identity authority: it issues and verifies the
// JWTs the rest of the system trusts, and it owns the handful of database lookups
// that must happen before a tenant is known.
//
// Almost all data access in the app is tenant-scoped through db.WithTenant, which
// pins Row-Level Security to one clinic. Authentication is the exception: when a
// request arrives to log in, refresh a session, or accept an invitation, no
// clinic is established yet — that is precisely what these operations determine.
// This file is that boundary. It calls the SECURITY DEFINER functions from
// migration 00004 on the raw pool (RLS-exempt by design) and returns only the
// narrow column sets those functions expose, so no full tenant row ever crosses
// the boundary.
package auth

import (
	"context"
	"errors"
	"time"

	"github.com/PaulAjii/physio-assistant-ai/internal/db"
	"github.com/jackc/pgx/v5"
)

// errNotFound is returned by every lookup when no row matches. It is intentionally
// undifferentiated: for login in particular, the caller MUST fold it into a
// generic "invalid credentials" so a response never reveals whether an email is
// registered. Kept unexported because only this package acts on it.
var errNotFound = errors.New("auth: no matching row")

// loginCandidate is the projection auth_find_user_for_login returns: just enough
// to verify a password and, on success, establish identity and tenant. It is
// deliberately not models.User — PasswordHash must never ride along on a full
// entity, and these columns stop at the RLS boundary.
type loginCandidate struct {
	ID           string
	ClinicID     string
	Role         string
	FullName     string
	PasswordHash string
}

// refreshTokenRow is the projection auth_find_refresh_token returns. ExpiresAt and
// RevokedAt are surfaced (not pre-filtered) so the caller can tell an expired or
// already-revoked token apart from an unknown one — reuse of a revoked token is a
// security signal, not a silent miss.
type refreshTokenRow struct {
	ID        string
	UserID    string
	ClinicID  string
	ExpiresAt time.Time
	RevokedAt *time.Time
}

// invitationRow is the projection auth_find_invitation returns. Validity
// (AcceptedAt already set, ExpiresAt in the past) is judged by the caller.
type invitationRow struct {
	ID         string
	ClinicID   string
	Email      string
	Role       string
	ExpiresAt  time.Time
	AcceptedAt *time.Time
}

// lookups performs the three pre-tenant-context queries on the raw pool. Every
// other repository goes through db.WithTenant; this type is the single, audited
// exception, safe because the underlying functions run SECURITY DEFINER and
// return only narrow projections.
type lookups struct {
	db *db.DB
}

// newLookups constructs the lookup boundary over the given database.
func newLookups(database *db.DB) *lookups {
	return &lookups{db: database}
}

const findUserForLoginSQL = `
SELECT id, clinic_id, role, full_name, password_hash
FROM auth_find_user_for_login($1::text::citext)`

// findUserForLogin resolves an email to the one live user (in a live clinic) whose
// password can be verified. The email is bound as text and cast to citext in SQL
// so we never depend on a pgx codec for the citext extension type; the function
// does the case-insensitive match. A missing user is errNotFound.
func (l *lookups) findUserForLogin(ctx context.Context, email string) (*loginCandidate, error) {
	var u loginCandidate
	err := l.db.Pool().QueryRow(ctx, findUserForLoginSQL, email).Scan(
		&u.ID, &u.ClinicID, &u.Role, &u.FullName, &u.PasswordHash,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

const findRefreshTokenSQL = `
SELECT id, user_id, clinic_id, expires_at, revoked_at
FROM auth_find_refresh_token($1)`

// findRefreshToken resolves a token hash to its stored row. The raw token is never
// stored, so the caller hashes first and passes only the hash here.
func (l *lookups) findRefreshToken(ctx context.Context, tokenHash string) (*refreshTokenRow, error) {
	var t refreshTokenRow
	err := l.db.Pool().QueryRow(ctx, findRefreshTokenSQL, tokenHash).Scan(
		&t.ID, &t.UserID, &t.ClinicID, &t.ExpiresAt, &t.RevokedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errNotFound
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

const findInvitationSQL = `
SELECT id, clinic_id, email::text, role, expires_at, accepted_at
FROM auth_find_invitation($1)`

// findInvitation resolves an invitation token hash to its row so an
// unauthenticated invitee can accept it. email is cast to text in SQL for the same
// codec-independence reason as login.
func (l *lookups) findInvitation(ctx context.Context, tokenHash string) (*invitationRow, error) {
	var inv invitationRow
	err := l.db.Pool().QueryRow(ctx, findInvitationSQL, tokenHash).Scan(
		&inv.ID, &inv.ClinicID, &inv.Email, &inv.Role, &inv.ExpiresAt, &inv.AcceptedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errNotFound
	}
	if err != nil {
		return nil, err
	}
	return &inv, nil
}
