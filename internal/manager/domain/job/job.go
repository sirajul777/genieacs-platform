package job

import "time"

type Status string

const (
	StatusPending    Status = "pending"
	StatusInProgress Status = "in_progress"
	StatusCompleted  Status = "completed"
	StatusFailed     Status = "failed"
)

type Job struct {
	ID        string
	Type      string
	Payload   []byte
	Status    Status
	AgentID   *string
	Progress  int
	Result    []byte
	Error     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Event struct {
	ID        int64
	JobID     string
	Type      string
	Message   string
	CreatedAt time.Time
}
