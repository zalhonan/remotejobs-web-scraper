package jobs

import (
	"testing"

	"github.com/zalhonan/remotejobs-web-scraper/model"
)

func TestDetectMainTechnology(t *testing.T) {
	repo := &repository{}

	technologies := []model.Technology{
		{
			Technology: "Go",
			Keywords:   []string{"golang", "go разработчик"},
		},
		{
			Technology: "Python",
			Keywords:   []string{"python", "джуниор python"},
		},
	}

	stopWords := []model.StopWord{
		{
			Word: "стремитесь",
		},
		{
			Word: "адвокат",
		},
	}

	tests := []struct {
		name        string
		content     string
		expected    string
		description string
	}{
		{
			name:        "Technology found",
			content:     "Ищем опытного Go разработчика",
			expected:    "Go",
			description: "Должен определить правильную технологию",
		},
		{
			name:        "No technology match",
			content:     "Ищем руководителя проекта",
			expected:    "",
			description: "Должен вернуть пустую строку, если технология не найдена",
		},
		{
			name:        "Stop word found",
			content:     "Стремитесь к новым высотам с Go разработкой",
			expected:    "",
			description: "Должен вернуть пустую строку, если найдено стоп-слово",
		},
		{
			name:        "Stop word with uppercase",
			content:     "АДВОКАТ для Go разработчика",
			expected:    "",
			description: "Должен игнорировать регистр при поиске стоп-слов",
		},
		{
			name:        "Multiple technologies match but stop word",
			content:     "Нужен Go и Python разработчик, стремитесь получить работу",
			expected:    "",
			description: "Должен вернуть пустую строку, если есть стоп-слово, даже если найдены технологии",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := repo.detectMainTechnology(tc.content, technologies, stopWords)
			if result != tc.expected {
				t.Errorf("%s: ожидалось '%s', получено '%s'", tc.description, tc.expected, result)
			}
		})
	}
}
