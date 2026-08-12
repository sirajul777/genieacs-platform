package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	domain "github.com/sirajul777/genieacs-platform/internal/manager/domain/heartbeat"
	"github.com/sirajul777/genieacs-platform/internal/manager/repository/postgres/sqlc"
)

type HeartbeatRepository struct{ queries *sqlc.Queries }

func NewHeartbeatRepository(pool *pgxpool.Pool) *HeartbeatRepository {
	return &HeartbeatRepository{queries: sqlc.New(pool)}
}

func (r *HeartbeatRepository) Create(ctx context.Context, agentID string) (domain.Heartbeat, error) {
	created, err := r.queries.CreateHeartbeat(ctx, agentID)
	if err != nil {
		return domain.Heartbeat{}, err
	}
	return domain.Heartbeat{ID: created.ID, AgentID: created.AgentID, CreatedAt: created.CreatedAt}, nil
}
