package task

import (
	"context"
	"errors"
	"time"
)

var (
	ErrTaskNotFound = errors.New("task not found")
)

type Task struct {
	Id        string    `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type Manager interface {
	TaskCreator
	TaskUpdater
	TaskDeleter
	TaskGetter
}

type TaskCreator interface {
	CreateTask(ctx context.Context, name string) (Task, error)
}

type TaskUpdater interface {
	UpdateTask(ctx context.Context, id, name string) (Task, error)
}

type TaskDeleter interface {
	DeleteTask(ctx context.Context, id string) error
}

type TaskGetter interface {
	GetTask(ctx context.Context, id string) (Task, error)
	ListTasks(ctx context.Context) ([]Task, error)
}
