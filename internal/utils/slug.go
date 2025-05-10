package utils

import (
	"fmt"
	"regexp"
	"strings"
)

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9\s-]`)
var multipleSpacesRegex = regexp.MustCompile(`\s+`)

// Таблица транслитерации кириллических символов
var translitMap = map[rune]string{
	'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d", 'е': "e", 'ё': "yo", 'ж': "zh",
	'з': "z", 'и': "i", 'й': "j", 'к': "k", 'л': "l", 'м': "m", 'н': "n", 'о': "o",
	'п': "p", 'р': "r", 'с': "s", 'т': "t", 'у': "u", 'ф': "f", 'х': "h", 'ц': "c",
	'ч': "ch", 'ш': "sh", 'щ': "sch", 'ъ': "", 'ы': "y", 'ь': "", 'э': "e", 'ю': "yu",
	'я': "ya",
	'А': "A", 'Б': "B", 'В': "V", 'Г': "G", 'Д': "D", 'Е': "E", 'Ё': "Yo", 'Ж': "Zh",
	'З': "Z", 'И': "I", 'Й': "J", 'К': "K", 'Л': "L", 'М': "M", 'Н': "N", 'О': "O",
	'П': "P", 'Р': "R", 'С': "S", 'Т': "T", 'У': "U", 'Ф': "F", 'Х': "H", 'Ц': "C",
	'Ч': "Ch", 'Ш': "Sh", 'Щ': "Sch", 'Ъ': "", 'Ы': "Y", 'Ь': "", 'Э': "E", 'Ю': "Yu",
	'Я': "Ya",
}

// Транслитерирует кириллические символы в латинские
func transliterate(input string) string {
	var result strings.Builder
	result.Grow(len(input) * 2) // Примерная оценка размера результата

	for _, r := range input {
		if replacement, ok := translitMap[r]; ok {
			result.WriteString(replacement)
		} else {
			result.WriteRune(r)
		}
	}

	return result.String()
}

// GenerateSlug создает слаг в формате <id>-<title>
// Преобразует заголовок в нижний регистр, удаляет специальные символы и заменяет пробелы на дефисы
func GenerateSlug(id int64, title string) string {
	// Если заголовок пустой, просто возвращаем ID
	if title == "" {
		return fmt.Sprintf("%d", id)
	}

	// Транслитерируем кириллические символы в латинские
	slug := transliterate(title)

	// Преобразуем в нижний регистр
	slug = strings.ToLower(slug)

	// Удаляем все символы кроме букв, цифр, пробелов и дефисов
	slug = nonAlphanumericRegex.ReplaceAllString(slug, "")

	// Заменяем множественные пробелы одиночными
	slug = multipleSpacesRegex.ReplaceAllString(slug, " ")

	// Заменяем пробелы на дефисы
	slug = strings.ReplaceAll(slug, " ", "-")

	// Удаляем начальные и конечные дефисы
	slug = strings.Trim(slug, "-")

	// Ограничиваем длину слага (для заголовка) до 50 символов
	if len(slug) > 50 {
		slug = slug[:50]
		// Удаляем последний дефис, если он есть
		slug = strings.TrimSuffix(slug, "-")
	}

	// Формируем финальный слаг в формате <id>-<slug>
	return fmt.Sprintf("%d-%s", id, slug)
}
