package tasks

import (
	"context"

	"github.com/marcussss1/queue_service/internal/models"
)

type Usecase interface {
	GetTasks(ctx context.Context) []models.Task
	AppendTask(ctx context.Context, data models.AppendTaskRequest) models.Task
}
