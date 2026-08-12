package repository

import (
	"context"

	"github.com/sirajul777/genieacs-platform/internal/manager/domain/heartbeat"
)

type HeartbeatRepository interface {
	Create(ctx context.Context, agentID string) (heartbeat.Heartbeat, error)
}
