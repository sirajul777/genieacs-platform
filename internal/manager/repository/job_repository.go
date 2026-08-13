package repository

import (
	"context"
	"errors"

	"github.com/sirajul777/genieacs-platform/internal/manager/domain/job"
)

var ErrNoPendingJob = errors.New("no pending job")

type JobRepository interface {
	Create(ctx context.Context, item job.Job) (job.Job, error)
	AssignNext(ctx context.Context, agentID string) (job.Job, error)
	Progress(ctx context.Context, id string, progress int) (job.Job, error)
	Complete(ctx context.Context, id string, result []byte) (job.Job, error)
	Fail(ctx context.Context, id string, message string) (job.Job, error)
}

type JobEventRepository interface {
	Create(ctx context.Context, jobID, eventType, message string) (job.Event, error)
}
