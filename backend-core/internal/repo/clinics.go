package repo

import (
	"context"
	"errors"

	"github.com/PaulAjii/physio-assistant-ai/internal/db"
	"github.com/PaulAjii/physio-assistant-ai/internal/models"
	"github.com/jackc/pgx/v5"
)

// Clinics is the persistence contract for clinics. There is no Create: a clinic
// is created by the owner-privileged seed step, because as the runtime role an
// INSERT would fail the RLS WITH CHECK (the new clinic is not yet the current
// tenant context).
type Clinics interface {
	GetByID(ctx context.Context, id string) (*models.Clinic, error)
}

// ClinicRepo is the Postgres-backed implementation of Clinics.
type ClinicRepo struct {
	db *db.DB
}

// NewClinicRepo returns a ClinicRepo backed by the given database.
func NewClinicRepo(database *db.DB) *ClinicRepo {
	return &ClinicRepo{db: database}
}

// Compile-time check that ClinicRepo satisfies Clinics.
var _ Clinics = (*ClinicRepo)(nil)

const getClinicByIDSQL = `
SELECT id, name, created_at, updated_at
FROM clinics
WHERE id = $1 AND deleted_at IS NULL`

// GetByID reads a clinic. The id doubles as the tenant context, and the RLS
// policy on clinics is id = current_clinic_id(), so a clinic can only ever read
// itself — asking for any other id returns ErrNotFound.
func (r *ClinicRepo) GetByID(ctx context.Context, id string) (*models.Clinic, error) {
	var c models.Clinic
	err := r.db.WithTenant(ctx, id, func(tx pgx.Tx) error {
		return tx.QueryRow(ctx, getClinicByIDSQL, id).Scan(
			&c.ID, &c.Name, &c.CreatedAt, &c.UpdatedAt,
		)
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}
