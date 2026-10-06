package storage

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/gearmobile/todo-manager/internal/todo"
)

// LoadJSON читает задачи из файла в формате JSON
// Если файл ещё не существует — возвращает пустой список без ошибки
func LoadJSON(path string) ([]todo.Task, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []todo.Task{}, nil // первый запуск
		}
		return nil, err
	}

	var tasks []todo.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

// SaveJSON записывает задачи в файл в формате JSON
func SaveJSON(path string, tasks []todo.Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 644)
}
