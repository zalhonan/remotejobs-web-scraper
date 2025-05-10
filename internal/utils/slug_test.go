package utils

import (
	"testing"
)

func TestTransliterate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Простая транслитерация",
			input:    "Привет мир",
			expected: "Privet mir",
		},
		{
			name:     "Смешанный текст",
			input:    "Hello Мир",
			expected: "Hello Mir",
		},
		{
			name:     "С специальными символами",
			input:    "Тест (спец-символы) и цифры 123",
			expected: "Test (spec-simvoly) i cifry 123",
		},
		{
			name:     "Длинный русский текст",
			input:    "Съешь ещё этих мягких французских булок, да выпей чаю",
			expected: "Sesh eschyo etih myagkih francuzskih bulok, da vypej chayu",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := transliterate(test.input)
			if result != test.expected {
				t.Errorf("Ожидалось: %s, получено: %s", test.expected, result)
			}
		})
	}
}

func TestGenerateSlug(t *testing.T) {
	tests := []struct {
		name     string
		id       int64
		title    string
		expected string
	}{
		{
			name:     "С обычным заголовком",
			id:       123,
			title:    "Senior Golang Developer",
			expected: "123-senior-golang-developer",
		},
		{
			name:     "С пустым заголовком",
			id:       456,
			title:    "",
			expected: "456",
		},
		{
			name:     "С заголовком, содержащим спецсимволы",
			id:       789,
			title:    "JavaScript & React.js Developer (Remote)",
			expected: "789-javascript-reactjs-developer-remote",
		},
		{
			name:     "С длинным заголовком (должен обрезаться до 50 символов)",
			id:       101112,
			title:    "Looking for a Very Skilled and Experienced Full Stack JavaScript Developer for a Long-Term Project with Many Complex Requirements and Interesting Challenges",
			expected: "101112-looking-for-a-very-skilled-and-experienced-full-st",
		},
		{
			name:     "С заголовком на русском",
			id:       131415,
			title:    "Ищем Senior Java разработчика",
			expected: "131415-ischem-senior-java-razrabotchika",
		},
		{
			name:     "С множественными пробелами",
			id:       161718,
			title:    "Python   Developer    with   Django",
			expected: "161718-python-developer-with-django",
		},
		{
			name:     "С смешанным русским и английским текстом",
			id:       192021,
			title:    "Senior Разработчик Go/Golang (удалённо)",
			expected: "192021-senior-razrabotchik-gogolang-udalyonno",
		},
		{
			name:     "С русским текстом и специальными символами",
			id:       222324,
			title:    "Требуется Разработчик на C++ (50% времени)",
			expected: "222324-trebuetsya-razrabotchik-na-c-50-vremeni",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := GenerateSlug(test.id, test.title)
			if result != test.expected {
				t.Errorf("Ожидалось: %s, получено: %s", test.expected, result)
			}
		})
	}
}
