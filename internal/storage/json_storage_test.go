package storage

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gearmobile/todo-manager/internal/todo"
)

// writeTestFile — вспомогательная функция: записывает содержимое во временный файл
func writeTestFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("не удалось создать тестовый файл %q: %v", path, err)
	}
	return path
}

// TestLoadJSON — чтение JSON: отсутствие файла, null, повреждённые данные, некорректные типы
func TestLoadJSON(t *testing.T) {
	tests := []struct {
		name string
		// содержимое файла; пустая строка — файл не создавать
		content string
		want    []todo.Task
		wantErr bool
		errSub  string
	}{
		{
			name:    "отсутствующий файл возвращает пустой список без ошибки",
			content: "",
			want:    []todo.Task{},
			wantErr: false,
		},
		{
			name:    "пустой массив читается как пустой список",
			content: "[]",
			want:    []todo.Task{},
			wantErr: false,
		},
		{
			name:    "json null возвращает пустой неп-nil срез",
			content: "null",
			want:    []todo.Task{},
			wantErr: false,
		},
		{
			name:    "валидный список задач читается идентично",
			content: `[{"id":1,"description":"первая","done":true},{"id":2,"description":"вторая","done":false}]`,
			want: []todo.Task{
				{ID: 1, Description: "первая", Done: true},
				{ID: 2, Description: "вторая", Done: false},
			},
			wantErr: false,
		},
		{
			name:    "повреждённый json возвращает ошибку декодирования",
			content: `{"id":1,"description":"незакрытый`,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "усечённый json возвращает ошибку",
			content: `[{"id":1,"description":"пер`,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "done строкой вместо bool возвращает ошибку",
			content: `[{"id":1,"description":"задача","done":"yes"}]`,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "id строкой вместо числа возвращает ошибку",
			content: `[{"id":"один","description":"задача","done":false}]`,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "tasks.json")
			if tt.content != "" {
				path = writeTestFile(t, dir, "tasks.json", tt.content)
			}

			got, err := LoadJSON(path)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("LoadJSON() = %v, nil; ожидается ошибка (подстрока %q)", got, tt.errSub)
				}
				if tt.errSub != "" && !strings.Contains(err.Error(), tt.errSub) {
					t.Errorf("LoadJSON() ошибка %q не содержит подстроку %q", err.Error(), tt.errSub)
				}
				return
			}
			if err != nil {
				t.Fatalf("LoadJSON() вернула неожиданную ошибку: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadJSON() = %v, ожидается %v", got, tt.want)
			}
		})
	}
}

// TestSaveLoadJSONRoundTrip — сохранение и чтение возвращают идентичные данные
func TestSaveLoadJSONRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		tasks []todo.Task
	}{
		{
			name:  "пустой список",
			tasks: []todo.Task{},
		},
		{
			name: "несколько задач",
			tasks: []todo.Task{
				{ID: 1, Description: "первая", Done: true},
				{ID: 2, Description: "вторая", Done: false},
				{ID: 10, Description: "описание с, запятой и \"кавычками\"", Done: false},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "tasks.json")

			if err := SaveJSON(path, tt.tasks); err != nil {
				t.Fatalf("SaveJSON() вернула неожиданную ошибку: %v", err)
			}
			got, err := LoadJSON(path)
			if err != nil {
				t.Fatalf("LoadJSON() вернула неожиданную ошибку: %v", err)
			}
			if !reflect.DeepEqual(got, tt.tasks) {
				t.Errorf("round-trip: LoadJSON() = %v, ожидается %v", got, tt.tasks)
			}
		})
	}
}

// TestSaveJSONNoTempFiles — writeFileAtomic не оставляет временных файлов
func TestSaveJSONNoTempFiles(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tasks.json")

	if err := SaveJSON(path, []todo.Task{{ID: 1, Description: "задача"}}); err != nil {
		t.Fatalf("SaveJSON() вернула неожиданную ошибку: %v", err)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("не удалось прочитать каталог %q: %v", dir, err)
	}
	for _, e := range entries {
		name := e.Name()
		if name != "tasks.json" {
			t.Errorf("в каталоге остались посторонние (временные) файлы: %q", name)
		}
	}
}
