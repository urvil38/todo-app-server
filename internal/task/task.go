package task

import (
	"context"
	"fmt"
	"time"
)

type Task struct {
	Id        string    `json:"id"`
	Name      string    `json:"task_name"`
	CreatedAt TimeStamp `json:"created_at"`
	UpdatedAt TimeStamp `json:"updated_at"`
}

type TimeStamp time.Time

func (t TimeStamp) MarshalJSON() ([]byte, error) {
	stemp := fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02T15:04:05"))
	return []byte(stemp), nil
}

type Manager interface {
	CreateTask(ctx context.Context, name string) (Task, error)
	UpdateTask(ctx context.Context, id, name string) (Task, error)
	DeleteTask(ctx context.Context, id string) error
	ListTasks(ctx context.Context) []Task
	GetTask(ctx context.Context, id string) (Task, error)
}
