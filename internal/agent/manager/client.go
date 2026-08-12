package manager

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var ErrNoJob = errors.New("no job available")

type Client struct {
	baseURL string
	token   string
	http    *http.Client
}

type Job struct {
	ID      string          `json:"id"`
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
	Status  string          `json:"status"`
}

type jobResponse struct {
	ID      string          `json:"id"`
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
	Status  string          `json:"status"`
}

func NewClient(baseURL, token string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}
	return &Client{baseURL: strings.TrimRight(baseURL, "/"), token: token, http: httpClient}
}

func (c *Client) Poll(ctx context.Context, agentID string) (Job, error) {
	var out jobResponse
	err := c.post(ctx, "/api/v1/jobs/poll", map[string]string{"agent_id": agentID}, &out)
	if err != nil {
		if errors.Is(err, ErrNoJob) {
			return Job{}, ErrNoJob
		}
		return Job{}, err
	}
	return Job{ID: out.ID, Type: out.Type, Payload: out.Payload, Status: out.Status}, nil
}

func (c *Client) Progress(ctx context.Context, jobID string, progress int, message string) error {
	return c.post(ctx, "/api/v1/jobs/"+jobID+"/progress", map[string]any{"progress": progress, "message": message}, nil)
}

func (c *Client) Complete(ctx context.Context, jobID string, result any) error {
	return c.post(ctx, "/api/v1/jobs/"+jobID+"/complete", map[string]any{"result": result}, nil)
}

func (c *Client) Fail(ctx context.Context, jobID, message string) error {
	return c.post(ctx, "/api/v1/jobs/"+jobID+"/fail", map[string]string{"error": message}, nil)
}

func (c *Client) post(ctx context.Context, path string, payload any, out any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode == http.StatusNotFound {
		return ErrNoJob
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("manager request %s: status=%d body=%s", path, resp.StatusCode, strings.TrimSpace(string(data)))
	}
	if out == nil || len(data) == 0 {
		return nil
	}
	if err := json.Unmarshal(data, out); err != nil {
		return fmt.Errorf("decode manager response: %w", err)
	}
	return nil
}
