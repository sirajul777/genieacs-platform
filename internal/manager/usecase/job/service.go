package job

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"

	domain "github.com/sirajul777/genieacs-platform/internal/manager/domain/job"
	"github.com/sirajul777/genieacs-platform/internal/manager/repository"
)

var (
	ErrInvalidJob      = errors.New("job type is required")
	ErrInvalidProgress = errors.New("progress must be between 0 and 100")
)

type Service struct {
	jobs   repository.JobRepository
	events repository.JobEventRepository
}

func NewService(jobs repository.JobRepository, events repository.JobEventRepository) *Service {
	return &Service{jobs: jobs, events: events}
}

type CreateInput struct {
	Type    string
	Payload []byte
}
type AssignNextInput struct{ AgentID string }
type ProgressInput struct {
	ID       string
	Progress int
}
type CompleteInput struct {
	ID     string
	Result []byte
}
type FailInput struct {
	ID    string
	Error string
}

func (s *Service) Create(ctx context.Context, input CreateInput) (domain.Job, error) {
	if input.Type == "" {
		return domain.Job{}, ErrInvalidJob
	}
	created, err := s.jobs.Create(ctx, domain.Job{ID: newID(), Type: input.Type, Payload: input.Payload, Status: domain.StatusPending})
	if err != nil {
		return domain.Job{}, err
	}
	_, _ = s.events.Create(ctx, created.ID, "created", "job created")
	return created, nil
}

func (s *Service) AssignNext(ctx context.Context, input AssignNextInput) (domain.Job, error) {
	item, err := s.jobs.AssignNext(ctx, input.AgentID)
	if err != nil {
		return domain.Job{}, err
	}
	_, _ = s.events.Create(ctx, item.ID, "assigned", input.AgentID)
	return item, nil
}

func (s *Service) Progress(ctx context.Context, input ProgressInput) (domain.Job, error) {
	if input.Progress < 0 || input.Progress > 100 {
		return domain.Job{}, ErrInvalidProgress
	}
	item, err := s.jobs.Progress(ctx, input.ID, input.Progress)
	if err != nil {
		return domain.Job{}, err
	}
	_, _ = s.events.Create(ctx, item.ID, "progress", "job progressed")
	return item, nil
}

func (s *Service) Complete(ctx context.Context, input CompleteInput) (domain.Job, error) {
	item, err := s.jobs.Complete(ctx, input.ID, input.Result)
	if err != nil {
		return domain.Job{}, err
	}
	_, _ = s.events.Create(ctx, item.ID, "completed", "job completed")
	return item, nil
}

func (s *Service) Fail(ctx context.Context, input FailInput) (domain.Job, error) {
	item, err := s.jobs.Fail(ctx, input.ID, input.Error)
	if err != nil {
		return domain.Job{}, err
	}
	_, _ = s.events.Create(ctx, item.ID, "failed", input.Error)
	return item, nil
}

func newID() string {
	buf := make([]byte, 16)
	_, _ = rand.Read(buf)
	buf[6] = (buf[6] & 0x0f) | 0x40
	buf[8] = (buf[8] & 0x3f) | 0x80
	return hex.EncodeToString(buf[0:4]) + "-" + hex.EncodeToString(buf[4:6]) + "-" + hex.EncodeToString(buf[6:8]) + "-" + hex.EncodeToString(buf[8:10]) + "-" + hex.EncodeToString(buf[10:16])
}
