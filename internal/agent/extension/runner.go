package extension

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("extension not found")

type Runner struct {
	registry *Registry
}

func NewRunner(registry *Registry) *Runner { return &Runner{registry: registry} }

func (r *Runner) Execute(ctx context.Context, name string, req Request) (Response, error) {
	if err := ctx.Err(); err != nil {
		return Response{}, err
	}
	ext, ok := r.registry.Get(name)
	if !ok {
		return Response{}, ErrNotFound
	}
	result, err := ext.Execute(ctx, req)
	if err != nil {
		return Response{}, err
	}
	if err := ctx.Err(); err != nil {
		return Response{}, err
	}
	return result, nil
}
