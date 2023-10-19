package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// @Summary		Get Tasks
// @Description	Get Tasks
// @Tags			Tasks
// @Accept			json
// @Produce		json
// @Success		200	{object}	models.Task
// @Router			/api/v1/tasks [get]
func (h handler) GetTasksHandler(ctx echo.Context) error {
	tasks := h.taskUsecase.GetTasks(ctx.Request().Context())
	return ctx.JSON(http.StatusOK, tasks)
}
