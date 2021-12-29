package task

import (
	"container/list"
	"context"
	"errors"
	"strconv"
	"sync"
	"time"

	"go.opencensus.io/trace"
)

var (
	ErrTaskNotFound = errors.New("task not found")
)

type InMemory struct {
	mu      sync.Mutex
	tasks   list.List
	mTask   map[string]*list.Element
	counter int
}

func NewInMemoryTaskManager() *InMemory {
	return &InMemory{
		mTask: make(map[string]*list.Element),
		tasks: *list.New(),
	}
}

func (i *InMemory) CreateTask(ctx context.Context, name string) (task Task, err error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	ctx, span := trace.StartSpan(ctx, "InMemoryTaskManager.CreateTask")
	defer span.End()

	i.counter++
	task = Task{
		Id:        strconv.Itoa(i.counter),
		Name:      name,
		CreatedAt: TimeStamp(time.Now()),
		UpdatedAt: TimeStamp(time.Now()),
	}
	ee := i.tasks.PushFront(&task)
	i.mTask[task.Id] = ee
	recordTaskCreate(context.Background())
	return task, nil
}

func (i *InMemory) DeleteTask(ctx context.Context, id string) (err error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	ctx, span := trace.StartSpan(ctx, "InMemoryTaskManager.DeleteTask")
	defer span.End()

	task, ok := i.mTask[id]
	if !ok {
		return ErrTaskNotFound
	}

	i.tasks.Remove(task)
	delete(i.mTask, id)
	recordTaskDelete(context.Background())
	return nil

}

func (i *InMemory) GetTask(ctx context.Context, id string) (_ Task, err error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	ctx, span := trace.StartSpan(ctx, "InMemoryTaskManager.GetTask")
	defer span.End()

	task, ok := i.mTask[id]
	if !ok {
		return Task{}, ErrTaskNotFound
	}
	return *task.Value.(*Task), nil
}

func (i *InMemory) UpdateTask(ctx context.Context, id, name string) (_ Task, err error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	ctx, span := trace.StartSpan(ctx, "InMemoryTaskManager.UpdateTask")
	defer span.End()

	lElement, ok := i.mTask[id]
	if !ok {
		return Task{}, ErrTaskNotFound
	}

	task := lElement.Value.(*Task)
	task.Name = name
	task.UpdatedAt = TimeStamp(time.Now())
	recordTaskUpdate(context.Background())
	return *task, nil

}

func (i *InMemory) ListTasks(ctx context.Context) []Task {
	i.mu.Lock()
	defer i.mu.Unlock()

	ctx, span := trace.StartSpan(ctx, "InMemoryTaskManager.ListTasks")
	defer span.End()

	tasks := make([]Task, i.tasks.Len())
	k := i.tasks.Len() - 1
	for t := i.tasks.Front(); t != nil; t = t.Next() {
		if t.Value != nil {
			tasks[k] = *t.Value.(*Task)
			k--
		}
	}

	return tasks
}
