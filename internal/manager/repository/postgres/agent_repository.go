package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	domain "github.com/sirajul777/genieacs-platform/internal/manager/domain/agent"
	"github.com/sirajul777/genieacs-platform/internal/manager/repository/postgres/sqlc"
)

type AgentRepository struct{ queries *sqlc.Queries }

func NewAgentRepository(pool *pgxpool.Pool) *AgentRepository {
	return &AgentRepository{queries: sqlc.New(pool)}
}

func (r *AgentRepository) Create(ctx context.Context, item domain.Agent) (domain.Agent, error) {
	created, err := r.queries.CreateAgent(ctx, sqlc.CreateAgentParams{ID: item.ID, Name: item.Name, Endpoint: item.Endpoint, Status: string(item.Status)})
	if err != nil {
		return domain.Agent{}, err
	}
	return toDomain(created), nil
}

func (r *AgentRepository) GetByID(ctx context.Context, id string) (domain.Agent, error) {
	item, err := r.queries.GetAgentByID(ctx, id)
	if err != nil {
		return domain.Agent{}, err
	}
	return toDomain(item), nil
}

func (r *AgentRepository) List(ctx context.Context) ([]domain.Agent, error) {
	items, err := r.queries.ListAgents(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]domain.Agent, 0, len(items))
	for _, item := range items {
		result = append(result, toDomain(item))
	}
	return result, nil
}

func toDomain(item sqlc.Agent) domain.Agent {
	return domain.Agent{ID: item.ID, Name: item.Name, Endpoint: item.Endpoint, Status: domain.Status(item.Status), CreatedAt: item.CreatedAt, UpdatedAt: item.UpdatedAt}
}
