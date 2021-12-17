package task

import (
	"container/list"
	"errors"
	"strconv"
	"sync"
	"time"
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

func (i *InMemory) CreateTask(name string) (task Task, err error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.counter++
	task = Task{
		Id:        strconv.Itoa(i.counter),
		Name:      name,
		CreatedAt: TimeStamp(time.Now()),
		UpdatedAt: TimeStamp(time.Now()),
	}
	ee := i.tasks.PushFront(&task)
	i.mTask[task.Id] = ee
	return task, nil
}

func (i *InMemory) DeleteTask(id string) (err error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	task, ok := i.mTask[id]
	if !ok {
		return ErrTaskNotFound
	}

	i.tasks.Remove(task)
	delete(i.mTask, id)
	return nil

}

func (i *InMemory) GetTask(id string) (_ Task, err error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	task, ok := i.mTask[id]
	if !ok {
		return Task{}, ErrTaskNotFound
	}
	return *task.Value.(*Task), nil
}

func (i *InMemory) UpdateTask(id, name string) (_ Task, err error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	lElement, ok := i.mTask[id]
	if !ok {
		return Task{}, ErrTaskNotFound
	}

	task := lElement.Value.(*Task)
	task.Name = name
	task.UpdatedAt = TimeStamp(time.Now())

	return *task, nil

}

func (i *InMemory) ListTasks() []Task {
	i.mu.Lock()
	defer i.mu.Unlock()

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
