package tasks

import (
	"context"
	"github.com/marcussss1/queue_service/internal/models"
)

type Repository interface {
	GetTasks(ctx context.Context) []models.Task
	AppendTask(ctx context.Context, data models.AppendTaskRequest) models.Task
	GetTaskToProgress(ctx context.Context) (models.Task, error)
	UpdateTask(ctx context.Context, task models.Task)
	DeleteTask(ctx context.Context, id int)
}
