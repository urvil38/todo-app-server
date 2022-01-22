package memory

import (
	"container/list"
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/urvil38/todo-app/internal/task"
	"go.opencensus.io/trace"
)

type TaskManager struct {
	mu      sync.Mutex
	tasks   list.List
	mTask   map[string]*list.Element
	counter int
}

func NewTaskManager() *TaskManager {
	return &TaskManager{
		mTask: make(map[string]*list.Element),
		tasks: *list.New(),
	}
}

func (i *TaskManager) CreateTask(ctx context.Context, name string) (t task.Task, err error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	ctx, span := trace.StartSpan(ctx, "memory.CreateTask")
	defer span.End()

	i.counter++
	t = task.Task{
		Id:        strconv.Itoa(i.counter),
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	ee := i.tasks.PushFront(&t)
	i.mTask[t.Id] = ee
	task.RecordTaskCreate(context.Background())
	return t, nil
}

func (i *TaskManager) DeleteTask(ctx context.Context, id string) (err error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	ctx, span := trace.StartSpan(ctx, "memory.DeleteTask")
	defer span.End()

	t, ok := i.mTask[id]
	if !ok {
		return task.ErrTaskNotFound
	}

	i.tasks.Remove(t)
	delete(i.mTask, id)
	task.RecordTaskDelete(context.Background())
	return nil

}

func (i *TaskManager) GetTask(ctx context.Context, id string) (_ task.Task, err error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	ctx, span := trace.StartSpan(ctx, "memory.GetTask")
	defer span.End()

	t, ok := i.mTask[id]
	if !ok {
		return task.Task{}, task.ErrTaskNotFound
	}
	return *t.Value.(*task.Task), nil
}

func (i *TaskManager) UpdateTask(ctx context.Context, id, name string) (_ task.Task, err error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	ctx, span := trace.StartSpan(ctx, "memory.UpdateTask")
	defer span.End()

	lElement, ok := i.mTask[id]
	if !ok {
		return task.Task{}, task.ErrTaskNotFound
	}

	t := lElement.Value.(*task.Task)
	t.Name = name
	t.UpdatedAt = time.Now()
	task.RecordTaskUpdate(context.Background())
	return *t, nil

}

func (i *TaskManager) ListTasks(ctx context.Context) ([]task.Task, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	ctx, span := trace.StartSpan(ctx, "memory.ListTasks")
	defer span.End()

	tasks := make([]task.Task, i.tasks.Len())
	k := i.tasks.Len() - 1
	for t := i.tasks.Front(); t != nil; t = t.Next() {
		if t.Value != nil {
			tasks[k] = *t.Value.(*task.Task)
			k--
		}
	}

	return tasks, nil
}
