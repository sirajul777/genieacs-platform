package http

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	domain "github.com/sirajul777/genieacs-platform/internal/manager/domain/agent"
	"github.com/sirajul777/genieacs-platform/internal/manager/domain/heartbeat"
	agentusecase "github.com/sirajul777/genieacs-platform/internal/manager/usecase/agent"
)

type handlerAgentRepository struct{ items map[string]domain.Agent }

func newHandlerAgentRepository() *handlerAgentRepository {
	return &handlerAgentRepository{items: map[string]domain.Agent{}}
}
func (r *handlerAgentRepository) Create(ctx context.Context, item domain.Agent) (domain.Agent, error) {
	r.items[item.ID] = item
	return item, nil
}
func (r *handlerAgentRepository) GetByID(ctx context.Context, id string) (domain.Agent, error) {
	item, ok := r.items[id]
	if !ok {
		return domain.Agent{}, errors.New("not found")
	}
	return item, nil
}
func (r *handlerAgentRepository) List(ctx context.Context) ([]domain.Agent, error) { return nil, nil }
func (r *handlerAgentRepository) UpdateLastSeen(ctx context.Context, id string) (domain.Agent, error) {
	item := r.items[id]
	now := time.Now()
	item.LastSeen = &now
	r.items[id] = item
	return item, nil
}

type handlerHeartbeatRepository struct{ count int }

func (r *handlerHeartbeatRepository) Create(ctx context.Context, agentID string) (heartbeat.Heartbeat, error) {
	r.count++
	return heartbeat.Heartbeat{ID: int64(r.count), AgentID: agentID, CreatedAt: time.Now()}, nil
}

func TestRegisterHandler(t *testing.T) {
	handler := NewAgentHandler(agentusecase.NewUseCase(newHandlerAgentRepository(), &handlerHeartbeatRepository{}))
	req := httptest.NewRequest(http.MethodPost, "/api/v1/agents/register", bytes.NewBufferString(`{"name":"node-1","endpoint":"http://agent"}`))
	rec := httptest.NewRecorder()
	handler.Register(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "token") {
		t.Fatalf("expected response to include token: %s", rec.Body.String())
	}
}

func TestHeartbeatHandler(t *testing.T) {
	agents := newHandlerAgentRepository()
	heartbeats := &handlerHeartbeatRepository{}
	usecase := agentusecase.NewUseCase(agents, heartbeats)
	handler := NewAgentHandler(usecase)
	token := "test-token"
	hash := sha256.Sum256([]byte(token))
	agents.items["agent-1"] = domain.Agent{ID: "agent-1", Name: "node-1", Endpoint: "http://agent", TokenHash: hash[:]}
	req := httptest.NewRequest(http.MethodPost, "/api/v1/agents/heartbeat", bytes.NewBufferString(`{"agent_id":"agent-1"}`))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	handler.Heartbeat(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	if heartbeats.count != 1 {
		t.Fatalf("expected heartbeat to be recorded")
	}
}

func TestHeartbeatHandlerRejectsMissingToken(t *testing.T) {
	handler := NewAgentHandler(agentusecase.NewUseCase(newHandlerAgentRepository(), &handlerHeartbeatRepository{}))
	req := httptest.NewRequest(http.MethodPost, "/api/v1/agents/heartbeat", bytes.NewBufferString(`{"agent_id":"agent-1"}`))
	rec := httptest.NewRecorder()
	handler.Heartbeat(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}
