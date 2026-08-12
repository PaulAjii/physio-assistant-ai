package models

import "time"

// Invitation is an admin's invite for someone to join their clinic. The invitee
// accepts with the raw token (delivered out-of-band for now); only its hash is
// stored, same reasoning as RefreshToken. An invitation is pending while
// AcceptedAt is nil and ExpiresAt is in the future. InvitedBy is the admin who
// issued it and is nullable so an invite outlives a deleted inviter.
type Invitation struct {
	ID         string     `json:"id"`
	ClinicID   string     `json:"clinic_id"`
	Email      string     `json:"email"`
	Role       string     `json:"role"`
	TokenHash  string     `json:"-"`
	InvitedBy  *string    `json:"invited_by,omitempty"`
	ExpiresAt  time.Time  `json:"expires_at"`
	AcceptedAt *time.Time `json:"accepted_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}
