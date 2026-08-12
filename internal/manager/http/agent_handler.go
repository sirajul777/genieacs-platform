package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	agentusecase "github.com/sirajul777/genieacs-platform/internal/manager/usecase/agent"
	"github.com/sirajul777/genieacs-platform/internal/shared/response"
)

type AgentHandler struct{ usecase *agentusecase.UseCase }

func NewAgentHandler(usecase *agentusecase.UseCase) *AgentHandler {
	return &AgentHandler{usecase: usecase}
}

type registerRequest struct {
	Name     string `json:"name"`
	Endpoint string `json:"endpoint"`
}
type heartbeatRequest struct {
	AgentID string `json:"agent_id"`
}

type agentResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Endpoint string `json:"endpoint"`
	Status   string `json:"status"`
}

func (h *AgentHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json body"})
		return
	}
	out, err := h.usecase.Register(r.Context(), agentusecase.RegisterInput{Name: req.Name, Endpoint: req.Endpoint})
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, agentusecase.ErrInvalidRegistration) {
			status = http.StatusBadRequest
		}
		response.JSON(w, status, map[string]string{"error": err.Error()})
		return
	}
	response.JSON(w, http.StatusCreated, map[string]any{
		"agent": agentResponse{
			ID:       out.Agent.ID,
			Name:     out.Agent.Name,
			Endpoint: out.Agent.Endpoint,
			Status:   string(out.Agent.Status),
		},
		"token": out.Token,
	})
}

func (h *AgentHandler) Heartbeat(w http.ResponseWriter, r *http.Request) {
	var req heartbeatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json body"})
		return
	}
	token := bearerToken(r.Header.Get("Authorization"))
	if token == "" {
		response.JSON(w, http.StatusUnauthorized, map[string]string{"error": "missing bearer token"})
		return
	}
	if err := h.usecase.Heartbeat(r.Context(), agentusecase.HeartbeatInput{AgentID: req.AgentID, Token: token}); err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, agentusecase.ErrUnauthorized) {
			status = http.StatusUnauthorized
		}
		response.JSON(w, status, map[string]string{"error": err.Error()})
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func bearerToken(header string) string {
	prefix := "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(header, prefix))
}
