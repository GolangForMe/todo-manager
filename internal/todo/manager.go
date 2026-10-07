package todo

import (
	"fmt"
)

// Manager — список задач и операции над ними
type Manager struct {
	tasks []Task
}

// Создает экземпляр структуры
func NewManager(tasks []Task) *Manager {
	return &Manager{tasks: tasks}
}

// Возвращает список задач
func (m *Manager) Tasks() []Task {
	return m.tasks
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
