package worker

import (
	"context"
	"errors"
)

var ErrNilProgressReporter = errors.New("progress reporter is required")

// Runtime resolves a job type and executes it with cancellation propagated
// directly from the Agent lifecycle.
type Runtime struct {
	registry *Registry
}

func NewRuntime(registry *Registry) *Runtime {
	return &Runtime{registry: registry}
}

func (r *Runtime) Execute(ctx context.Context, job Job, progress ProgressReporter) (Result, error) {
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	if progress == nil {
		return Result{}, ErrNilProgressReporter
	}

	w, err := r.registry.Get(job.Type)
	if err != nil {
		return Result{}, err
	}

	result, err := w.Execute(ctx, job, progress)
	if err != nil {
		return Result{}, err
	}
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	return result, nil
}
