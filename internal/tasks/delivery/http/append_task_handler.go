package http

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/marcussss1/queue_service/internal/models"
)

// @Summary		Append Task
// @Tags			Tasks
// @Description	Append Task
// @Accept			json
// @Produce		json
// @Param			data	body		models.AppendTaskRequest true "Task"
// @Success		201		{object}	models.Task
// @Failure		500		{object}	error
// @Router			/api/v1/append_task [post]
func (h handler) AppendTaskHandler(ctx echo.Context) error {
	var data models.AppendTaskRequest

	err := ctx.Bind(&data)
	if err != nil {
		return fmt.Errorf("ctx.Bind: %w", err)
	}

	task := h.taskUsecase.AppendTask(ctx.Request().Context(), data)
	return ctx.JSON(http.StatusCreated, task)
}
