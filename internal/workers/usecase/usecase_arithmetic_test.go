package usecase

import (
	"context"
	"github.com/marcussss1/queue_service/internal/tasks/repository"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// проверяем что функция правильно работает
func TestUsecase_Arithmetic_Success(t *testing.T) {
	tasksRepository := repository.NewTasksRepository()
	workersUsecase := NewWorkersUsecase(tasksRepository, 1)

	// добавляем задачу в очередь
	tasksRepository.AppendTask(context.TODO(), defaultTask)

	// берем задачу для воркера
	task, err := tasksRepository.GetTaskToProgress(context.TODO())
	require.NoError(t, err)

	workersUsecase.arithmetic(task)

	// смотрим что всё верно
	tasks := tasksRepository.GetTasks(context.TODO())
	require.EqualValues(t, 30, tasks[0].StartValue)
	require.EqualValues(t, 5, tasks[0].CurrentIteration)
}

// проверяем что в очереди задача меняет свое состояние при подсчете прогрессии
func TestUsecase_Arithmetic_SuccessState(t *testing.T) {
	tasksRepository := repository.NewTasksRepository()
	workersUsecase := NewWorkersUsecase(tasksRepository, 1)

	// добавляем задачу в очередь
	tasksRepository.AppendTask(context.TODO(), defaultTask)

	// берем задачу для воркера
	task, err := tasksRepository.GetTaskToProgress(context.TODO())
	require.NoError(t, err)

	task.Interval = 5 // 500 ms

	go func() {
		workersUsecase.arithmetic(task)
	}()

	// ждем время большее чем время интервала
	time.Sleep(time.Millisecond * 600)

	// смотрим что задача поменяласа состояние
	tasks := tasksRepository.GetTasks(context.TODO())
	require.EqualValues(t, 15, tasks[0].StartValue)
	require.EqualValues(t, 2, tasks[0].CurrentIteration)
}
