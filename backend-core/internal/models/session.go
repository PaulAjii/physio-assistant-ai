package models

import "time"

// Session statuses. A session starts as draft, moves to recording while audio
// is being captured, and ends at finalized — the terminal, signed state.
const (
	SessionStatusDraft     = "draft"
	SessionStatusRecording = "recording"
	SessionStatusFinalized = "finalized"
)

// Session is one physiotherapy visit and the spine the rest of the clinical data
// hangs off (media, transcript, assessment). ClinicianID is the user who owns the
// visit. The signing fields are nil until finalize, at which point the database
// CHECK guarantees they are all set together — so a finalized session is always
// attributable to a real clinician (SignedBy) at a known time (SignedAt).
type Session struct {
	ID          string     `json:"id"`
	ClinicID    string     `json:"clinic_id"`
	ClinicianID string     `json:"clinician_id"`
	PatientRef  *string    `json:"patient_ref,omitempty"`
	Status      string     `json:"status"`
	Language    *string    `json:"language,omitempty"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	FinalizedAt *time.Time `json:"finalized_at,omitempty"`
	SignedBy    *string    `json:"signed_by,omitempty"`
	SignedAt    *time.Time `json:"signed_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
