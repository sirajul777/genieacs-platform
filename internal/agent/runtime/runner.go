package runtime

import (
	"context"
	"errors"
	"time"

	"github.com/sirajul777/genieacs-platform/internal/agent/manager"
	"github.com/sirajul777/genieacs-platform/internal/agent/worker"
)

var ErrStopped = errors.New("agent runner stopped")

type Runner struct {
	manager  *manager.Client
	runtime  *worker.Runtime
	agentID  string
	interval time.Duration
}

func NewRunner(client *manager.Client, runtime *worker.Runtime, agentID string, interval time.Duration) *Runner {
	if interval <= 0 {
		interval = 2 * time.Second
	}
	return &Runner{manager: client, runtime: runtime, agentID: agentID, interval: interval}
}

func (r *Runner) Run(ctx context.Context) error {
	if r.agentID == "" {
		return errors.New("agent id is required")
	}

	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	for {
		if err := r.runOnce(ctx); err != nil && !errors.Is(err, manager.ErrNoJob) {
			if ctx.Err() != nil {
				return ctx.Err()
			}
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}

func (r *Runner) runOnce(ctx context.Context) error {
	job, err := r.manager.Poll(ctx, r.agentID)
	if err != nil {
		return err
	}

	progress := &reporter{client: r.manager, jobID: job.ID}
	result, err := r.runtime.Execute(ctx, worker.Job{ID: job.ID, Type: job.Type, Payload: job.Payload}, progress)
	if err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return r.manager.Fail(ctx, job.ID, err.Error())
	}
	return r.manager.Complete(ctx, job.ID, result.Payload)
}

type reporter struct {
	client *manager.Client
	jobID  string
}

func (r *reporter) Report(ctx context.Context, progress int, message string) error {
	if progress < 0 {
		progress = 0
	}
	if progress > 100 {
		progress = 100
	}
	return r.client.Progress(ctx, r.jobID, progress, message)
}
