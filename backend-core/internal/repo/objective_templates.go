package repo

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/PaulAjii/physio-assistant-ai/internal/db"
	"github.com/PaulAjii/physio-assistant-ai/internal/models"
	"github.com/jackc/pgx/v5"
)

// CreateTemplateParams is the input for defining an objective-exam scaffold. Body
// is the ordered JSON array of field definitions; leave it nil to accept the
// column default ('[]').
type CreateTemplateParams struct {
	Name string
	Body json.RawMessage
}

// UpdateTemplateParams edits a template. IsActive lets a clinic retire a scaffold
// without deleting it, so historical sessions that reference it keep their shape.
type UpdateTemplateParams struct {
	Name     string
	Body     json.RawMessage
	IsActive bool
}

// ObjectiveTemplates is the persistence contract for clinic-configurable exam
// scaffolds.
type ObjectiveTemplates interface {
	Create(ctx context.Context, clinicID string, p CreateTemplateParams) (*models.ObjectiveTemplate, error)
	GetByID(ctx context.Context, clinicID, id string) (*models.ObjectiveTemplate, error)
	List(ctx context.Context, clinicID string) ([]models.ObjectiveTemplate, error)
	Update(ctx context.Context, clinicID, id string, p UpdateTemplateParams) (*models.ObjectiveTemplate, error)
}

// ObjectiveTemplateRepo is the Postgres-backed implementation of ObjectiveTemplates.
type ObjectiveTemplateRepo struct {
	db *db.DB
}

// NewObjectiveTemplateRepo returns an ObjectiveTemplateRepo backed by the given database.
func NewObjectiveTemplateRepo(database *db.DB) *ObjectiveTemplateRepo {
	return &ObjectiveTemplateRepo{db: database}
}

// Compile-time check that ObjectiveTemplateRepo satisfies ObjectiveTemplates.
var _ ObjectiveTemplates = (*ObjectiveTemplateRepo)(nil)

const templateColumns = `id, clinic_id, name, body, is_active, created_at, updated_at`

func scanTemplate(row pgx.Row, t *models.ObjectiveTemplate) error {
	return row.Scan(
		&t.ID, &t.ClinicID, &t.Name, &t.Body, &t.IsActive, &t.CreatedAt, &t.UpdatedAt,
	)
}

const createTemplateSQL = `
INSERT INTO objective_templates (clinic_id, name, body)
VALUES ($1, $2, COALESCE($3, '[]'::jsonb))
RETURNING ` + templateColumns

// Create defines a new template. A second live template with the same name in the
// clinic violates the partial unique index and surfaces as ErrConflict.
func (r *ObjectiveTemplateRepo) Create(ctx context.Context, clinicID string, p CreateTemplateParams) (*models.ObjectiveTemplate, error) {
	var t models.ObjectiveTemplate
	err := r.db.WithTenant(ctx, clinicID, func(tx pgx.Tx) error {
		return scanTemplate(tx.QueryRow(ctx, createTemplateSQL, clinicID, p.Name, p.Body), &t)
	})
	if err != nil {
		if isUniqueViolation(err) {
			return nil, ErrConflict
		}
		return nil, err
	}
	return &t, nil
}

const getTemplateByIDSQL = `
SELECT ` + templateColumns + `
FROM objective_templates
WHERE id = $1 AND deleted_at IS NULL`

// GetByID fetches a live template within clinicID.
func (r *ObjectiveTemplateRepo) GetByID(ctx context.Context, clinicID, id string) (*models.ObjectiveTemplate, error) {
	var t models.ObjectiveTemplate
	err := r.db.WithTenant(ctx, clinicID, func(tx pgx.Tx) error {
		return scanTemplate(tx.QueryRow(ctx, getTemplateByIDSQL, id), &t)
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

const listTemplatesSQL = `
SELECT ` + templateColumns + `
FROM objective_templates
WHERE deleted_at IS NULL
ORDER BY name`

// List returns all live templates for the clinic, ordered by name. Retired
// templates (is_active = false) are included so an admin UI can manage them; a
// caller building a pick-list should filter on IsActive. RLS scopes the read to
// the current tenant, so no clinic_id predicate is needed.
func (r *ObjectiveTemplateRepo) List(ctx context.Context, clinicID string) ([]models.ObjectiveTemplate, error) {
	var list []models.ObjectiveTemplate
	err := r.db.WithTenant(ctx, clinicID, func(tx pgx.Tx) error {
		rows, err := tx.Query(ctx, listTemplatesSQL)
		if err != nil {
			return err
		}
		list, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (models.ObjectiveTemplate, error) {
			var t models.ObjectiveTemplate
			err := scanTemplate(row, &t)
			return t, err
		})
		return err
	})
	if err != nil {
		return nil, err
	}
	return list, nil
}

const updateTemplateSQL = `
UPDATE objective_templates
SET name = $2, body = COALESCE($3, '[]'::jsonb), is_active = $4
WHERE id = $1 AND deleted_at IS NULL
RETURNING ` + templateColumns

// Update edits a template. Renaming onto another live template's name in the same
// clinic surfaces as ErrConflict; a missing id returns ErrNotFound.
func (r *ObjectiveTemplateRepo) Update(ctx context.Context, clinicID, id string, p UpdateTemplateParams) (*models.ObjectiveTemplate, error) {
	var t models.ObjectiveTemplate
	err := r.db.WithTenant(ctx, clinicID, func(tx pgx.Tx) error {
		return scanTemplate(tx.QueryRow(ctx, updateTemplateSQL, id, p.Name, p.Body, p.IsActive), &t)
	})
	if err != nil {
		if isUniqueViolation(err) {
			return nil, ErrConflict
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &t, nil
}
