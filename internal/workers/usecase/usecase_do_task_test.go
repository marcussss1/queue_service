package usecase

import (
	"context"
	"github.com/marcussss1/queue_service/internal/models"
	"github.com/marcussss1/queue_service/internal/pkg/e"
	"github.com/marcussss1/queue_service/internal/tasks/repository"
	"github.com/stretchr/testify/require"
	"runtime"
	"sync"
	"testing"
	"time"
)

// 1 воркер, 1 задача, проверяем, что она успешно выполнилась
func TestUsecase_OneWorker_OneTask_Success(t *testing.T) {
	tasksRepository := repository.NewTasksRepository()
	workersUsecase := NewWorkersUsecase(tasksRepository, 1)

	// добавляем задачу в очередь
	tasksRepository.AppendTask(context.TODO(), defaultTask)

	// берем задачу для воркера
	task, err := tasksRepository.GetTaskToProgress(context.TODO())
	require.NoError(t, err)

	// выполняем задачу
	workersUsecase.doTask(task)

	// смотрим что задача выполнилась
	tasks := tasksRepository.GetTasks(context.TODO())
	require.EqualValues(t, models.TASK_COMPLETED, tasks[0].Status)
	require.NotEmpty(t, tasks[0].FinishedAt)

	// ждем пока задача стерется из памяти
	time.Sleep(time.Millisecond * time.Duration(defaultTask.TTL*150))

	// смотрим что задача действительно стерлась
	tasks = tasksRepository.GetTasks(context.TODO())
	require.Empty(t, tasks)

	// смотрим что не утекли горутины
	require.EqualValues(t, 2, runtime.NumGoroutine())
}

// 20 воркеров, 20 задач, проверяем, что они успешно выполнились параллельно
func TestUsecase_TwoWorkers_TwoTasks_Success(t *testing.T) {
	tasksRepository := repository.NewTasksRepository()
	workersUsecase := NewWorkersUsecase(tasksRepository, 20)

	wg := &sync.WaitGroup{}
	wg.Add(20)
	for idx := 0; idx < 20; idx++ {
		// добавляем задачу в очередь
		tasksRepository.AppendTask(context.TODO(), defaultTask)

		go func() {
			// берем задачу для воркера
			task, err := tasksRepository.GetTaskToProgress(context.TODO())
			require.NoError(t, err)

			// выполняем задачу
			workersUsecase.doTask(task)
			wg.Done()
		}()
	}

	wg.Wait()

	// смотрим что все 20 задач выполнились
	tasks := tasksRepository.GetTasks(context.TODO())
	for _, task := range tasks {
		require.EqualValues(t, models.TASK_COMPLETED, task.Status)
		require.NotEmpty(t, task.FinishedAt)
	}

	// ждем пока задачи стерутся из памяти
	time.Sleep(time.Millisecond * time.Duration(defaultTask.TTL*150))

	// смотрим что задачи действительно стерлась
	tasks = tasksRepository.GetTasks(context.TODO())
	require.Empty(t, tasks)

	// смотрим что не утекли горутины
	require.EqualValues(t, 2, runtime.NumGoroutine())
}

// 2 воркера, 4 задачи, проверяем, что у нас действительно очередь и ограничение на количество воркеров
func TestUsecase_TwoWorkers_FourTasks_Success(t *testing.T) {
	tasksRepository := repository.NewTasksRepository()
	workersUsecase := NewWorkersUsecase(tasksRepository, 2)

	for idx := 0; idx < 4; idx++ {
		// добавляем задачу в очередь
		tasksRepository.AppendTask(context.TODO(), defaultTask)

		go func() {
			// берем задачу для воркера
			task, err := tasksRepository.GetTaskToProgress(context.TODO())
			require.NoError(t, err)

			// выполняем задачу
			workersUsecase.doTask(task)
		}()
	}

	// ждем время достаточное для выполнения первых двух задач, но не всех четырех
	time.Sleep(time.Millisecond * 600)

	// смотрим что все 2 задачи выполнились, а 2 выполняются в данный момент
	var (
		completedTasks  int
		inProgressTasks int
	)
	tasks := tasksRepository.GetTasks(context.TODO())
	for _, task := range tasks {
		switch task.Status {
		case models.TASK_COMPLETED:
			completedTasks++
		case models.TASK_IN_PROGRESS:
			inProgressTasks++
		}
	}
	require.EqualValues(t, 2, completedTasks)
	require.EqualValues(t, 2, inProgressTasks)

	// ждем пока первые две задачи стерутся из памяти
	time.Sleep(time.Millisecond * time.Duration(defaultTask.TTL*150))

	// смотрим что первые две задачи действительно стерлись
	tasks = tasksRepository.GetTasks(context.TODO())
	require.EqualValues(t, 2, len(tasks))

	// смотрим что две задачи еще выполняются
	require.EqualValues(t, 4, runtime.NumGoroutine())
}

// нет задач, проверяем что пришла ошибка
func TestUsecase_OneWorker_ZeroTasks_Success(t *testing.T) {
	tasksRepository := repository.NewTasksRepository()

	// берем задачу для воркера
	_, err := tasksRepository.GetTaskToProgress(context.TODO())
	require.EqualValues(t, e.ErrNoFreeTasks.Error(), err.Error())
}
