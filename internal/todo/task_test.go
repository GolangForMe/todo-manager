package todo

import (
	"strings"
	"testing"
)

// TestValidate — табличные тесты валидации одной задачи
func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		task    Task
		wantErr bool
		errSub  string // подстрока, ожидаемая в тексте ошибки
	}{
		{
			name:    "корректная задача",
			task:    Task{ID: 1, Description: "купить молоко", Done: false},
			wantErr: false,
		},
		{
			name:    "нулевой id",
			task:    Task{ID: 0, Description: "задача"},
			wantErr: true,
			errSub:  "положительным",
		},
		{
			name:    "отрицательный id",
			task:    Task{ID: -3, Description: "задача"},
			wantErr: true,
			errSub:  "положительным",
		},
		{
			name:    "пустое описание",
			task:    Task{ID: 1, Description: ""},
			wantErr: true,
			errSub:  "описание",
		},
		{
			name:    "описание только из пробелов",
			task:    Task{ID: 1, Description: "   \t\n "},
			wantErr: true,
			errSub:  "описание",
		},
		{
			name:    "описание с пробелами по краям — допустимо",
			task:    Task{ID: 1, Description: "  задача  "},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.task.Validate()
			if tt.wantErr {
				if err == nil {
					t.Fatalf("Validate() = nil, ожидается ошибка (подстрока %q)", tt.errSub)
				}
				if tt.errSub != "" && !strings.Contains(err.Error(), tt.errSub) {
					t.Errorf("Validate() ошибка %q не содержит подстроку %q", err.Error(), tt.errSub)
				}
				return
			}
			if err != nil {
				t.Fatalf("Validate() вернула неожиданную ошибку: %v", err)
			}
		})
	}
}

// TestValidateTasks — табличные тесты валидации набора задач
func TestValidateTasks(t *testing.T) {
	tests := []struct {
		name    string
		tasks   []Task
		wantErr bool
		errSub  string
	}{
		{
			name: "корректный набор",
			tasks: []Task{
				{ID: 1, Description: "первая"},
				{ID: 2, Description: "вторая"},
			},
			wantErr: false,
		},
		{
			name:    "пустой набор допустим",
			tasks:   []Task{},
			wantErr: false,
		},
		{
			name: "задача с id <= 0",
			tasks: []Task{
				{ID: 1, Description: "первая"},
				{ID: 0, Description: "вторая"},
			},
			wantErr: true,
			errSub:  "положительным",
		},
		{
			name: "задача с пустым описанием",
			tasks: []Task{
				{ID: 1, Description: ""},
			},
			wantErr: true,
			errSub:  "описание",
		},
		{
			name: "дубликаты id",
			tasks: []Task{
				{ID: 1, Description: "первая"},
				{ID: 2, Description: "вторая"},
				{ID: 1, Description: "третья"},
			},
			wantErr: true,
			errSub:  "дубликат id 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTasks(tt.tasks)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("ValidateTasks() = nil, ожидается ошибка (подстрока %q)", tt.errSub)
				}
				if tt.errSub != "" && !strings.Contains(err.Error(), tt.errSub) {
					t.Errorf("ValidateTasks() ошибка %q не содержит подстроку %q", err.Error(), tt.errSub)
				}
				return
			}
			if err != nil {
				t.Fatalf("ValidateTasks() вернула неожиданную ошибку: %v", err)
			}
		})
	}
}
