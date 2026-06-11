package jobs

import (
	"fmt"
	"sync"
	"time"

	"github.com/ogilcher/lunar-deploy-agent/internal/deploy/engine"
)

type Store struct {
	jobs  map[string]*DeploymentJob
	mutex sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		jobs: map[string]*DeploymentJob{},
	}
}

func (s *Store) CreateJob(
	deploymentName string,
) *DeploymentJob {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	job := &DeploymentJob{
		ID:             fmt.Sprintf("%d", time.Now().UnixNano()),
		DeploymentName: deploymentName,
		Status:         JobQueued,
		CreatedAt:      time.Now(),
	}

	s.jobs[job.ID] = job

	return job
}

func (s *Store) GetJob(
	id string,
) (*DeploymentJob, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	job, exists := s.jobs[id]
	return job, exists
}

func (s *Store) ListJobs() []*DeploymentJob {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	results := []*DeploymentJob{}

	for _, job := range s.jobs {
		results = append(results, job)
	}

	return results
}

func (s *Store) MarkRunning(id string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if job, exists := s.jobs[id]; exists {
		job.Status = JobRunning
		job.StartedAt = new(time.Now())
	}
}

func (s *Store) MarkSucceeded(
	id string,
	result *engine.DeploymentResult,
) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if job, exists := s.jobs[id]; exists {
		job.Status = JobSucceeded
		job.FinishedAt = new(time.Now())
		job.Result = result
	}
}

func (s *Store) MarkFailed(
	id string,
	result *engine.DeploymentResult,
	err error,
) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if job, exists := s.jobs[id]; exists {
		job.Status = JobFailed
		job.FinishedAt = new(time.Now())
		job.Result = result

		if err != nil {
			job.Error = err.Error()
		}
	}
}

func (s *Store) MarkCancelled(id string) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	job, exists := s.jobs[id]
	if !exists {
		return false
	}

	if job.Status == JobSucceeded || job.Status == JobFailed {
		return false
	}

	job.Status = JobCancelled
	job.FinishedAt = new(time.Now())

	return true
}

func (s *Store) IsCancelled(id string) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	job, exists := s.jobs[id]

	if !exists {
		return false
	}

	return job.Status == JobCancelled
}
