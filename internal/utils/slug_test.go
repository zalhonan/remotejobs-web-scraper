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
		name           string
		id             int64
		title          string
		mainTechnology string
		expected       string
	}{
		{
			name:           "С обычным заголовком",
			id:             123,
			title:          "Senior Golang Developer",
			mainTechnology: "Golang",
			expected:       "123-senior-golang-developer",
		},
		{
			name:           "С пустым заголовком и указанной технологией",
			id:             456,
			title:          "",
			mainTechnology: "JavaScript",
			expected:       "456-javascript",
		},
		{
			name:           "С пустым заголовком и пустой технологией",
			id:             789,
			title:          "",
			mainTechnology: "",
			expected:       "789",
		},
		{
			name:           "С заголовком, содержащим спецсимволы",
			id:             789,
			title:          "JavaScript & React.js Developer (Remote)",
			mainTechnology: "React",
			expected:       "789-javascript-reactjs-developer-remote",
		},
		{
			name:           "С длинным заголовком (должен обрезаться до 50 символов)",
			id:             101112,
			title:          "Looking for a Very Skilled and Experienced Full Stack JavaScript Developer for a Long-Term Project with Many Complex Requirements and Interesting Challenges",
			mainTechnology: "JavaScript",
			expected:       "101112-looking-for-a-very-skilled-and-experienced-full-st",
		},
		{
			name:           "С заголовком на русском",
			id:             131415,
			title:          "Ищем Senior Java разработчика",
			mainTechnology: "Java",
			expected:       "131415-ischem-senior-java-razrabotchika",
		},
		{
			name:           "С множественными пробелами",
			id:             161718,
			title:          "Python   Developer    with   Django",
			mainTechnology: "Python",
			expected:       "161718-python-developer-with-django",
		},
		{
			name:           "С смешанным русским и английским текстом",
			id:             192021,
			title:          "Senior Разработчик Go/Golang (удалённо)",
			mainTechnology: "Golang",
			expected:       "192021-senior-razrabotchik-gogolang-udalyonno",
		},
		{
			name:           "С русским текстом и специальными символами",
			id:             222324,
			title:          "Требуется Разработчик на C++ (50% времени)",
			mainTechnology: "C++",
			expected:       "222324-trebuetsya-razrabotchik-na-c-50-vremeni",
		},
		{
			name:           "С заголовком 'вакансия' и указанной технологией",
			id:             252627,
			title:          "Вакансия",
			mainTechnology: "PHP",
			expected:       "252627-vakansiya-php",
		},
		{
			name:           "С заголовком 'vacancy' и указанной технологией",
			id:             282930,
			title:          "Vacancy",
			mainTechnology: "Python",
			expected:       "282930-vacancy-python",
		},
		{
			name:           "С заголовком 'вакансия' и пустой технологией",
			id:             313233,
			title:          "Вакансия",
			mainTechnology: "",
			expected:       "313233-vakansiya",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := GenerateSlug(test.id, test.title, test.mainTechnology)
			if result != test.expected {
				t.Errorf("Ожидалось: %s, получено: %s", test.expected, result)
			}
		})
	}
}
