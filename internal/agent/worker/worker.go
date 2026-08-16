package worker

import "context"

// Job is the execution contract delivered to a worker by the Agent runtime.
// Payload is opaque to the runtime and owned by the worker implementation.
type Job struct {
	ID      string
	Type    string
	Payload []byte
}

// Result contains the worker output. The runtime does not interpret it.
type Result struct {
	Payload []byte
}

// ProgressReporter lets a worker publish bounded execution progress without
// coupling the worker to HTTP, persistence, or the Manager implementation.
type ProgressReporter interface {
	Report(ctx context.Context, progress int, message string) error
}

// Worker executes one job of a single, stable type.
type Worker interface {
	Type() string
	Execute(ctx context.Context, job Job, progress ProgressReporter) (Result, error)
}
