package models

import "time"

// Clinic is a tenant: the root every other tenant-scoped row hangs off. Clinics
// are created by the owner-privileged seed step, not by the application, so
// there is no app-level create path (see repo.ClinicRepo).
type Clinic struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
