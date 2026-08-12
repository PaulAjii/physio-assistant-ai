package repo

import (
	"context"

	"github.com/PaulAjii/physio-assistant-ai/internal/db"
	"github.com/PaulAjii/physio-assistant-ai/internal/models"
	"github.com/jackc/pgx/v5"
)

// CreateSegmentParams is one chunk of speech-to-text output. Speaker, StartMs and
// EndMs are optional (the STT provider may not supply them). IsFinal is false for
// interim results and true once the provider commits the text.
type CreateSegmentParams struct {
	SessionID string
	Speaker   *string
	Content   string
	StartMs   *int
	EndMs     *int
	IsFinal   bool
}

// TranscriptSegments is the persistence contract for the append-only transcript
// log. Segments are never updated or deleted; interim ones are simply superseded
// by later final ones, so there is no Update or soft-delete here.
type TranscriptSegments interface {
	Create(ctx context.Context, clinicID string, p CreateSegmentParams) (*models.TranscriptSegment, error)
	ListBySession(ctx context.Context, clinicID, sessionID string) ([]models.TranscriptSegment, error)
}

// TranscriptSegmentRepo is the Postgres-backed implementation of TranscriptSegments.
type TranscriptSegmentRepo struct {
	db *db.DB
}

// NewTranscriptSegmentRepo returns a TranscriptSegmentRepo backed by the given database.
func NewTranscriptSegmentRepo(database *db.DB) *TranscriptSegmentRepo {
	return &TranscriptSegmentRepo{db: database}
}

// Compile-time check that TranscriptSegmentRepo satisfies TranscriptSegments.
var _ TranscriptSegments = (*TranscriptSegmentRepo)(nil)

const segmentColumns = `id, clinic_id, session_id, speaker, content, start_ms, end_ms, is_final, created_at`

func scanSegment(row pgx.Row, s *models.TranscriptSegment) error {
	return row.Scan(
		&s.ID, &s.ClinicID, &s.SessionID, &s.Speaker, &s.Content,
		&s.StartMs, &s.EndMs, &s.IsFinal, &s.CreatedAt,
	)
}

const createSegmentSQL = `
INSERT INTO transcript_segments (clinic_id, session_id, speaker, content, start_ms, end_ms, is_final)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING ` + segmentColumns

// Create appends a segment to a session's transcript.
func (r *TranscriptSegmentRepo) Create(ctx context.Context, clinicID string, p CreateSegmentParams) (*models.TranscriptSegment, error) {
	var s models.TranscriptSegment
	err := r.db.WithTenant(ctx, clinicID, func(tx pgx.Tx) error {
		return scanSegment(tx.QueryRow(ctx, createSegmentSQL,
			clinicID, p.SessionID, p.Speaker, p.Content, p.StartMs, p.EndMs, p.IsFinal,
		), &s)
	})
	if err != nil {
		return nil, err
	}
	return &s, nil
}

const listSegmentsBySessionSQL = `
SELECT ` + segmentColumns + `
FROM transcript_segments
WHERE session_id = $1
ORDER BY created_at`

// ListBySession returns a session's segments in the order they were appended.
// Callers building a clean transcript filter on IsFinal; interim rows are kept so
// a live view can show text as it streams in.
func (r *TranscriptSegmentRepo) ListBySession(ctx context.Context, clinicID, sessionID string) ([]models.TranscriptSegment, error) {
	var list []models.TranscriptSegment
	err := r.db.WithTenant(ctx, clinicID, func(tx pgx.Tx) error {
		rows, err := tx.Query(ctx, listSegmentsBySessionSQL, sessionID)
		if err != nil {
			return err
		}
		list, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (models.TranscriptSegment, error) {
			var s models.TranscriptSegment
			err := scanSegment(row, &s)
			return s, err
		})
		return err
	})
	if err != nil {
		return nil, err
	}
	return list, nil
}
