package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	jobusecase "github.com/sirajul777/genieacs-platform/internal/manager/usecase/job"
	"github.com/sirajul777/genieacs-platform/internal/shared/response"
)

type JobHandler struct{ service *jobusecase.Service }

func NewJobHandler(service *jobusecase.Service) *JobHandler { return &JobHandler{service: service} }

type createJobRequest struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}
type pollJobRequest struct {
	AgentID string `json:"agent_id"`
}
type progressJobRequest struct {
	Progress int `json:"progress"`
}
type completeJobRequest struct {
	Result json.RawMessage `json:"result"`
}
type failJobRequest struct {
	Error string `json:"error"`
}

func (h *JobHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json body"})
		return
	}
	job, err := h.service.Create(r.Context(), jobusecase.CreateInput{Type: req.Type, Payload: defaultJSON(req.Payload)})
	writeJob(w, http.StatusCreated, job, err)
}
func (h *JobHandler) Poll(w http.ResponseWriter, r *http.Request) {
	var req pollJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json body"})
		return
	}
	job, err := h.service.AssignNext(r.Context(), jobusecase.AssignNextInput{AgentID: req.AgentID})
	writeJob(w, http.StatusOK, job, err)
}
func (h *JobHandler) Progress(w http.ResponseWriter, r *http.Request) {
	var req progressJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json body"})
		return
	}
	job, err := h.service.Progress(r.Context(), jobusecase.ProgressInput{ID: jobID(r.URL.Path), Progress: req.Progress})
	writeJob(w, http.StatusOK, job, err)
}
func (h *JobHandler) Complete(w http.ResponseWriter, r *http.Request) {
	var req completeJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json body"})
		return
	}
	job, err := h.service.Complete(r.Context(), jobusecase.CompleteInput{ID: jobID(r.URL.Path), Result: defaultJSON(req.Result)})
	writeJob(w, http.StatusOK, job, err)
}
func (h *JobHandler) Fail(w http.ResponseWriter, r *http.Request) {
	var req failJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json body"})
		return
	}
	job, err := h.service.Fail(r.Context(), jobusecase.FailInput{ID: jobID(r.URL.Path), Error: req.Error})
	writeJob(w, http.StatusOK, job, err)
}

func writeJob(w http.ResponseWriter, status int, body any, err error) {
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, jobusecase.ErrInvalidJob) || errors.Is(err, jobusecase.ErrInvalidProgress) {
			code = http.StatusBadRequest
		}
		response.JSON(w, code, map[string]string{"error": err.Error()})
		return
	}
	response.JSON(w, status, body)
}
func defaultJSON(raw json.RawMessage) []byte {
	if len(raw) == 0 {
		return []byte(`{}`)
	}
	return raw
}
func jobID(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 4 {
		return parts[3]
	}
	return ""
}
