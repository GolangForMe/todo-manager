package todo

import (
	"reflect"
	"strings"
	"testing"
)

// TestAdd — добавление первой задачи и последующих: ID инкрементируется
func TestAdd(t *testing.T) {
	tests := []struct {
		name string
		// описания, добавляемые по очереди
		descriptions []string
		// ожидаемые ID в порядке добавления
		wantIDs []int
	}{
		{
			name:         "первая задача получает id 1",
			descriptions: []string{"первая"},
			wantIDs:      []int{1},
		},
		{
			name:         "id последующих задач инкрементируются",
			descriptions: []string{"первая", "вторая", "третья"},
			wantIDs:      []int{1, 2, 3},
		},
		{
			name:         "пустое описание добавляется без валидации в Add",
			descriptions: []string{""},
			wantIDs:      []int{1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewManager(nil)
			gotIDs := []int{}
			for _, desc := range tt.descriptions {
				task := m.Add(desc)
				gotIDs = append(gotIDs, task.ID)
			}
			if !reflect.DeepEqual(gotIDs, tt.wantIDs) {
				t.Errorf("Add() ID = %v, ожидается %v", gotIDs, tt.wantIDs)
			}
		})
	}
}

// TestNextIDAfterDelete — генерация ID после удаления задач
func TestNextIDAfterDelete(t *testing.T) {
	tests := []struct {
		name string
		// стартовый список задач (передаётся в NewManager)
		start []Task
		// id задачи, которую удаляем перед проверкой nextID
		deleteID int
		// ожидаемый ID следующей добавленной задачи
		wantNextID int
	}{
		{
			name: "удаление задачи с максимальным id — nextID продолжает нумерацию",
			start: []Task{
				{ID: 1, Description: "первая"},
				{ID: 2, Description: "вторая"},
				{ID: 3, Description: "третья"},
			},
			deleteID:   3,
			wantNextID: 3,
		},
		{
			name: "удаление задачи из середины — nextID считается по максимуму оставшихся",
			start: []Task{
				{ID: 1, Description: "первая"},
				{ID: 2, Description: "вторая"},
				{ID: 3, Description: "третья"},
			},
			deleteID:   2,
			wantNextID: 4,
		},
		{
			name:       "удаление единственной задачи — следующая получает id 1",
			start:      []Task{{ID: 5, Description: "единственная"}},
			deleteID:   5,
			wantNextID: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewManager(tt.start)
			if err := m.Delete(tt.deleteID); err != nil {
				t.Fatalf("Delete(%d) вернула неожиданную ошибку: %v", tt.deleteID, err)
			}
			task := m.Add("новая после удаления")
			if task.ID != tt.wantNextID {
				t.Errorf("ID новой задачи = %d, ожидается %d", task.ID, tt.wantNextID)
			}
		})
	}
}

// TestListFilters — фильтры all/done/pending: содержимое и порядок
func TestListFilters(t *testing.T) {
	start := []Task{
		{ID: 1, Description: "первая", Done: false},
		{ID: 2, Description: "вторая", Done: true},
		{ID: 3, Description: "третья", Done: false},
		{ID: 4, Description: "четвёртая", Done: true},
	}

	tests := []struct {
		name   string
		filter string
		want   []Task
	}{
		{
			name:   "фильтр all возвращает все задачи в исходном порядке",
			filter: "all",
			want: []Task{
				{ID: 1, Description: "первая", Done: false},
				{ID: 2, Description: "вторая", Done: true},
				{ID: 3, Description: "третья", Done: false},
				{ID: 4, Description: "четвёртая", Done: true},
			},
		},
		{
			name:   "фильтр done возвращает только выполненные",
			filter: "done",
			want: []Task{
				{ID: 2, Description: "вторая", Done: true},
				{ID: 4, Description: "четвёртая", Done: true},
			},
		},
		{
			name:   "фильтр pending возвращает только невыполненные",
			filter: "pending",
			want: []Task{
				{ID: 1, Description: "первая", Done: false},
				{ID: 3, Description: "третья", Done: false},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewManager(start)
			got := m.List(tt.filter)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List(%q) = %v, ожидается %v", tt.filter, got, tt.want)
			}
		})
	}
}

// TestListEmptyManager — пустой менеджер для всех фильтров возвращает пустой список
func TestListEmptyManager(t *testing.T) {
	tests := []struct {
		name   string
		filter string
	}{
		{name: "all пустой", filter: "all"},
		{name: "done пустой", filter: "done"},
		{name: "pending пустой", filter: "pending"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewManager(nil)
			got := m.List(tt.filter)
			if got == nil {
				t.Fatalf("List(%q) = nil, ожидается пустой срез", tt.filter)
			}
			if len(got) != 0 {
				t.Errorf("List(%q) = %v, ожидается пустой список", tt.filter, got)
			}
		})
	}
}

