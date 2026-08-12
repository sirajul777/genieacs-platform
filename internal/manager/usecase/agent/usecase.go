package agent

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"

	domain "github.com/sirajul777/genieacs-platform/internal/manager/domain/agent"
	"github.com/sirajul777/genieacs-platform/internal/manager/repository"
)

var (
	ErrInvalidRegistration = errors.New("agent name and endpoint are required")
	ErrUnauthorized        = errors.New("invalid agent bearer token")
)

type RegisterInput struct{ Name, Endpoint string }
type RegisterOutput struct {
	Agent domain.Agent
	Token string
}
type HeartbeatInput struct{ AgentID, Token string }

type UseCase struct {
	agents     repository.AgentRepository
	heartbeats repository.HeartbeatRepository
}

func NewUseCase(agents repository.AgentRepository, heartbeats repository.HeartbeatRepository) *UseCase {
	return &UseCase{agents: agents, heartbeats: heartbeats}
}

func (u *UseCase) Register(ctx context.Context, input RegisterInput) (RegisterOutput, error) {
	if input.Name == "" || input.Endpoint == "" {
		return RegisterOutput{}, ErrInvalidRegistration
	}
	token, hash, err := newToken()
	if err != nil {
		return RegisterOutput{}, err
	}
	item := domain.Agent{ID: newID(), Name: input.Name, Endpoint: input.Endpoint, Status: domain.StatusUnknown, TokenHash: hash}
	created, err := u.agents.Create(ctx, item)
	if err != nil {
		return RegisterOutput{}, err
	}
	return RegisterOutput{Agent: created, Token: token}, nil
}

func (u *UseCase) Heartbeat(ctx context.Context, input HeartbeatInput) error {
	item, err := u.agents.GetByID(ctx, input.AgentID)
	if err != nil {
		return err
	}
	provided := sha256.Sum256([]byte(input.Token))
	if subtle.ConstantTimeCompare(provided[:], item.TokenHash) != 1 {
		return ErrUnauthorized
	}
	if _, err := u.heartbeats.Create(ctx, input.AgentID); err != nil {
		return err
	}
	_, err = u.agents.UpdateLastSeen(ctx, input.AgentID)
	return err
}

func newToken() (string, []byte, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", nil, err
	}
	token := hex.EncodeToString(buf)
	hash := sha256.Sum256([]byte(token))
	return token, hash[:], nil
}

func newID() string {
	buf := make([]byte, 16)
	_, _ = rand.Read(buf)
	buf[6] = (buf[6] & 0x0f) | 0x40
	buf[8] = (buf[8] & 0x3f) | 0x80
	return hex.EncodeToString(buf[0:4]) + "-" + hex.EncodeToString(buf[4:6]) + "-" + hex.EncodeToString(buf[6:8]) + "-" + hex.EncodeToString(buf[8:10]) + "-" + hex.EncodeToString(buf[10:16])
}
