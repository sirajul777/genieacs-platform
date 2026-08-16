package extension

import (
	"context"
	"testing"
)

type testExtension struct{}

func (testExtension) Name() string { return "genieacs" }
func (testExtension) Execute(_ context.Context, req Request) (Response, error) {
	return Response{Payload: req.Payload}, nil
}

func TestRegistryAndRunner(t *testing.T) {
	r := NewRegistry(testExtension{})
	runner := NewRunner(r)
	got, err := runner.Execute(context.Background(), "genieacs", Request{Payload: []byte(`{"ok":true}`)})
	if err != nil {
		t.Fatal(err)
	}
	if string(got.Payload) != `{"ok":true}` {
		t.Fatalf("unexpected payload: %s", got.Payload)
	}
}

func TestRunnerUnknownExtension(t *testing.T) {
	_, err := NewRunner(NewRegistry()).Execute(context.Background(), "missing", Request{})
	if err != ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
