package models

import "time"

// RefreshToken is the server-side record of an issued refresh token. Only a hash
// of the token is stored (TokenHash), never the token itself, so the database
// cannot be used to mint valid sessions. ClinicID is denormalised from the user
// so RLS can isolate these rows without a join. A token is valid while
// RevokedAt is nil and ExpiresAt is in the future.
type RefreshToken struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	ClinicID  string     `json:"clinic_id"`
	TokenHash string     `json:"-"`
	ExpiresAt time.Time  `json:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}
