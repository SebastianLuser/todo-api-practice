package domain

import "time"

const (
	StatusPending    Status = "pending"
	StatusInProgress Status = "in_progress"
	StatusCompleted  Status = "completed"

	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

type (
	Status string

	Priority string

	Todo struct {
		ID          string
		Title       string
		Description string
		Status      Status
		Priority    Priority
		CreatedAt   time.Time
		UpdatedAt   time.Time
	}
)

func (s Status) IsValid() bool {
	switch s {
	case StatusPending, StatusInProgress, StatusCompleted:
		return true
	}
	return false
}

func (p Priority) IsValid() bool {
	switch p {
	case PriorityLow, PriorityMedium, PriorityHigh:
		return true
	}
	return false
}
