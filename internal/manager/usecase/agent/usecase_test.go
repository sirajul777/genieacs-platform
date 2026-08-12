package agent

import (
	"context"
	"crypto/sha256"
	"errors"
	"testing"
	"time"

	domain "github.com/sirajul777/genieacs-platform/internal/manager/domain/agent"
	"github.com/sirajul777/genieacs-platform/internal/manager/domain/heartbeat"
)

type fakeAgentRepository struct {
	items       map[string]domain.Agent
	lastSeenSet bool
}

func newFakeAgentRepository() *fakeAgentRepository {
	return &fakeAgentRepository{items: map[string]domain.Agent{}}
}
func (r *fakeAgentRepository) Create(ctx context.Context, item domain.Agent) (domain.Agent, error) {
	r.items[item.ID] = item
	return item, nil
}
func (r *fakeAgentRepository) GetByID(ctx context.Context, id string) (domain.Agent, error) {
	item, ok := r.items[id]
	if !ok {
		return domain.Agent{}, errors.New("not found")
	}
	return item, nil
}
func (r *fakeAgentRepository) List(ctx context.Context) ([]domain.Agent, error) {
	out := []domain.Agent{}
	for _, item := range r.items {
		out = append(out, item)
	}
	return out, nil
}
func (r *fakeAgentRepository) UpdateLastSeen(ctx context.Context, id string) (domain.Agent, error) {
	item := r.items[id]
	now := time.Now()
	item.LastSeen = &now
	r.items[id] = item
	r.lastSeenSet = true
	return item, nil
}

type fakeHeartbeatRepository struct{ count int }

func (r *fakeHeartbeatRepository) Create(ctx context.Context, agentID string) (heartbeat.Heartbeat, error) {
	r.count++
	return heartbeat.Heartbeat{ID: int64(r.count), AgentID: agentID, CreatedAt: time.Now()}, nil
}

func TestRegisterGeneratesTokenAndStoresOnlyHash(t *testing.T) {
	agents := newFakeAgentRepository()
	usecase := NewUseCase(agents, &fakeHeartbeatRepository{})
	out, err := usecase.Register(context.Background(), RegisterInput{Name: "node-1", Endpoint: "http://agent"})
	if err != nil {
		t.Fatal(err)
	}
	if len(out.Token) != 64 {
		t.Fatalf("expected 64 hex chars for 32-byte token, got %d", len(out.Token))
	}
	stored := agents.items[out.Agent.ID]
	if string(stored.TokenHash) == out.Token {
		t.Fatal("stored raw token instead of hash")
	}
	hash := sha256.Sum256([]byte(out.Token))
	if string(stored.TokenHash) != string(hash[:]) {
		t.Fatal("stored hash does not match returned token")
	}
}

func TestHeartbeatValidatesTokenRecordsHeartbeatAndUpdatesLastSeen(t *testing.T) {
	agents := newFakeAgentRepository()
	heartbeats := &fakeHeartbeatRepository{}
	usecase := NewUseCase(agents, heartbeats)
	registered, err := usecase.Register(context.Background(), RegisterInput{Name: "node-1", Endpoint: "http://agent"})
	if err != nil {
		t.Fatal(err)
	}
	if err := usecase.Heartbeat(context.Background(), HeartbeatInput{AgentID: registered.Agent.ID, Token: registered.Token}); err != nil {
		t.Fatal(err)
	}
	if heartbeats.count != 1 {
		t.Fatalf("expected heartbeat to be recorded")
	}
	if !agents.lastSeenSet {
		t.Fatalf("expected last_seen to be updated")
	}
}

func TestHeartbeatRejectsInvalidToken(t *testing.T) {
	agents := newFakeAgentRepository()
	heartbeats := &fakeHeartbeatRepository{}
	usecase := NewUseCase(agents, heartbeats)
	registered, err := usecase.Register(context.Background(), RegisterInput{Name: "node-1", Endpoint: "http://agent"})
	if err != nil {
		t.Fatal(err)
	}
	err = usecase.Heartbeat(context.Background(), HeartbeatInput{AgentID: registered.Agent.ID, Token: "wrong"})
	if !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
	if heartbeats.count != 0 {
		t.Fatalf("invalid token should not record heartbeat")
	}
}
