package repository

import (
	"context"

	"github.com/sirajul777/genieacs-platform/internal/manager/domain/agent"
)

type AgentRepository interface {
	Create(ctx context.Context, item agent.Agent) (agent.Agent, error)
	GetByID(ctx context.Context, id string) (agent.Agent, error)
	List(ctx context.Context) ([]agent.Agent, error)
}
