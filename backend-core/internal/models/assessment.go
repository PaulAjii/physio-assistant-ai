package models

import "time"

type ObjectiveFinding struct {
	Category string `json:"category"`
	Result   string `json:"result"`
	Notes    string `json:"notes"`
	Test string `json:"test"`
}

type AssessmentSubmission struct {
	JobID string `json:"jobID"`
	Complaint string `json:"complaint"`
	ObjectiveFindings []ObjectiveFinding `json:"objective_findings"`
	SubjectiveFindings Subjective `json:"subjective_findings"`
}

type CollatedAssessment struct {
	ID string `json:"id"`
	Complaint string `json:"complaint"`
	ObjectiveFindings []ObjectiveFinding `json:"objective_findings"`
	SubjectiveFindings Subjective `json:"subjective_findings"`
	AIDraft AIData `json:"ai_draft"`
	CreatedAt time.Time `json:"created_at"`
}