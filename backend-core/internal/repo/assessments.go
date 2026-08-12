package repo

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/PaulAjii/physio-assistant-ai/internal/db"
	"github.com/PaulAjii/physio-assistant-ai/internal/models"
	"github.com/jackc/pgx/v5"
)

// CreateAssessmentParams is the input for the first save of a session's note.
// Subjective and Objective are raw JSON sections; leave them nil to accept the
// column defaults ('{}'). AISummary and PainIntensity are optional.
type CreateAssessmentParams struct {
	SessionID     string
	Subjective    json.RawMessage
	Objective     json.RawMessage
	AISummary     *string
	PainIntensity *int
}

// UpdateAssessmentParams is the input for editing the note. It replaces the
// editable content wholesale; the repo bumps version so each edit is visible.
type UpdateAssessmentParams struct {
	Subjective    json.RawMessage
	Objective     json.RawMessage
	AISummary     *string
	PainIntensity *int
}

// Assessments is the persistence contract for the structured note. There is one
// live assessment per session, so it is fetched by session rather than by its own
// id in the normal flow.
type Assessments interface {
	Create(ctx context.Context, clinicID string, p CreateAssessmentParams) (*models.AssessmentRecord, error)
	GetBySessionID(ctx context.Context, clinicID, sessionID string) (*models.AssessmentRecord, error)
	Update(ctx context.Context, clinicID, id string, p UpdateAssessmentParams) (*models.AssessmentRecord, error)
}

// AssessmentRepo is the Postgres-backed implementation of Assessments.
type AssessmentRepo struct {
	db *db.DB
}

// NewAssessmentRepo returns an AssessmentRepo backed by the given database.
func NewAssessmentRepo(database *db.DB) *AssessmentRepo {
	return &AssessmentRepo{db: database}
}

// Compile-time check that AssessmentRepo satisfies Assessments.
var _ Assessments = (*AssessmentRepo)(nil)

const assessmentColumns = `id, clinic_id, session_id, subjective, objective,
	ai_summary, pain_intensity, version, created_at, updated_at`

func scanAssessment(row pgx.Row, a *models.AssessmentRecord) error {
	return row.Scan(
		&a.ID, &a.ClinicID, &a.SessionID, &a.Subjective, &a.Objective,
		&a.AISummary, &a.PainIntensity, &a.Version, &a.CreatedAt, &a.UpdatedAt,
	)
}

const createAssessmentSQL = `
INSERT INTO assessments (clinic_id, session_id, subjective, objective, ai_summary, pain_intensity)
VALUES ($1, $2, COALESCE($3, '{}'::jsonb), COALESCE($4, '{}'::jsonb), $5, $6)
RETURNING ` + assessmentColumns

// Create writes the first assessment for a session. A second live assessment for
// the same session violates the partial unique index and surfaces as ErrConflict.
func (r *AssessmentRepo) Create(ctx context.Context, clinicID string, p CreateAssessmentParams) (*models.AssessmentRecord, error) {
	var a models.AssessmentRecord
	err := r.db.WithTenant(ctx, clinicID, func(tx pgx.Tx) error {
		return scanAssessment(tx.QueryRow(ctx, createAssessmentSQL,
			clinicID, p.SessionID, p.Subjective, p.Objective, p.AISummary, p.PainIntensity,
		), &a)
	})
	if err != nil {
		if isUniqueViolation(err) {
			return nil, ErrConflict
		}
		return nil, err
	}
	return &a, nil
}

const getAssessmentBySessionSQL = `
SELECT ` + assessmentColumns + `
FROM assessments
WHERE session_id = $1 AND deleted_at IS NULL`

// GetBySessionID returns the live assessment for a session, or ErrNotFound if the
// session has none yet.
func (r *AssessmentRepo) GetBySessionID(ctx context.Context, clinicID, sessionID string) (*models.AssessmentRecord, error) {
	var a models.AssessmentRecord
	err := r.db.WithTenant(ctx, clinicID, func(tx pgx.Tx) error {
		return scanAssessment(tx.QueryRow(ctx, getAssessmentBySessionSQL, sessionID), &a)
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}

const updateAssessmentSQL = `
UPDATE assessments
SET subjective = COALESCE($2, '{}'::jsonb),
    objective = COALESCE($3, '{}'::jsonb),
    ai_summary = $4,
    pain_intensity = $5,
    version = version + 1
WHERE id = $1 AND deleted_at IS NULL
RETURNING ` + assessmentColumns

// Update replaces the editable content and bumps version. An id that is missing
// or belongs to another tenant returns ErrNotFound.
func (r *AssessmentRepo) Update(ctx context.Context, clinicID, id string, p UpdateAssessmentParams) (*models.AssessmentRecord, error) {
	var a models.AssessmentRecord
	err := r.db.WithTenant(ctx, clinicID, func(tx pgx.Tx) error {
		return scanAssessment(tx.QueryRow(ctx, updateAssessmentSQL,
			id, p.Subjective, p.Objective, p.AISummary, p.PainIntensity,
		), &a)
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}
