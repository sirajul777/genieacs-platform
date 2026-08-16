package genieacs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/sirajul777/genieacs-platform/internal/agent/extension"
)

const Name = "genieacs"

type Extension struct {
	Command string
}

func New(command string) *Extension { return &Extension{Command: command} }

func (e *Extension) Name() string { return Name }

type request struct {
	Action string          `json:"action"`
	Args   map[string]any  `json:"args"`
}

type response struct {
	Action  string `json:"action"`
	Output  string `json:"output,omitempty"`
	Success bool   `json:"success"`
}

func (e *Extension) Execute(ctx context.Context, req extension.Request) (extension.Response, error) {
	if strings.TrimSpace(e.Command) == "" {
		return extension.Response{}, errors.New("genieacs command is required")
	}
	var in request
	if len(req.Payload) == 0 {
		return extension.Response{}, errors.New("payload is required")
	}
	if err := json.Unmarshal(req.Payload, &in); err != nil {
		return extension.Response{}, fmt.Errorf("decode genieacs request: %w", err)
	}
	if in.Action == "" {
		return extension.Response{}, errors.New("action is required")
	}
	args := []string{in.Action}
	for key, value := range in.Args {
		args = append(args, key, fmt.Sprint(value))
	}
	cmd := exec.CommandContext(ctx, e.Command, args...)
	output, err := cmd.CombinedOutput()
	out := response{Action: in.Action, Output: string(output), Success: err == nil}
	data, marshalErr := json.Marshal(out)
	if marshalErr != nil {
		return extension.Response{}, marshalErr
	}
	if err != nil {
		return extension.Response{Payload: data}, fmt.Errorf("genieacs action %s: %w", in.Action, err)
	}
	return extension.Response{Payload: data}, nil
}
