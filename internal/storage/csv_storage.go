package storage

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/gearmobile/todo-manager/internal/todo"
)

// csvHeader — ожидаемая первая строка CSV-файла
var csvHeader = []string{"id", "description", "done"}

// SaveCSV записывает задачи в файл в формате CSV
func SaveCSV(path string, tasks []todo.Task) error {
	records := [][]string{csvHeader}
	for _, t := range tasks {
		records = append(records, []string{
			strconv.Itoa(t.ID),
			t.Description,
			strconv.FormatBool(t.Done),
		})
	}

	return writeFileAtomic(path, func(w io.Writer) error {
		cw := csv.NewWriter(w)
		if err := cw.WriteAll(records); err != nil {
			return err
		}
		cw.Flush()
		return cw.Error()
	})
}

// LoadCSVImport читает задачи из пользовательского CSV-файла импорта
func LoadCSVImport(path string) ([]todo.Task, error) {
	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("файл импорта %q не найден: %w", path, err)
		}
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать CSV из %q: %w", path, err)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("файл импорта %q пуст: ожидается заголовок %q", path, strings.Join(csvHeader, ","))
	}
	if !slices.Equal(records[0], csvHeader) {
		return nil, fmt.Errorf("строка 1: неверный заголовок CSV: ожидается %q, получено %q",
			strings.Join(csvHeader, ","), strings.Join(records[0], ","))
	}

	// сначала полностью проверяем все строки, только потом собираем результат
	tasks := make([]todo.Task, 0, len(records)-1)
	for i, rec := range records[1:] {
		line := i + 2
		if len(rec) != 3 {
			return nil, fmt.Errorf("строка %d: ожидается 3 поля, получено %d", line, len(rec))
		}
		id, err := strconv.Atoi(rec[0])
		if err != nil {
			return nil, fmt.Errorf("строка %d: некорректный id %q: %w", line, rec[0], err)
		}
		done, err := strconv.ParseBool(rec[2])
		if err != nil {
			return nil, fmt.Errorf("строка %d: некорректное значение done %q: %w", line, rec[2], err)
		}
		tasks = append(tasks, todo.Task{
			ID:          id,
			Description: rec[1],
			Done:        done,
		})
	}
	if err := todo.ValidateTasks(tasks); err != nil {
		return nil, fmt.Errorf("некорректные данные в %q: %w", path, err)
	}
	return tasks, nil
}
