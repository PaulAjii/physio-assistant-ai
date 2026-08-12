package models

import (
	"encoding/json"
	"time"
)

// ObjectiveTemplate is a clinic-configurable scaffold for the objective exam.
// Body is an ordered JSON array of field definitions (key/label/type/unit/...),
// kept as raw JSON because the shape is owned by the clinic. The clinician speaks
// findings; matching keys get filled and the rest are entered by hand. IsActive
// lets a clinic retire a template without deleting it (and losing history).
type ObjectiveTemplate struct {
	ID        string          `json:"id"`
	ClinicID  string          `json:"clinic_id"`
	Name      string          `json:"name"`
	Body      json.RawMessage `json:"body"`
	IsActive  bool            `json:"is_active"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}
