package utils

import (
	"time"

	"github.com/PaulAjii/physio-assistant-ai/internal/models"
	"github.com/google/uuid"
)

func CollateAssessment(submission *models.AssessmentSubmission, draft models.AIResponse) *models.CollatedAssessment {
	return &models.CollatedAssessment{
		ID:        uuid.New().String(),
		Complaint: submission.Complaint,
		Assessment: models.Assessment{
			PresentingComplaint:    submission.Complaint,
			HistoryOfComplaint:     submission.CorrectedSubjective.HistoryOfComplaint,
			PainProfile:            submission.CorrectedSubjective.PainProfile,
			RedFlags:               submission.CorrectedSubjective.RedFlags,
			AssociatedSymptoms:     submission.CorrectedSubjective.AssociatedSymptoms,
			RelevantMedicalHistory: submission.CorrectedSubjective.RelevantMedicalHistory,
			DrugHistory:            submission.CorrectedSubjective.DrugHistory,
			PastSurgicalHistory:    submission.CorrectedSubjective.PastSurgicalHistory,
			FamilyHistory:          submission.CorrectedSubjective.FamilyHistory,
			SocialHistory:          submission.CorrectedSubjective.SocialHistory,
			ObjectiveFindings:      submission.ObjectiveFindings,
		},
		AIDraft:   draft,
		CreatedAt: time.Now(),
	}
}
