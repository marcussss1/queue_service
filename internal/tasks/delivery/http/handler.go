package http

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/marcussss1/queue_service/internal/models"
	desc "github.com/marcussss1/queue_service/internal/tasks"
	"net/http"
)

type handler struct {
	taskUsecase desc.Usecase
}

func (h handler) GetTasksHandler(ctx echo.Context) error {
	tasks := h.taskUsecase.GetTasks(ctx.Request().Context())
	return ctx.JSON(http.StatusOK, tasks)
}

func (h handler) AppendTaskHandler(ctx echo.Context) error {
	var data models.AppendTaskRequest

	err := ctx.Bind(&data)
	if err != nil {
		return fmt.Errorf("ctx.Bind: %w", err)
	}

	task := h.taskUsecase.AppendTask(ctx.Request().Context(), data)
	return ctx.JSON(http.StatusCreated, task)
}

func NewTaskHandler(e *echo.Echo, taskUsecase desc.Usecase) handler {
	h := handler{taskUsecase: taskUsecase}

	e.GET("api/v1/tasks", h.GetTasksHandler)
	e.POST("api/v1/append_task", h.AppendTaskHandler)

	return h
}
