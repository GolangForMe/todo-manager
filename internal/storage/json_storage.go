package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

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

// LoadJSONImport читает задачи из пользовательского файла импорта в формате JSON
func LoadJSONImport(path string) ([]todo.Task, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("файл импорта %q не найден: %w", path, err)
		}
		return nil, err
	}

	var tasks []todo.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("не удалось разобрать JSON из %q: %w", path, err)
	}
	return tasks, nil
}

// SaveJSON записывает задачи в файл в формате JSON.
func SaveJSON(path string, tasks []todo.Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	return writeFileAtomic(path, func(w io.Writer) error {
		_, err := w.Write(data)
		return err
	})
}

// writeFileAtomic записывает данные через функцию write во временный файл
func writeFileAtomic(path string, write func(io.Writer) error) error {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, filepath.Base(path)+".tmp-*")
	if err != nil {
		return fmt.Errorf("не удалось создать временный файл для %q: %w", path, err)
	}
	tmpName := tmp.Name()

	if err := write(tmp); err != nil {
		return errors.Join(err, tmp.Close(), os.Remove(tmpName))
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpName)
		return fmt.Errorf("не удалось закрыть временный файл для %q: %w", path, err)
	}
	if err := os.Rename(tmpName, path); err != nil {
		os.Remove(tmpName)
		return fmt.Errorf("не удалось заменить файл %q: %w", path, err)
	}
	return nil
}
