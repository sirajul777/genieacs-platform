package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/sirajul777/genieacs-platform/internal/agent/extension"
	"github.com/sirajul777/genieacs-platform/internal/agent/genieacs"
)

type GenieACSWorker struct {
	NameValue string
	Runner    *extension.Runner
}

func NewGenieACSWorker(runner *extension.Runner, jobType string) *GenieACSWorker {
	return &GenieACSWorker{NameValue: jobType, Runner: runner}
}

func (w *GenieACSWorker) Type() string { return w.NameValue }

func (w *GenieACSWorker) Execute(ctx context.Context, job Job, progress ProgressReporter) (Result, error) {
	if w.Runner == nil {
		return Result{}, errors.New("extension runner is required")
	}
	if err := progress.Report(ctx, 10, "job started"); err != nil {
		return Result{}, err
	}

	var payload struct {
		Args map[string]any `json:"args"`
	}
	if len(job.Payload) > 0 {
		if err := json.Unmarshal(job.Payload, &payload); err != nil {
			return Result{}, fmt.Errorf("decode job payload: %w", err)
		}
	}

	action := ""
	switch w.NameValue {
	case "genieacs.install":
		action = "install"
	case "genieacs.verify":
		action = "verify"
	default:
		return Result{}, fmt.Errorf("unsupported GenieACS worker: %s", w.NameValue)
	}

	if err := progress.Report(ctx, 25, action+" requested"); err != nil {
		return Result{}, err
	}
	out, err := w.Runner.Execute(ctx, genieacs.Name, extension.Request{Action: action, Payload: mustJSON(struct {
		Action string         `json:"action"`
		Args   map[string]any `json:"args"`
	}{Action: action, Args: payload.Args})})
	if err != nil {
		return Result{Payload: out.Payload}, err
	}
	if err := progress.Report(ctx, 100, action+" completed"); err != nil {
		return Result{}, err
	}
	return Result{Payload: out.Payload}, nil
}

func mustJSON(v any) []byte {
	b, _ := json.Marshal(v)
	return b
}
