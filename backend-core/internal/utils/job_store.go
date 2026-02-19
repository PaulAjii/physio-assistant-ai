package utils

import "sync"

type JobStatus struct {
	Status bool
	Result []byte
	Error  error
}

type JobStore struct {
	mu      sync.Mutex
	jobs    map[string]*JobStatus
	signals map[string]chan struct{}
}

var Store = &JobStore{
	jobs:    make(map[string]*JobStatus),
	signals: make(map[string]chan struct{}),
}

func (s *JobStore) Register(jobId string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.jobs[jobId] = &JobStatus{}
	s.signals[jobId] = make(chan struct{})
}

func (s *JobStore) Complete(jobId string, result []byte, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if job, ok := s.jobs[jobId]; ok {
		job.Status = true
		job.Result = result
		job.Error = err
		close(s.signals[jobId])
	}
}

func (s *JobStore) Wait(jobID string) <-chan struct{} {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.signals[jobID]
}

func (s *JobStore) Get(jobID string) *JobStatus {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.jobs[jobID]
}

func (s *JobStore) Delete(jobID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.jobs, jobID)
	delete(s.signals, jobID)
}
