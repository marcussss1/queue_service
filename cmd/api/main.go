package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/marcussss1/queue_service/internal/config"
	"github.com/marcussss1/queue_service/internal/tasks/delivery/http"
	tasks_repository "github.com/marcussss1/queue_service/internal/tasks/repository"
	tasks_usecase "github.com/marcussss1/queue_service/internal/tasks/usecase"
	workers_usecase "github.com/marcussss1/queue_service/internal/workers/usecase"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("не найден .env файл")
	}
}

func main() {
	yamlPath, exists := os.LookupEnv("YAML_PATH")
	if !exists {
		log.Fatal("переменная YAML_PATH не найдена")
	}

	yamlFile, err := os.ReadFile(yamlPath)
	if err != nil {
		log.Fatal(err)
	}

	var config config.Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	e.GET("/", func(ctx echo.Context) error {
		log.Info("для добавления задачи в очередь используйте api/v1/append_task")
		log.Info("для получения списка задач используйте api/v1/tasks")
		return nil
	})

	tasksRepository := tasks_repository.NewTasksRepository()
	tasksUsecase := tasks_usecase.NewTasksUsecase(tasksRepository)
	workersUsecase := workers_usecase.NewWorkersUsecase(tasksRepository, config.Server.WorkersNum)
	http.NewTaskHandler(e, tasksUsecase)

	go workersUsecase.Run()

	e.Logger.Fatal(e.Start(config.Server.Port))
}
