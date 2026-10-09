package todo

import (
	"fmt"
)

// Manager — список задач и операции над ними
type Manager struct {
	tasks []Task
}

// NewManager создает экземпляр структуры. Входной срез копируется:
// последующие изменения переданного среза не влияют на менеджер,
// а изменения менеджера — на исходный срез.
func NewManager(tasks []Task) *Manager {
	tasksCopy := make([]Task, len(tasks))
	copy(tasksCopy, tasks)
	return &Manager{tasks: tasksCopy}
}

// Tasks возвращает копию списка задач. Возвращается независимый срез:
// изменения в нём не влияют на состояние Manager.
func (m *Manager) Tasks() []Task {
	tasksCopy := make([]Task, len(m.tasks))
	copy(tasksCopy, m.tasks)
	return tasksCopy
}

// Add добавляет новую зщадачу в список задач
func (m *Manager) Add(description string) Task {
	t := Task{
		ID:          m.nextID(),
		Description: description,
		Done:        false,
	}
	m.tasks = append(m.tasks, t)
	return t
}

// nextID вычисляет следующий свободный ID
func (m *Manager) nextID() int {
	maxID := 0
	for _, t := range m.tasks {
		if t.ID > maxID {
			maxID = t.ID
		}
	}
	return maxID + 1
}

// List возвращает задачи по фильтру: все (all), выполненные (done), ожидание(pending)
func (m *Manager) List(filter string) []Task {
	switch filter {
	case "done":
		output := []Task{}
		for _, t := range m.tasks {
			if t.Done {
				output = append(output, t)
			}
		}
		return output
	case "pending":
		output := []Task{}
		for _, t := range m.tasks {
			if !t.Done {
				output = append(output, t)
			}
		}
		return output
	default:
		return m.tasks
	}
}

// Complete помечает задачу выполненной
func (m *Manager) Complete(id int) error {
	for i := range m.tasks {
		// обращаемся к элементу слайса по его индексу; t.Done - было бы ошибкой, тк t - это КОПИЯ элемента слайса
		if m.tasks[i].ID == id {
			m.tasks[i].Done = true
			return nil
		}
	}
	return fmt.Errorf("task с id %d не найдена", id)
}

// Delete удаляет задачу
func (m *Manager) Delete(id int) error {
	for i, t := range m.tasks {
		if t.ID == id {
			m.tasks = append(m.tasks[:i], m.tasks[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("task с id %d не найдена", id)
}
