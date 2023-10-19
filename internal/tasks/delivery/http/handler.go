package http

import (
	"github.com/labstack/echo/v4"
	desc "github.com/marcussss1/queue_service/internal/tasks"
)

type handler struct {
	taskUsecase desc.Usecase
}

func NewTaskHandler(e *echo.Echo, taskUsecase desc.Usecase) handler {
	h := handler{taskUsecase: taskUsecase}

	e.GET("api/v1/tasks", h.GetTasksHandler)
	e.POST("api/v1/append_task", h.AppendTaskHandler)

	return h
}
