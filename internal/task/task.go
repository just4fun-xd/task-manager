package task

import (
	"context"
	"time"
)

type TaskStatus string

const (
	StatusNew        TaskStatus = "new"
	StatusInProgress TaskStatus = "in_progress"
	StatusDone       TaskStatus = "done"
)

type Task struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Created     time.Time  `json:"created"`
	Status      TaskStatus `json:"status"`
}

type TaskRepository interface {
	Add(ctx context.Context, task *Task) error
	GetAll(ctx context.Context) ([]Task, error)
	GetById(ctx context.Context, id int) (*Task, error)
	Update(ctx context.Context, task *Task) error
	Delete(ctx context.Context, id int) error
}

func (s TaskStatus) IsValid() bool {
	switch s {
	case StatusNew, StatusInProgress, StatusDone:
		return true
	default:
		return false
	}
}