// TestComplete — завершение задач: успех, повторное завершение, несуществующий id
func TestComplete(t *testing.T) {
	start := []Task{
		{ID: 1, Description: "первая"},
		{ID: 2, Description: "вторая", Done: true},
	}

	tests := []struct {
		name    string
		id      int
		wantErr bool
		errSub  string
		// ожидаемое состояние задач после операции
		wantTasks []Task
	}{
		{
			name:    "успешное завершение невыполненной задачи",
			id:      1,
			wantErr: false,
			wantTasks: []Task{
				{ID: 1, Description: "первая", Done: true},
				{ID: 2, Description: "вторая", Done: true},
			},
		},
		{
			name:    "повторное завершение уже выполненной задачи не ломается",
			id:      2,
			wantErr: false,
			wantTasks: []Task{
				{ID: 1, Description: "первая", Done: false},
				{ID: 2, Description: "вторая", Done: true},
			},
		},
		{
			name:    "несуществующий id возвращает ошибку и не меняет состояние",
			id:      42,
			wantErr: true,
			errSub:  "не найдена",
			wantTasks: []Task{
				{ID: 1, Description: "первая", Done: false},
				{ID: 2, Description: "вторая", Done: true},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewManager(start)
			err := m.Complete(tt.id)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("Complete(%d) = nil, ожидается ошибка", tt.id)
				}
				if !strings.Contains(err.Error(), tt.errSub) {
					t.Errorf("Complete(%d) ошибка %q не содержит подстроку %q", tt.id, err.Error(), tt.errSub)
				}
			} else if err != nil {
				t.Fatalf("Complete(%d) вернула неожиданную ошибку: %v", tt.id, err)
			}
			if got := m.Tasks(); !reflect.DeepEqual(got, tt.wantTasks) {
				t.Errorf("состояние после Complete(%d) = %v, ожидается %v", tt.id, got, tt.wantTasks)
			}
		})
	}
}

// TestDelete — удаление задач: успех с сохранением порядка, несуществующий id
func TestDelete(t *testing.T) {
	start := []Task{
		{ID: 1, Description: "первая"},
		{ID: 2, Description: "вторая"},
		{ID: 3, Description: "третья"},
	}

	tests := []struct {
		name      string
		id        int
		wantErr   bool
		errSub    string
		wantTasks []Task
	}{
		{
			name:    "удаление из середины сохраняет порядок остальных",
			id:      2,
			wantErr: false,
			wantTasks: []Task{
				{ID: 1, Description: "первая"},
				{ID: 3, Description: "третья"},
			},
		},
		{
			name:    "удаление первой задачи сохраняет порядок остальных",
			id:      1,
			wantErr: false,
			wantTasks: []Task{
				{ID: 2, Description: "вторая"},
				{ID: 3, Description: "третья"},
			},
		},
		{
			name:    "несуществующий id возвращает ошибку и не меняет состояние",
			id:      99,
			wantErr: true,
			errSub:  "не найдена",
			wantTasks: []Task{
				{ID: 1, Description: "первая"},
				{ID: 2, Description: "вторая"},
				{ID: 3, Description: "третья"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewManager(start)
			err := m.Delete(tt.id)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("Delete(%d) = nil, ожидается ошибка", tt.id)
				}
				if !strings.Contains(err.Error(), tt.errSub) {
					t.Errorf("Delete(%d) ошибка %q не содержит подстроку %q", tt.id, err.Error(), tt.errSub)
				}
			} else if err != nil {
				t.Fatalf("Delete(%d) вернула неожиданную ошибку: %v", tt.id, err)
			}
			if got := m.Tasks(); !reflect.DeepEqual(got, tt.wantTasks) {
				t.Errorf("состояние после Delete(%d) = %v, ожидается %v", tt.id, got, tt.wantTasks)
			}
		})
	}
}

// TestTasksIsolation — Tasks() возвращает копию: изменения возвращённого среза не влияют на менеджер
func TestTasksIsolation(t *testing.T) {
	m := NewManager([]Task{
		{ID: 1, Description: "первая"},
		{ID: 2, Description: "вторая"},
	})

	got := m.Tasks()
	got[0].Description = "изменено извне"
	got[0].Done = true
	got = append(got, Task{ID: 999, Description: "лишняя"})

	want := []Task{
		{ID: 1, Description: "первая", Done: false},
		{ID: 2, Description: "вторая", Done: false},
	}
	if actual := m.Tasks(); !reflect.DeepEqual(actual, want) {
		t.Errorf("Tasks() после изменения внешнего среза = %v, ожидается %v", actual, want)
	}
}

// TestNewManagerIsolation — NewManager копирует входной срез: изменения оригинала не влияют на менеджер
func TestNewManagerIsolation(t *testing.T) {
	start := []Task{
		{ID: 1, Description: "первая"},
		{ID: 2, Description: "вторая"},
	}

	m := NewManager(start)
	start[0].Description = "изменено извне"
	start[0].Done = true
	start = append(start, Task{ID: 999, Description: "лишняя"})

	want := []Task{
		{ID: 1, Description: "первая", Done: false},
		{ID: 2, Description: "вторая", Done: false},
	}
	if actual := m.Tasks(); !reflect.DeepEqual(actual, want) {
		t.Errorf("Tasks() после изменения исходного среза = %v, ожидается %v", actual, want)
	}
}
