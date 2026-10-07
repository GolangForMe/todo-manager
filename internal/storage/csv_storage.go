package storage

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/gearmobile/todo-manager/internal/todo"
)

func SaveCSV(path string, tasks []todo.Task) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	records := [][]string{{"id", "description", "done"}}
	for _, t := range tasks {
		records = append(records, []string{
			strconv.Itoa(t.ID),
			t.Description,
			strconv.FormatBool(t.Done),
		})
	}
	w := csv.NewWriter(f)
	return w.WriteAll(records)
}

func LoadCSV(path string) ([]todo.Task, error) {
	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []todo.Task{}, nil
		}
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	tasks := []todo.Task{}
	for i, rec := range records {
		if i == 0 {
			continue
		}
		id, err := strconv.Atoi(rec[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad id %q: %w", i+1, rec[0], err)
		}
		done, err := strconv.ParseBool(rec[2])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad done %q: %w", i+1, rec[2], err)
		}
		tasks = append(tasks, todo.Task{
			ID:          id,
			Description: rec[1],
			Done:        done,
		})
	}
	return tasks, nil
}
