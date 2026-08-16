package job

import "time"

type Status string

const (
	StatusPending    Status = "pending"
	StatusInProgress Status = "in_progress"
	StatusCompleted  Status = "completed"
	StatusFailed     Status = "failed"
	StatusCanceled   Status = "canceled"
)

func (s Status) IsTerminal() bool {
	switch s {
	case StatusCompleted, StatusFailed, StatusCanceled:
		return true
	default:
		return false
	}
}

func (s Status) CanTransitionTo(next Status) bool {
	if s.IsTerminal() {
		return false
	}
	if s == next {
		return true
	}

	switch s {
	case StatusPending:
		return next == StatusInProgress || next == StatusCanceled
	case StatusInProgress:
		return next == StatusCompleted || next == StatusFailed || next == StatusCanceled
	default:
		return false
	}
}

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
