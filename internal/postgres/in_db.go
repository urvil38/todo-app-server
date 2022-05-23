package postgres

import (
	"context"
	"database/sql"

	"github.com/urvil38/todo-app/internal/config"
	"github.com/urvil38/todo-app/internal/database"
	"github.com/urvil38/todo-app/internal/log"
	"github.com/urvil38/todo-app/internal/task"
	"go.opencensus.io/trace"
)

type TaskManager struct {
	db *DB
}

func NewTaskManager(ctx context.Context, cfg config.Config) *TaskManager {
	db, err := OpenDB(ctx, &cfg)
	if err != nil {
		log.Logger.Fatal(err)
	}
	return &TaskManager{
		db: db,
	}
}

func (tm *TaskManager) CreateTask(ctx context.Context, name string) (_ task.Task, err error) {
	ctx, span := trace.StartSpan(ctx, "db.CreateTask")
	defer span.End()

	taskArgs := database.StructScanner(task.Task{})
	var t task.Task

	err = tm.db.db.QueryRow(ctx, `
	INSERT INTO tasks(
		name)
	VALUES ($1)
	RETURNING id, name, created_at, updated_at
	`, name).Scan(taskArgs(&t)...)
	if err != nil {
		return t, err
	}
	task.RecordTaskCreate(ctx)
	return t, nil
}

func (tm *TaskManager) UpdateTask(ctx context.Context, id, name string) (task.Task, error) {
	ctx, span := trace.StartSpan(ctx, "db.UpdateTask")
	defer span.End()

	taskArgs := database.StructScanner(task.Task{})
	var t task.Task

	err := tm.db.db.QueryRow(ctx, "UPDATE tasks SET name = $1 WHERE id = $2 RETURNING id, name, created_at, updated_at", name, id).Scan(taskArgs(&t)...)
	if err != nil {
		if err == sql.ErrNoRows {
			return t, task.ErrTaskNotFound
		} else {
			return t, err
		}
	}
	task.RecordTaskUpdate(ctx)
	return t, nil
}

func (tm *TaskManager) DeleteTask(ctx context.Context, id string) error {
	ctx, span := trace.StartSpan(ctx, "db.DeleteTask")
	defer span.End()

	taskArgs := database.StructScanner(task.Task{})
	var t task.Task

	err := tm.db.db.QueryRow(ctx, "DELETE FROM tasks WHERE id = $1 RETURNING id, name, created_at, updated_at", id).Scan(taskArgs(&t)...)
	if err != nil {
		if err == sql.ErrNoRows {
			return task.ErrTaskNotFound
		} else {
			return err
		}
	}
	
	task.RecordTaskDelete(ctx)
	return nil
}

func (tm *TaskManager) ListTasks(ctx context.Context) ([]task.Task, error) {
	ctx, span := trace.StartSpan(ctx, "db.ListTasks")
	defer span.End()

	var tasks []task.Task

	collect := func(rows *sql.Rows) error {
		var t task.Task
		taskArgs := database.StructScanner(task.Task{})
		if err := rows.Scan(taskArgs(&t)...); err != nil {
			return err
		}
		tasks = append(tasks, t)
		return nil
	}

	err := tm.db.db.RunQueryIncrementally(ctx, "SELECT id, name, created_at, updated_at FROM tasks", 5000, collect)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (tm *TaskManager) GetTask(ctx context.Context, id string) (task.Task, error) {
	ctx, span := trace.StartSpan(ctx, "db.GetTask")
	defer span.End()

	taskArgs := database.StructScanner(task.Task{})
	var t task.Task

	err := tm.db.db.QueryRow(ctx, "SELECT id, name, created_at, updated_at FROM tasks WHERE id = $1", id).Scan(taskArgs(&t)...)
	if err != nil {
		if err == sql.ErrNoRows {
			return t, task.ErrTaskNotFound
		} else {
			return t, err
		}
	}

	return t, nil
}
