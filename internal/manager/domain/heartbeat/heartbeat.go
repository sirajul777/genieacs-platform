package heartbeat

import "time"

type Heartbeat struct {
	ID        int64
	AgentID   string
	CreatedAt time.Time
}
