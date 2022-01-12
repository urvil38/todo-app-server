package task

import (
	"context"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

var (
	taskCreateCount = stats.Int64(
		"todo_app/task/create/count",
		"Number of tasks created",
		stats.UnitDimensionless,
	)

	taskUpdateCount = stats.Int64(
		"todo_app/task/update/count",
		"Number of tasks updated",
		stats.UnitDimensionless,
	)

	taskDeleteCount = stats.Int64(
		"todo_app/task/delete/count",
		"Number of tasks deleted",
		stats.UnitDimensionless,
	)

	TaskCreatedCountView = &view.View{
		Name:        "todo_app/task/create/count",
		Measure:     taskCreateCount,
		Aggregation: view.Count(),
		Description: "Number of tasks created",
	}
	TaskUpdatedCountView = &view.View{
		Name:        "todo_app/task/update/count",
		Measure:     taskUpdateCount,
		Aggregation: view.Count(),
		Description: "Number of tasks updated",
	}
	TaskDeletedCountView = &view.View{
		Name:        "todo_app/task/delete/count",
		Measure:     taskDeleteCount,
		Aggregation: view.Count(),
		Description: "Number of tasks deleted",
	}
)

func RecordTaskCreate(ctx context.Context) {
	stats.Record(ctx, taskCreateCount.M(1))
}

func RecordTaskUpdate(ctx context.Context) {
	stats.Record(ctx, taskUpdateCount.M(1))
}

func RecordTaskDelete(ctx context.Context) {
	stats.Record(ctx, taskDeleteCount.M(1))
}
