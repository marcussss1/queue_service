package utils

import "github.com/marcussss1/queue_service/internal/models"

func FromMapToSlice(dict map[int]models.Task) []models.Task {
	var idx int
	slice := make([]models.Task, len(dict))

	for _, v := range dict {
		slice[idx] = v
		idx++
	}

	return slice
}
