package utils

import (
	"time"

	"github.com/google/uuid"
	"github.com/PaulAjii/physio-assistant-ai/internal/models"
)

func CollateAssessment(submission *models.AssessmentSubmission, draft models.AIData) *models.CollatedAssessment {
	return &models.CollatedAssessment{
		ID: uuid.New().String(),
		Complaint: submission.Complaint,
		ObjectiveFindings: submission.ObjectiveFindings,
		SubjectiveFindings: submission.SubjectiveFindings,
		AIDraft: draft,
		CreatedAt: time.Now(),
	}
}