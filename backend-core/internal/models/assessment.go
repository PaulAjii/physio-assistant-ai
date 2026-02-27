package models

import "time"

type ObjectiveFinding struct {
	Category string `json:"category"`
	Result   string `json:"result"`
	Notes    string `json:"notes"`
	Test     string `json:"test"`
}

type AssessmentSubmission struct {
	JobID               string             `json:"jobID"`
	Complaint           string             `json:"complaint"`
	ObjectiveFindings   []ObjectiveFinding `json:"objective_findings"`
	CorrectedSubjective Subjective         `json:"corrected_subjective"`
}

type Assessment struct {
	PresentingComplaint    string             `json:"presenting_complaint"`
	HistoryOfComplaint     string             `json:"history_of_complaint"`
	PainProfile            PainProfile        `json:"pain_profile"`
	RedFlags               []string           `json:"red_flags"`
	AssociatedSymptoms     []string           `json:"associated_symptoms"`
	RelevantMedicalHistory []string           `json:"relevant_medical_history"`
	DrugHistory            []string           `json:"drug_history"`
	PastSurgicalHistory    []string           `json:"past_surgical_history"`
	FamilyHistory          string             `json:"family_history"`
	SocialHistory          []string           `json:"social_history"`
	ObjectiveFindings      []ObjectiveFinding `json:"objective_findings"`
}

type CollatedAssessment struct {
	ID         string     `json:"id"`
	Complaint  string     `json:"complaint"`
	Assessment Assessment `json:"assessment"`
	AIDraft    AIResponse `json:"ai_draft"`
	CreatedAt  time.Time  `json:"created_at"`
}
