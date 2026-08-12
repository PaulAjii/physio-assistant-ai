package models

import "time"

// TranscriptSegment is one chunk of live speech-to-text output. Segments are an
// append-only log for a session: interim results are stored with IsFinal=false
// and later superseded by final ones, so there is no soft-delete or update — the
// record is immutable once written. Speaker, StartMs and EndMs are nil when the
// STT provider does not supply them.
type TranscriptSegment struct {
	ID        string    `json:"id"`
	ClinicID  string    `json:"clinic_id"`
	SessionID string    `json:"session_id"`
	Speaker   *string   `json:"speaker,omitempty"`
	Content   string    `json:"content"`
	StartMs   *int      `json:"start_ms,omitempty"`
	EndMs     *int      `json:"end_ms,omitempty"`
	IsFinal   bool      `json:"is_final"`
	CreatedAt time.Time `json:"created_at"`
}
