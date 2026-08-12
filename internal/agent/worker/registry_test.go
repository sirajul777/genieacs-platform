package worker

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

type testWorker struct {
	typeName string
	fn       func(context.Context) (Result, error)
}

func (w testWorker) Type() string { return w.typeName }

func (w testWorker) Execute(ctx context.Context, _ Job, _ ProgressReporter) (Result, error) {
	return w.fn(ctx)
}

type testProgress struct{}

func (testProgress) Report(context.Context, int, string) error { return nil }

func TestRegistryRejectsDuplicateTypes(t *testing.T) {
	r := NewRegistry()
	w := testWorker{typeName: "test.example", fn: func(context.Context) (Result, error) { return Result{}, nil }}

	if err := r.Register(w); err != nil {
		t.Fatalf("first registration: %v", err)
	}
	if err := r.Register(w); !errors.Is(err, ErrDuplicateType) {
		t.Fatalf("expected duplicate error, got %v", err)
	}
}

func TestRegistryTypesAreDeterministic(t *testing.T) {
	r := NewRegistry()
	for _, typeName := range []string{"z.example", "a.example", "m.example"} {
		if err := r.Register(testWorker{typeName: typeName}); err != nil {
			t.Fatal(err)
		}
	}

	got := r.Types()
	want := []string{"a.example", "m.example", "z.example"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("types = %v, want %v", got, want)
	}
}

func TestRuntimeRejectsUnknownWorker(t *testing.T) {
	r := NewRegistry()
	runtime := NewRuntime(r)

	_, err := runtime.Execute(context.Background(), Job{ID: "job-1", Type: "missing.example"}, testProgress{})
	if !errors.Is(err, ErrWorkerNotFound) {
		t.Fatalf("expected worker not found, got %v", err)
	}
}

func TestRuntimePropagatesCancellation(t *testing.T) {
	r := NewRegistry()
	cancelled := make(chan struct{})
	if err := r.Register(testWorker{
		typeName: "cancel.example",
		fn: func(ctx context.Context) (Result, error) {
			<-ctx.Done()
			close(cancelled)
			return Result{}, ctx.Err()
		},
	}); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := NewRuntime(r).Execute(ctx, Job{ID: "job-1", Type: "cancel.example"}, testProgress{})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context cancellation, got %v", err)
	}
	select {
	case <-cancelled:
		t.Fatalf("worker should not start after cancellation")
	default:
	}
}

func TestRuntimeRejectsNilProgressReporter(t *testing.T) {
	r := NewRegistry()
	if err := r.Register(testWorker{typeName: "test.example"}); err != nil {
		t.Fatal(err)
	}

	_, err := NewRuntime(r).Execute(context.Background(), Job{Type: "test.example"}, nil)
	if !errors.Is(err, ErrNilProgressReporter) {
		t.Fatalf("expected nil progress reporter error, got %v", err)
	}
}
