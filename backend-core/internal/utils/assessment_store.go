package utils

import (
	"sync"
	"github.com/PaulAjii/physio-assistant-ai/internal/models"
)

type AssessmentStore struct {
	mu sync.RWMutex
	assessments map[string]*models.CollatedAssessment
	drafts map[string]*models.AIData
}

var Assessments = &AssessmentStore{
	assessments: make(map[string]*models.CollatedAssessment),
	drafts: make(map[string]*models.AIData),
}

func (s *AssessmentStore) SaveDraft(jobID string, draft *models.AIData) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.drafts[jobID] = draft
}

func (s *AssessmentStore) GetDraft(jobID string) *models.AIData {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.drafts[jobID]
}

func (s *AssessmentStore) SaveAssessment(assessment *models.CollatedAssessment) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.assessments[assessment.ID] = assessment
}

func (s *AssessmentStore) GetAssessment(id string) *models.CollatedAssessment {
	s.mu.RLock()
	defer s.mu.RLock()
	return s.assessments[id]
}

func (s *AssessmentStore) GetAllAssessments() []*models.CollatedAssessment {
	s.mu.RLock()
	defer s.mu.RUnlock()
	assessments := make([]*models.CollatedAssessment, 0, len(s.assessments))
	 for _, assessment := range s.assessments {
		assessments = append(assessments, assessment)
	 }
	 return assessments
}