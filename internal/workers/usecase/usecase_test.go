package usecase

import "github.com/marcussss1/queue_service/internal/models"

var (
	defaultTask = models.AppendTaskRequest{
		NumElements: 5,
		Delta:       5,
		StartValue:  5,
		Interval:    1, // 100 ms
		TTL:         1, // 100 ms
	}
)
