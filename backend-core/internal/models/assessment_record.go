package models

import (
	"encoding/json"
	"time"
)

// AssessmentRecord is the persisted row for the assessments table — the durable
// replacement for the in-memory CollatedAssessment store. It wraps the clinical
// content as two jsonb sections plus persistence metadata; the content shape
// itself is modelled by Assessment and its parts. The sections are kept as raw
// JSON here because the shape is clinic-defined (via objective templates) and
// evolving, so it can persist without a schema change.
//
// There is one live record per session (partial unique index), and Version bumps
// on every edit so pre-finalize changes are visible. PainIntensity is a pointer
// so "not reported" (nil) is distinct from a reported 0; the database bounds it
// to the 0-10 scale.
type AssessmentRecord struct {
	ID            string          `json:"id"`
	ClinicID      string          `json:"clinic_id"`
	SessionID     string          `json:"session_id"`
	Subjective    json.RawMessage `json:"subjective"`
	Objective     json.RawMessage `json:"objective"`
	AISummary     *string         `json:"ai_summary,omitempty"`
	PainIntensity *int            `json:"pain_intensity,omitempty"`
	Version       int             `json:"version"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}
