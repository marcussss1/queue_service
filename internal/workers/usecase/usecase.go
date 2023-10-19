package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/marcussss1/queue_service/internal/models"
	"github.com/marcussss1/queue_service/internal/pkg/e"
	"github.com/marcussss1/queue_service/internal/tasks"
	log "github.com/sirupsen/logrus"
)

type usecase struct {
	workers        chan struct{}
	taskRepository tasks.Repository
}

func NewWorkersUsecase(taskRepository tasks.Repository, workersNum int) usecase {
	return usecase{workers: make(chan struct{}, workersNum), taskRepository: taskRepository}
}

func (u usecase) Run() {
	for {
		task, err := u.taskRepository.GetTaskToProgress(context.TODO())
		if errors.Is(err, e.ErrNoFreeTasks) {
			time.Sleep(time.Second) // условное время сна
			continue
		}

		go u.doTask(task)
	}
}

func (u usecase) doTask(task models.Task) {
	log.Info(fmt.Sprintf("[ЗАДАЧА: %d] [СТАТУС: В РАБОТЕ]", task.ID))

	u.workers <- struct{}{}
	defer func() {
		<-u.workers
		log.Info(fmt.Sprintf("[ЗАДАЧА: %d] [СТАТУС: ЗАВЕРШЕНА]", task.ID))
	}()

	task = u.arithmetic(task)
	task.Status = models.TASK_COMPLETED
	task.FinishedAt = time.Now().String()

	u.taskRepository.UpdateTask(context.TODO(), task)
	u.deleteTaskViaTTL(task)
}

func (u usecase) arithmetic(task models.Task) models.Task {
	for i := 0; i < task.NumElements; i++ {
		task.StartValue += task.Delta
		task.CurrentIteration++
		u.taskRepository.UpdateTask(context.TODO(), task)
		time.Sleep(time.Millisecond * time.Duration(task.Interval*100))
	}

	return task
}

func (u usecase) deleteTaskViaTTL(task models.Task) {
	time.AfterFunc(time.Millisecond*time.Duration(task.TTL*100), func() {
		u.taskRepository.DeleteTask(context.TODO(), task.ID)
		log.Info(fmt.Sprintf("[ЗАДАЧА: %d] [СТАТУС: УДАЛЕНА ИЗ ОЧЕРЕДИ]", task.ID))
	})
}
