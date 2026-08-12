package agent

import "time"

type Agent struct {
	ID        string
	Name      string
	Endpoint  string
	Status    Status
	TokenHash []byte
	CreatedAt time.Time
	UpdatedAt time.Time
	LastSeen  *time.Time
}

type Status string

const (
	StatusUnknown Status = "unknown"
	StatusOnline  Status = "online"
	StatusOffline Status = "offline"
)
