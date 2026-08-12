package repo

import (
	"context"
	"errors"

	"github.com/PaulAjii/physio-assistant-ai/internal/db"
	"github.com/PaulAjii/physio-assistant-ai/internal/models"
	"github.com/jackc/pgx/v5"
)

// CreateSessionParams is the input for opening a session. The clinic is passed
// separately (tenant context + clinic_id). ClinicianID is the authenticated user
// who owns the visit; PatientRef and Language are optional.
type CreateSessionParams struct {
	ClinicianID string
	PatientRef  *string
	Language    *string
}

// Sessions is the persistence contract for consultation sessions.
type Sessions interface {
	Create(ctx context.Context, clinicID string, p CreateSessionParams) (*models.Session, error)
	GetByID(ctx context.Context, clinicID, id string) (*models.Session, error)
	Finalize(ctx context.Context, clinicID, id, signedBy string) (*models.Session, error)
}

// SessionRepo is the Postgres-backed implementation of Sessions.
type SessionRepo struct {
	db *db.DB
}

// NewSessionRepo returns a SessionRepo backed by the given database.
func NewSessionRepo(database *db.DB) *SessionRepo {
	return &SessionRepo{db: database}
}

// Compile-time check that SessionRepo satisfies Sessions.
var _ Sessions = (*SessionRepo)(nil)

// The full column list is shared by every read/write so the scan target lines up
// exactly with what the query returns.
const sessionColumns = `id, clinic_id, clinician_id, patient_ref, status, language,
	started_at, finalized_at, signed_by, signed_at, created_at, updated_at`

func scanSession(row pgx.Row, s *models.Session) error {
	return row.Scan(
		&s.ID, &s.ClinicID, &s.ClinicianID, &s.PatientRef, &s.Status, &s.Language,
		&s.StartedAt, &s.FinalizedAt, &s.SignedBy, &s.SignedAt, &s.CreatedAt, &s.UpdatedAt,
	)
}

const createSessionSQL = `
INSERT INTO consultation_sessions (clinic_id, clinician_id, patient_ref, language)
VALUES ($1, $2, $3, $4)
RETURNING ` + sessionColumns

// Create opens a draft session for the clinician in clinicID.
func (r *SessionRepo) Create(ctx context.Context, clinicID string, p CreateSessionParams) (*models.Session, error) {
	var s models.Session
	err := r.db.WithTenant(ctx, clinicID, func(tx pgx.Tx) error {
		return scanSession(tx.QueryRow(ctx, createSessionSQL,
			clinicID, p.ClinicianID, p.PatientRef, p.Language,
		), &s)
	})
	if err != nil {
		return nil, err
	}
	return &s, nil
}

const getSessionByIDSQL = `
SELECT ` + sessionColumns + `
FROM consultation_sessions
WHERE id = $1 AND deleted_at IS NULL`

// GetByID fetches a live session within clinicID; RLS also confines the read to
// that tenant, so an id from another clinic returns ErrNotFound.
func (r *SessionRepo) GetByID(ctx context.Context, clinicID, id string) (*models.Session, error) {
	var s models.Session
	err := r.db.WithTenant(ctx, clinicID, func(tx pgx.Tx) error {
		return scanSession(tx.QueryRow(ctx, getSessionByIDSQL, id), &s)
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

const finalizeSessionSQL = `
UPDATE consultation_sessions
SET status = 'finalized', finalized_at = now(), signed_by = $2, signed_at = now()
WHERE id = $1
RETURNING ` + sessionColumns

// Finalize signs and closes a session. signedBy is the authenticated clinician,
// which the write records in signed_by (the database CHECK then guarantees a
// finalized session is always signed). It runs as select-then-update in one tx so
// the errors are precise: a missing session is ErrNotFound, an already-finalized
// one is ErrConflict, and only a genuinely open session is stamped.
func (r *SessionRepo) Finalize(ctx context.Context, clinicID, id, signedBy string) (*models.Session, error) {
	var s models.Session
	err := r.db.WithTenant(ctx, clinicID, func(tx pgx.Tx) error {
		var status string
		err := tx.QueryRow(ctx,
			`SELECT status FROM consultation_sessions WHERE id = $1 AND deleted_at IS NULL`, id,
		).Scan(&status)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		if err != nil {
			return err
		}
		if status == models.SessionStatusFinalized {
			return ErrConflict
		}
		return scanSession(tx.QueryRow(ctx, finalizeSessionSQL, id, signedBy), &s)
	})
	if err != nil {
		return nil, err
	}
	return &s, nil
}
