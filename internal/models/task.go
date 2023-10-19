package models

const (
	// TASK_IN_QUEUE Задача в очереди
	TASK_IN_QUEUE = iota
	// TASK_IN_PROGRESS Задача в процессе
	TASK_IN_PROGRESS
	// TASK_COMPLETED Задача завершена
	TASK_COMPLETED
)

type AppendTaskRequest struct {
	NumElements int     `json:"n"`   // количество элементов
	Delta       float64 `json:"d"`   // дельта между элементами последовательности
	StartValue  float64 `json:"n1"`  // стартовое значение
	Interval    float64 `json:"I"`   // интервал в секундах между итерациями
	TTL         float64 `json:"TTL"` // время хранения результата в секундах
}

type Task struct {
	ID               int     `json:"id"`                     // номер в очереди
	Status           int     `json:"status"`                 // статус: В процессе/В очереди/Завершена
	NumElements      int     `json:"n"`                      // количество элементов
	Delta            float64 `json:"d"`                      // дельта между элементами последовательности
	StartValue       float64 `json:"n1"`                     // стартовое значение
	Interval         float64 `json:"I"`                      // интервал в секундах между итерациями
	TTL              float64 `json:"TTL"`                    // время хранения результата в секундах
	CurrentIteration int     `json:"current_iteration"`      // текущая итерация
	CreatedAt        string  `json:"created_at"`             // время постановки задачи
	StartedAt        string  `json:"started_at,omitempty"`   // время старта задачи
	FinishedAt       string  `json:"completed_at,omitempty"` // время окончания задачи
}
