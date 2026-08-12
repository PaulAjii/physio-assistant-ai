package repo

import (
	"context"
	"time"

	"github.com/PaulAjii/physio-assistant-ai/internal/db"
	"github.com/PaulAjii/physio-assistant-ai/internal/models"
	"github.com/jackc/pgx/v5"
)

// CreateRefreshTokenParams is the input for storing a freshly issued token. The
// raw token is generated and hashed by the auth layer; the repo only ever sees
// TokenHash. UserID identifies the owner; the clinic is passed separately as the
// tenant context (and written as clinic_id to satisfy the RLS WITH CHECK).
type CreateRefreshTokenParams struct {
	UserID    string
	TokenHash string
	ExpiresAt time.Time
}

// RefreshTokens is the persistence contract for refresh tokens. Lookup by hash
// is intentionally absent: it happens at refresh time, before a tenant context
// exists, via the SECURITY DEFINER function in the auth layer.
type RefreshTokens interface {
	Create(ctx context.Context, clinicID string, p CreateRefreshTokenParams) (*models.RefreshToken, error)
	Revoke(ctx context.Context, clinicID, id string) error
}

// RefreshTokenRepo is the Postgres-backed implementation of RefreshTokens.
type RefreshTokenRepo struct {
	db *db.DB
}

// NewRefreshTokenRepo returns a RefreshTokenRepo backed by the given database.
func NewRefreshTokenRepo(database *db.DB) *RefreshTokenRepo {
	return &RefreshTokenRepo{db: database}
}

// Compile-time check that RefreshTokenRepo satisfies RefreshTokens.
var _ RefreshTokens = (*RefreshTokenRepo)(nil)

const createRefreshTokenSQL = `
INSERT INTO refresh_tokens (user_id, clinic_id, token_hash, expires_at)
VALUES ($1, $2, $3, $4)
RETURNING id, user_id, clinic_id, token_hash, expires_at, revoked_at, created_at`

// Create stores a new refresh token for a user in clinicID.
func (r *RefreshTokenRepo) Create(ctx context.Context, clinicID string, p CreateRefreshTokenParams) (*models.RefreshToken, error) {
	var t models.RefreshToken
	err := r.db.WithTenant(ctx, clinicID, func(tx pgx.Tx) error {
		return tx.QueryRow(ctx, createRefreshTokenSQL,
			p.UserID, clinicID, p.TokenHash, p.ExpiresAt,
		).Scan(
			&t.ID, &t.UserID, &t.ClinicID, &t.TokenHash,
			&t.ExpiresAt, &t.RevokedAt, &t.CreatedAt,
		)
	})
	if err != nil {
		return nil, err
	}
	return &t, nil
}

const revokeRefreshTokenSQL = `
UPDATE refresh_tokens
SET revoked_at = now()
WHERE id = $1 AND revoked_at IS NULL`

// Revoke marks a token revoked (on rotation or logout). Revoking an id that does
// not exist in this clinic, or is already revoked, affects no rows and returns
// ErrNotFound so the caller can tell a no-op from a success.
func (r *RefreshTokenRepo) Revoke(ctx context.Context, clinicID, id string) error {
	return r.db.WithTenant(ctx, clinicID, func(tx pgx.Tx) error {
		tag, err := tx.Exec(ctx, revokeRefreshTokenSQL, id)
		if err != nil {
			return err
		}
		if tag.RowsAffected() == 0 {
			return ErrNotFound
		}
		return nil
	})
}
