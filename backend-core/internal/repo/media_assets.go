package repo

import (
	"context"
	"errors"

	"github.com/PaulAjii/physio-assistant-ai/internal/db"
	"github.com/PaulAjii/physio-assistant-ai/internal/models"
	"github.com/jackc/pgx/v5"
)

// CreateMediaAssetParams reserves a storage key for a recording before the blob
// is uploaded. Kind is models.MediaKind*; ContentType is optional at this stage.
// The asset starts in the pending state.
type CreateMediaAssetParams struct {
	SessionID   string
	Kind        string
	StorageKey  string
	ContentType *string
}

// MediaAssets is the persistence contract for recording metadata. The blobs live
// in object storage; these rows only track the key and upload lifecycle.
type MediaAssets interface {
	Create(ctx context.Context, clinicID string, p CreateMediaAssetParams) (*models.MediaAsset, error)
	GetByID(ctx context.Context, clinicID, id string) (*models.MediaAsset, error)
	ListBySession(ctx context.Context, clinicID, sessionID string) ([]models.MediaAsset, error)
	MarkUploaded(ctx context.Context, clinicID, id string, contentType *string, bytes int64) (*models.MediaAsset, error)
	MarkFailed(ctx context.Context, clinicID, id string) (*models.MediaAsset, error)
}

// MediaAssetRepo is the Postgres-backed implementation of MediaAssets.
type MediaAssetRepo struct {
	db *db.DB
}

// NewMediaAssetRepo returns a MediaAssetRepo backed by the given database.
func NewMediaAssetRepo(database *db.DB) *MediaAssetRepo {
	return &MediaAssetRepo{db: database}
}

// Compile-time check that MediaAssetRepo satisfies MediaAssets.
var _ MediaAssets = (*MediaAssetRepo)(nil)

const mediaColumns = `id, clinic_id, session_id, kind, storage_key, content_type,
	bytes, status, created_at, updated_at`

func scanMediaAsset(row pgx.Row, m *models.MediaAsset) error {
	return row.Scan(
		&m.ID, &m.ClinicID, &m.SessionID, &m.Kind, &m.StorageKey, &m.ContentType,
		&m.Bytes, &m.Status, &m.CreatedAt, &m.UpdatedAt,
	)
}

const createMediaAssetSQL = `
INSERT INTO media_assets (clinic_id, session_id, kind, storage_key, content_type)
VALUES ($1, $2, $3, $4, $5)
RETURNING ` + mediaColumns

// Create reserves a pending asset. A reused storage_key hits the global unique
// index and surfaces as ErrConflict.
func (r *MediaAssetRepo) Create(ctx context.Context, clinicID string, p CreateMediaAssetParams) (*models.MediaAsset, error) {
	var m models.MediaAsset
	err := r.db.WithTenant(ctx, clinicID, func(tx pgx.Tx) error {
		return scanMediaAsset(tx.QueryRow(ctx, createMediaAssetSQL,
			clinicID, p.SessionID, p.Kind, p.StorageKey, p.ContentType,
		), &m)
	})
	if err != nil {
		if isUniqueViolation(err) {
			return nil, ErrConflict
		}
		return nil, err
	}
	return &m, nil
}

const getMediaAssetByIDSQL = `
SELECT ` + mediaColumns + `
FROM media_assets
WHERE id = $1 AND deleted_at IS NULL`

// GetByID fetches a live asset within clinicID.
func (r *MediaAssetRepo) GetByID(ctx context.Context, clinicID, id string) (*models.MediaAsset, error) {
	var m models.MediaAsset
	err := r.db.WithTenant(ctx, clinicID, func(tx pgx.Tx) error {
		return scanMediaAsset(tx.QueryRow(ctx, getMediaAssetByIDSQL, id), &m)
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}

const listMediaBySessionSQL = `
SELECT ` + mediaColumns + `
FROM media_assets
WHERE session_id = $1 AND deleted_at IS NULL
ORDER BY created_at`

// ListBySession returns a session's live assets in upload order.
func (r *MediaAssetRepo) ListBySession(ctx context.Context, clinicID, sessionID string) ([]models.MediaAsset, error) {
	var list []models.MediaAsset
	err := r.db.WithTenant(ctx, clinicID, func(tx pgx.Tx) error {
		rows, err := tx.Query(ctx, listMediaBySessionSQL, sessionID)
		if err != nil {
			return err
		}
		list, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (models.MediaAsset, error) {
			var m models.MediaAsset
			err := scanMediaAsset(row, &m)
			return m, err
		})
		return err
	})
	if err != nil {
		return nil, err
	}
	return list, nil
}

const markMediaUploadedSQL = `
UPDATE media_assets
SET status = 'uploaded', content_type = COALESCE($2, content_type), bytes = $3
WHERE id = $1 AND deleted_at IS NULL
RETURNING ` + mediaColumns

// MarkUploaded confirms a completed upload, recording the final size and (if the
// caller learned it) the content type. A missing id returns ErrNotFound.
func (r *MediaAssetRepo) MarkUploaded(ctx context.Context, clinicID, id string, contentType *string, bytes int64) (*models.MediaAsset, error) {
	var m models.MediaAsset
	err := r.db.WithTenant(ctx, clinicID, func(tx pgx.Tx) error {
		return scanMediaAsset(tx.QueryRow(ctx, markMediaUploadedSQL, id, contentType, bytes), &m)
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}

const markMediaFailedSQL = `
UPDATE media_assets
SET status = 'failed'
WHERE id = $1 AND deleted_at IS NULL
RETURNING ` + mediaColumns

// MarkFailed records that an upload did not complete. A missing id returns
// ErrNotFound.
func (r *MediaAssetRepo) MarkFailed(ctx context.Context, clinicID, id string) (*models.MediaAsset, error) {
	var m models.MediaAsset
	err := r.db.WithTenant(ctx, clinicID, func(tx pgx.Tx) error {
		return scanMediaAsset(tx.QueryRow(ctx, markMediaFailedSQL, id), &m)
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}
