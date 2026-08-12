package models

type AIResponse struct {
	Message string `json:"message"`
	Data    AIData `json:"data"`
}

type AIData struct {
	Meta                    Meta       `json:"meta"`
	Subjective              Subjective `json:"subjective"`
	SuggestedObjectiveFocus []string   `json:"suggested_objective_focus"`
}

type Meta struct {
	LanguageDetected      string `json:"language_detected"`
	SpeakerIdentification string `json:"speaker_identification"`
}

type Subjective struct {
	PresentingComplaint    string      `json:"presenting_complaint"`
	HistoryOfComplaint     string      `json:"history_of_complaint"`
	PainProfile            PainProfile `json:"pain_profile"`
	RedFlags               []string    `json:"red_flags"`
	AssociatedSymptoms     []string    `json:"associated_symptoms"`
	RelevantMedicalHistory []string    `json:"relevant_medical_history"`
	DrugHistory            []string    `json:"drug_history"`
	PastSurgicalHistory    []string    `json:"past_surgical_history"`
	FamilyHistory          string      `json:"family_history"`
	SocialHistory          []string    `json:"social_history"`
}

type PainProfile struct {
	// Intensity is nullable: the AI schema declares it as number().nullable()
	// and returns null when pain intensity was not reported, so a pointer is
	// required to unmarshal null without erroring.
	Intensity   *int     `json:"intensity"`
	Quality     string   `json:"quality"`
	Aggravating []string `json:"aggravating"`
	Alleviating []string `json:"alleviating"`
	Duration    string   `json:"duration"`
	Location    []string `json:"location"`
}
