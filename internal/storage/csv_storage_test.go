package storage

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gearmobile/todo-manager/internal/todo"
)

// TestLoadJSONImport — строгий импорт JSON: файл обязан существовать, данные валидируются
func TestLoadJSONImport(t *testing.T) {
	tests := []struct {
		name    string
		content string // пустая строка — файл не создавать
		want    []todo.Task
		wantErr bool
		errSub  string
	}{
		{
			name:    "отсутствующий файл возвращает ошибку «файл импорта не найден»",
			content: "",
			want:    nil,
			wantErr: true,
			errSub:  "не найден",
		},
		{
			name:    "валидный файл возвращает задачи",
			content: `[{"id":1,"description":"первая","done":true}]`,
			want: []todo.Task{
				{ID: 1, Description: "первая", Done: true},
			},
			wantErr: false,
		},
		{
			name:    "json null возвращает пустой неп-nil список",
			content: "null",
			want:    []todo.Task{},
			wantErr: false,
		},
		{
			name:    "невалидные задачи (дубликат id) валидируются и отклоняются",
			content: `[{"id":1,"description":"первая","done":false},{"id":1,"description":"вторая","done":false}]`,
			want:    nil,
			wantErr: true,
			errSub:  "дубликат id 1",
		},
		{
			name:    "повреждённый json возвращает ошибку",
			content: `[{"id":1,`,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "import.json")
			if tt.content != "" {
				path = writeTestFile(t, dir, "import.json", tt.content)
			}

			got, err := LoadJSONImport(path)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("LoadJSONImport() = %v, nil; ожидается ошибка (подстрока %q)", got, tt.errSub)
				}
				if !strings.Contains(err.Error(), tt.errSub) {
					t.Errorf("LoadJSONImport() ошибка %q не содержит подстроку %q", err.Error(), tt.errSub)
				}
				// при ошибке задачи не теряются частично — результат пустой
				if got != nil {
					t.Errorf("LoadJSONImport() при ошибке вернула %v, ожидается nil (без частичного результата)", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("LoadJSONImport() вернула неожиданную ошибку: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadJSONImport() = %v, ожидается %v", got, tt.want)
			}
		})
	}
}

// TestLoadCSVImport — импорт CSV: заголовок, число полей, конвертация, валидация
func TestLoadCSVImport(t *testing.T) {
	tests := []struct {
		name    string
		content string // пустая строка — файл не создавать
		want    []todo.Task
		wantErr bool
		errSub  string
	}{
		{
			name:    "корректный csv возвращает задачи",
			content: "id,description,done\n1,первая,true\n2,вторая,false\n",
			want: []todo.Task{
				{ID: 1, Description: "первая", Done: true},
				{ID: 2, Description: "вторая", Done: false},
			},
			wantErr: false,
		},
		{
			name:    "отсутствующий файл возвращает ошибку «файл импорта не найден»",
			content: "",
			want:    nil,
			wantErr: true,
			errSub:  "не найден",
		},
		{
			name:    "неверный заголовок возвращает ошибку с указанием строки 1",
			content: "id,task,done\n1,первая,true\n",
			want:    nil,
			wantErr: true,
			errSub:  "строка 1",
		},
		{
			name:    "пустой файл без заголовка возвращает ошибку",
			content: "",
			want:    nil,
			wantErr: true,
		},
		{
			// csv.Reader по умолчанию (FieldsPerRecord) сам отклоняет строку
			// с числом полей, отличным от заголовка, ещё на этапе ReadAll;
			// явная проверка len(rec) != 3 в LoadCSVImport — защитный слой,
			// поэтому достаточно проверить, что ошибка возвращается без паники.
			name:    "запись с 2 полями возвращает ошибку без паники",
			content: "id,description,done\n1,только два поля\n",
			want:    nil,
			wantErr: true,
			errSub:  "",
		},
		{
			name:    "некорректный id (abc) возвращает ошибку с номером строки",
			content: "id,description,done\nabc,первая,true\n",
			want:    nil,
			wantErr: true,
			errSub:  "строка 2",
		},
		{
			name:    "некорректное done (maybe) возвращает ошибку с номером строки",
			content: "id,description,done\n1,первая,maybe\n",
			want:    nil,
			wantErr: true,
			errSub:  "строка 2",
		},
		{
			name:    "id <= 0 отклоняется валидацией",
			content: "id,description,done\n0,нулевая,false\n",
			want:    nil,
			wantErr: true,
			errSub:  "положительным",
		},
		{
			name:    "пустое описание отклоняется валидацией",
			content: "id,description,done\n1,,false\n",
			want:    nil,
			wantErr: true,
			errSub:  "описание",
		},
		{
			name:    "дубликаты id отклоняются валидацией",
			content: "id,description,done\n1,первая,false\n1,вторая,false\n",
			want:    nil,
			wantErr: true,
			errSub:  "дубликат id 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "import.csv")
			if tt.content != "" {
				path = writeTestFile(t, dir, "import.csv", tt.content)
			}

			got, err := LoadCSVImport(path)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("LoadCSVImport() = %v, nil; ожидается ошибка (подстрока %q)", got, tt.errSub)
				}
				if tt.errSub != "" && !strings.Contains(err.Error(), tt.errSub) {
					t.Errorf("LoadCSVImport() ошибка %q не содержит подстроку %q", err.Error(), tt.errSub)
				}
				// при ошибке частичный результат не возвращается
				if got != nil {
					t.Errorf("LoadCSVImport() при ошибке вернула %v, ожидается nil (без частичного результата)", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("LoadCSVImport() вернула неожиданную ошибку: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadCSVImport() = %v, ожидается %v", got, tt.want)
			}
		})
	}
}

// TestSaveCSVLoadCSVImportRoundTrip — сохранение в CSV и импорт возвращают идентичные данные
func TestSaveCSVLoadCSVImportRoundTrip(t *testing.T) {
	tasks := []todo.Task{
		{ID: 1, Description: "первая", Done: true},
		{ID: 2, Description: "описание с, запятой и \"кавычками\"", Done: false},
		{ID: 3, Description: "третья", Done: false},
	}

	path := filepath.Join(t.TempDir(), "tasks.csv")

	if err := SaveCSV(path, tasks); err != nil {
		t.Fatalf("SaveCSV() вернула неожиданную ошибку: %v", err)
	}
	got, err := LoadCSVImport(path)
	if err != nil {
		t.Fatalf("LoadCSVImport() вернула неожиданную ошибку: %v", err)
	}
	if !reflect.DeepEqual(got, tasks) {
		t.Errorf("round-trip: LoadCSVImport() = %v, ожидается %v", got, tasks)
	}
}

// TestSaveCSVEmptyList — сохранение пустого списка: файл с заголовком и без строк данных
func TestSaveCSVEmptyList(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tasks.csv")

	if err := SaveCSV(path, []todo.Task{}); err != nil {
		t.Fatalf("SaveCSV() вернула неожиданную ошибку: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("не удалось прочитать сохранённый файл %q: %v", path, err)
	}
	got := strings.TrimSpace(string(data))
	if got != "id,description,done" {
		t.Errorf("файл пустого списка = %q, ожидается только заголовок %q", got, "id,description,done")
	}
}
