package repository

import (
	"context"
	"github.com/marcussss1/queue_service/internal/models"
	"github.com/marcussss1/queue_service/internal/pkg/e"
	"github.com/marcussss1/queue_service/internal/pkg/utils"
	desc "github.com/marcussss1/queue_service/internal/tasks"
	"sync"
	"time"
)

type repository struct{}

var (
	lastTaskIdx int
	freeTaskIdx int
	mutex       sync.RWMutex
	tasks       map[int]models.Task
)

func NewTasksRepository() desc.Repository {
	tasks = make(map[int]models.Task)
	return repository{}
}

func (r repository) GetTasks(ctx context.Context) []models.Task {
	mutex.RLock()
	defer mutex.RUnlock()

	return utils.FromMapToSlice(tasks)
}

func (r repository) AppendTask(ctx context.Context, data models.AppendTaskRequest) models.Task {
	task := models.Task{
		ID:               lastTaskIdx,
		Status:           models.TASK_IN_QUEUE,
		NumElements:      data.NumElements,
		Delta:            data.Delta,
		StartValue:       data.StartValue,
		Interval:         data.Interval,
		TTL:              data.TTL,
		CurrentIteration: 0,
		CreatedAt:        time.Now().String(),
		StartedAt:        "",
		FinishedAt:       "",
	}

	mutex.Lock()
	tasks[lastTaskIdx] = task
	lastTaskIdx++
	defer mutex.Unlock()

	return task
}

func (r repository) GetTaskToProgress(ctx context.Context) (models.Task, error) {
	if freeTaskIdx >= lastTaskIdx {
		return models.Task{}, e.ErrNoFreeTasks
	}

	mutex.Lock()
	r.updateTaskToProgress()
	defer func() {
		mutex.Unlock()
		freeTaskIdx++
	}()

	return tasks[freeTaskIdx], nil
}

func (r repository) UpdateTask(ctx context.Context, task models.Task) {
	mutex.Lock()
	tasks[task.ID] = task
	mutex.Unlock()
}

func (r repository) DeleteTask(ctx context.Context, id int) {
	mutex.Lock()
	delete(tasks, id)
	mutex.Unlock()
}

// потокобезопасно, так как у каждого воркера уникальная задача
func (r repository) updateTaskToProgress() {
	task := tasks[freeTaskIdx]
	task.Status = models.TASK_IN_PROGRESS
	task.StartedAt = time.Now().String()
	tasks[freeTaskIdx] = task
}
