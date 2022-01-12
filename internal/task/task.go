package task

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrTaskNotFound = errors.New("task not found")
)

type Task struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt TimeStamp `json:"created_at"`
	UpdatedAt TimeStamp `json:"updated_at"`
}

type TimeStamp time.Time

func (t TimeStamp) MarshalJSON() ([]byte, error) {
	stemp := fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02T15:04:05"))
	return []byte(stemp), nil
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
	ListTasks(ctx context.Context) ([]Task,error)
}
