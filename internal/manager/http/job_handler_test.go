package http

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	domain "github.com/sirajul777/genieacs-platform/internal/manager/domain/job"
	jobusecase "github.com/sirajul777/genieacs-platform/internal/manager/usecase/job"
)

type handlerJobs struct{ items map[string]domain.Job }

func newHandlerJobs() *handlerJobs { return &handlerJobs{items: map[string]domain.Job{}} }
func (f *handlerJobs) Create(ctx context.Context, item domain.Job) (domain.Job, error) {
	f.items[item.ID] = item
	return item, nil
}
func (f *handlerJobs) AssignNext(ctx context.Context, agentID string) (domain.Job, error) {
	for id, item := range f.items {
		if item.Status == domain.StatusPending {
			item.Status = domain.StatusInProgress
			item.AgentID = &agentID
			f.items[id] = item
			return item, nil
		}
	}
	return domain.Job{}, errors.New("none")
}
func (f *handlerJobs) Progress(ctx context.Context, id string, progress int) (domain.Job, error) {
	item := f.items[id]
	item.Progress = progress
	f.items[id] = item
	return item, nil
}
func (f *handlerJobs) Complete(ctx context.Context, id string, result []byte) (domain.Job, error) {
	item := f.items[id]
	item.Status = domain.StatusCompleted
	item.Progress = 100
	f.items[id] = item
	return item, nil
}
func (f *handlerJobs) Fail(ctx context.Context, id string, message string) (domain.Job, error) {
	item := f.items[id]
	item.Status = domain.StatusFailed
	item.Error = message
	f.items[id] = item
	return item, nil
}

type handlerEvents struct{}

func (handlerEvents) Create(ctx context.Context, jobID, eventType, message string) (domain.Event, error) {
	return domain.Event{JobID: jobID, Type: eventType, Message: message}, nil
}

func TestJobHandlerLifecycle(t *testing.T) {
	jobs := newHandlerJobs()
	h := NewJobHandler(jobusecase.NewService(jobs, handlerEvents{}))
	rec := httptest.NewRecorder()
	h.Create(rec, httptest.NewRequest(http.MethodPost, "/api/v1/jobs", bytes.NewBufferString(`{"type":"generic","payload":{"a":1}}`)))
	if rec.Code != http.StatusCreated {
		t.Fatalf("create status %d: %s", rec.Code, rec.Body.String())
	}
	var id string
	for k := range jobs.items {
		id = k
	}
	rec = httptest.NewRecorder()
	h.Poll(rec, httptest.NewRequest(http.MethodPost, "/api/v1/jobs/poll", bytes.NewBufferString(`{"agent_id":"agent-1"}`)))
	if rec.Code != http.StatusOK {
		t.Fatalf("poll status %d", rec.Code)
	}
	rec = httptest.NewRecorder()
	h.Progress(rec, httptest.NewRequest(http.MethodPost, "/api/v1/jobs/"+id+"/progress", bytes.NewBufferString(`{"progress":25}`)))
	if rec.Code != http.StatusOK {
		t.Fatalf("progress status %d", rec.Code)
	}
	rec = httptest.NewRecorder()
	h.Complete(rec, httptest.NewRequest(http.MethodPost, "/api/v1/jobs/"+id+"/complete", bytes.NewBufferString(`{"result":{"ok":true}}`)))
	if rec.Code != http.StatusOK {
		t.Fatalf("complete status %d", rec.Code)
	}
}
func TestJobHandlerInvalidProgress(t *testing.T) {
	h := NewJobHandler(jobusecase.NewService(newHandlerJobs(), handlerEvents{}))
	rec := httptest.NewRecorder()
	h.Progress(rec, httptest.NewRequest(http.MethodPost, "/api/v1/jobs/job-1/progress", bytes.NewBufferString(`{"progress":101}`)))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", rec.Code)
	}
}
func TestJobHandlerFail(t *testing.T) {
	jobs := newHandlerJobs()
	h := NewJobHandler(jobusecase.NewService(jobs, handlerEvents{}))
	created, _ := h.service.Create(context.Background(), jobusecase.CreateInput{Type: "generic", Payload: []byte(`{}`)})
	_, _ = h.service.AssignNext(context.Background(), jobusecase.AssignNextInput{AgentID: "agent-1"})
	rec := httptest.NewRecorder()
	h.Fail(rec, httptest.NewRequest(http.MethodPost, "/api/v1/jobs/"+created.ID+"/fail", bytes.NewBufferString(`{"error":"boom"}`)))
	if rec.Code != http.StatusOK {
		t.Fatalf("fail status %d", rec.Code)
	}
}
