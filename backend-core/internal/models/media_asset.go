package models

import "time"

// Media kinds and upload statuses. An asset is created pending (the storage key
// is reserved), flips to uploaded once the client confirms the blob landed, or
// failed if the upload did not complete.
const (
	MediaKindAudio = "audio"
	MediaKindVideo = "video"

	MediaStatusPending  = "pending"
	MediaStatusUploaded = "uploaded"
	MediaStatusFailed   = "failed"
)

// MediaAsset is an audio or video recording for a session. The blob itself lives
// in object storage (MinIO/R2); Postgres only holds StorageKey, which is
// presigned on demand — nothing binary is stored in the database. ContentType and
// Bytes are nil until the upload is confirmed.
type MediaAsset struct {
	ID          string    `json:"id"`
	ClinicID    string    `json:"clinic_id"`
	SessionID   string    `json:"session_id"`
	Kind        string    `json:"kind"`
	StorageKey  string    `json:"storage_key"`
	ContentType *string   `json:"content_type,omitempty"`
	Bytes       *int64    `json:"bytes,omitempty"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
