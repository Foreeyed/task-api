package app

import "time"


const (
	StatusPending    string = "pending"
	StatusCompleted  string = "completed"
	StatusProcessing string = "processing"
	StatusFailed     string = "failed"
)

type Task struct {
	ID          string        `json:"id"`
	Status      string        `json:"status"`
	CreatedAt   time.Time     `json:"created_at"`
	StartedAt   *time.Time    `json:"started_at,omitempty"`
	CompletedAt *time.Time    `json:"completed_at,omitempty"`
	Duration    time.Duration `json:"duration,omitempty"`
	Result      *string       `json:"result,omitempty"`
	Error       *string       `json:"error,omitempty"`
}
