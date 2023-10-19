package usecase

import (
	"context"

	"github.com/marcussss1/queue_service/internal/models"
	desc "github.com/marcussss1/queue_service/internal/tasks"
)

type usecase struct {
	taskRepository desc.Repository
}

func NewTasksUsecase(taskRepository desc.Repository) desc.Usecase {
	return usecase{taskRepository: taskRepository}
}

func (u usecase) GetTasks(ctx context.Context) []models.Task {
	return u.taskRepository.GetTasks(ctx)
}

func (u usecase) AppendTask(ctx context.Context, data models.AppendTaskRequest) models.Task {
	return u.taskRepository.AppendTask(ctx, data)
}
