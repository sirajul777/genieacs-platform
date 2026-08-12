package job

import (
	"context"
	"errors"
	"testing"

	domain "github.com/sirajul777/genieacs-platform/internal/manager/domain/job"
)

type fakeJobs struct{ items map[string]domain.Job }

func newFakeJobs() *fakeJobs { return &fakeJobs{items: map[string]domain.Job{}} }
func (f *fakeJobs) Create(ctx context.Context, item domain.Job) (domain.Job, error) {
	f.items[item.ID] = item
	return item, nil
}
func (f *fakeJobs) AssignNext(ctx context.Context, agentID string) (domain.Job, error) {
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
func (f *fakeJobs) Progress(ctx context.Context, id string, progress int) (domain.Job, error) {
	item := f.items[id]
	item.Progress = progress
	f.items[id] = item
	return item, nil
}
func (f *fakeJobs) Complete(ctx context.Context, id string, result []byte) (domain.Job, error) {
	item := f.items[id]
	item.Status = domain.StatusCompleted
	item.Progress = 100
	item.Result = result
	f.items[id] = item
	return item, nil
}
func (f *fakeJobs) Fail(ctx context.Context, id string, message string) (domain.Job, error) {
	item := f.items[id]
	item.Status = domain.StatusFailed
	item.Error = message
	f.items[id] = item
	return item, nil
}

type fakeEvents struct{ count int }

func (f *fakeEvents) Create(ctx context.Context, jobID, eventType, message string) (domain.Event, error) {
	f.count++
	return domain.Event{ID: int64(f.count), JobID: jobID, Type: eventType, Message: message}, nil
}

func TestCreateJob(t *testing.T) {
	jobs := newFakeJobs()
	events := &fakeEvents{}
	svc := NewService(jobs, events)
	item, err := svc.Create(context.Background(), CreateInput{Type: "generic", Payload: []byte(`{"x":1}`)})
	if err != nil {
		t.Fatal(err)
	}
	if item.Type != "generic" || item.Status != domain.StatusPending {
		t.Fatalf("unexpected job: %#v", item)
	}
	if events.count != 1 {
		t.Fatal("expected event")
	}
}
func TestAssignProgressComplete(t *testing.T) {
	jobs := newFakeJobs()
	svc := NewService(jobs, &fakeEvents{})
	created, _ := svc.Create(context.Background(), CreateInput{Type: "generic", Payload: []byte(`{}`)})
	assigned, err := svc.AssignNext(context.Background(), AssignNextInput{AgentID: "agent-1"})
	if err != nil {
		t.Fatal(err)
	}
	if assigned.ID != created.ID || assigned.Status != domain.StatusInProgress {
		t.Fatalf("bad assignment")
	}
	progressed, err := svc.Progress(context.Background(), ProgressInput{ID: created.ID, Progress: 50})
	if err != nil {
		t.Fatal(err)
	}
	if progressed.Progress != 50 {
		t.Fatal("bad progress")
	}
	completed, err := svc.Complete(context.Background(), CompleteInput{ID: created.ID, Result: []byte(`{"ok":true}`)})
	if err != nil {
		t.Fatal(err)
	}
	if completed.Status != domain.StatusCompleted || completed.Progress != 100 {
		t.Fatal("bad complete")
	}
}
func TestFailAndInvalidProgress(t *testing.T) {
	jobs := newFakeJobs()
	svc := NewService(jobs, &fakeEvents{})
	created, _ := svc.Create(context.Background(), CreateInput{Type: "generic", Payload: []byte(`{}`)})
	_, _ = svc.AssignNext(context.Background(), AssignNextInput{AgentID: "agent-1"})
	if _, err := svc.Progress(context.Background(), ProgressInput{ID: created.ID, Progress: 101}); !errors.Is(err, ErrInvalidProgress) {
		t.Fatalf("expected invalid progress")
	}
	failed, err := svc.Fail(context.Background(), FailInput{ID: created.ID, Error: "boom"})
	if err != nil {
		t.Fatal(err)
	}
	if failed.Status != domain.StatusFailed || failed.Error != "boom" {
		t.Fatal("bad fail")
	}
}
