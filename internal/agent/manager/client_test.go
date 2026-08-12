package manager

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientPollAndLifecycle(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer token" { t.Fatalf("missing auth") }
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/api/v1/jobs/poll" { _, _ = w.Write([]byte(`{"id":"job-1","type":"test.example","payload":{"x":1},"status":"in_progress"}`)); return }
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	c := NewClient(server.URL, "token", server.Client())
	job, err := c.Poll(context.Background(), "agent-1")
	if err != nil || job.ID != "job-1" || job.Type != "test.example" { t.Fatalf("unexpected job: %#v err=%v", job, err) }
	if err := c.Progress(context.Background(), job.ID, 50, "half"); err != nil { t.Fatal(err) }
	if err := c.Complete(context.Background(), job.ID, map[string]any{"ok": true}); err != nil { t.Fatal(err) }
}

func TestClientPollNoJob(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) }))
	defer server.Close()
	_, err := NewClient(server.URL, "", server.Client()).Poll(context.Background(), "agent-1")
	if err != ErrNoJob { t.Fatalf("expected ErrNoJob, got %v", err) }
}
