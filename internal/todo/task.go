package todo

import (
	"fmt"
	"strings"
)

// Task — задача пользователя
type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

// Validate проверяет корректность задачи
func (t Task) Validate() error {
	if t.ID <= 0 {
		return fmt.Errorf("id должен быть положительным числом, получено %d", t.ID)
	}
	if strings.TrimSpace(t.Description) == "" {
		return fmt.Errorf("описание пустое или состоит только из пробелов")
	}
	return nil
}

// ValidateTasks проверяет корректность набора задач
func ValidateTasks(tasks []Task) error {
	seen := make(map[int]struct{}, len(tasks))
	for i, t := range tasks {
		if err := t.Validate(); err != nil {
			return fmt.Errorf("задача %d (описание %q): %w", i+1, t.Description, err)
		}
		if _, dup := seen[t.ID]; dup {
			return fmt.Errorf("дубликат id %d (задача %d)", t.ID, i+1)
		}
		seen[t.ID] = struct{}{}
	}
	return nil
}
