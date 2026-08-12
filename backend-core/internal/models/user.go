package models

import "time"

// User is a clinic member: an admin or a clinician. It is the domain entity
// shared across layers (repo, services, handlers). The API should not serialise
// it directly — a dto is used for responses — but PasswordHash carries json:"-"
// as a defence-in-depth backstop so it can never leak even if one does.
type User struct {
	ID           string    `json:"id"`
	ClinicID     string    `json:"clinic_id"`
	Email        string    `json:"email"`
	FullName     string    `json:"full_name"`
	Role         string    `json:"role"`
	LicenseNo    *string   `json:"license_no,omitempty"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
