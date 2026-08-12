package repo

import (
	"context"
	"time"

	"github.com/PaulAjii/physio-assistant-ai/internal/db"
	"github.com/PaulAjii/physio-assistant-ai/internal/models"
	"github.com/jackc/pgx/v5"
)

// CreateInvitationParams is the input for issuing an invitation. The raw token
// is generated and hashed by the auth layer; the repo only stores TokenHash.
// The clinic is passed separately as the tenant context (and written as
// clinic_id to satisfy the RLS WITH CHECK). InvitedBy is the issuing admin.
type CreateInvitationParams struct {
	Email     string
	Role      string
	TokenHash string
	InvitedBy *string
	ExpiresAt time.Time
}

// Invitations is the persistence contract for invitations. Lookup by hash is
// intentionally absent: accepting an invite happens before a tenant context
// exists, via the SECURITY DEFINER function in the auth layer.
type Invitations interface {
	Create(ctx context.Context, clinicID string, p CreateInvitationParams) (*models.Invitation, error)
	MarkAccepted(ctx context.Context, clinicID, id string) error
}

// InvitationRepo is the Postgres-backed implementation of Invitations.
type InvitationRepo struct {
	db *db.DB
}

// NewInvitationRepo returns an InvitationRepo backed by the given database.
func NewInvitationRepo(database *db.DB) *InvitationRepo {
	return &InvitationRepo{db: database}
}

// Compile-time check that InvitationRepo satisfies Invitations.
var _ Invitations = (*InvitationRepo)(nil)

const createInvitationSQL = `
INSERT INTO invitations (clinic_id, email, role, token_hash, invited_by, expires_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, clinic_id, email, role, token_hash, invited_by, expires_at, accepted_at, created_at`

// Create issues an invitation for clinicID. A second pending invite for the same
// email in the same clinic violates the partial unique index and surfaces as
// ErrConflict.
func (r *InvitationRepo) Create(ctx context.Context, clinicID string, p CreateInvitationParams) (*models.Invitation, error) {
	var inv models.Invitation
	err := r.db.WithTenant(ctx, clinicID, func(tx pgx.Tx) error {
		return tx.QueryRow(ctx, createInvitationSQL,
			clinicID, p.Email, p.Role, p.TokenHash, p.InvitedBy, p.ExpiresAt,
		).Scan(
			&inv.ID, &inv.ClinicID, &inv.Email, &inv.Role, &inv.TokenHash,
			&inv.InvitedBy, &inv.ExpiresAt, &inv.AcceptedAt, &inv.CreatedAt,
		)
	})
	if err != nil {
		if isUniqueViolation(err) {
			return nil, ErrConflict
		}
		return nil, err
	}
	return &inv, nil
}

const markInvitationAcceptedSQL = `
UPDATE invitations
SET accepted_at = now()
WHERE id = $1 AND accepted_at IS NULL`

// MarkAccepted stamps an invitation accepted. An id that does not exist in this
// clinic, or was already accepted, affects no rows and returns ErrNotFound so a
// double-accept is distinguishable from a fresh one.
func (r *InvitationRepo) MarkAccepted(ctx context.Context, clinicID, id string) error {
	return r.db.WithTenant(ctx, clinicID, func(tx pgx.Tx) error {
		tag, err := tx.Exec(ctx, markInvitationAcceptedSQL, id)
		if err != nil {
			return err
		}
		if tag.RowsAffected() == 0 {
			return ErrNotFound
		}
		return nil
	})
}
