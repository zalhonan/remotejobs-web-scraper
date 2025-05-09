package utils

import "unicode/utf8"

// ensureValidUTF8 проверяет и очищает строку, чтобы она содержала только валидные UTF-8 символы
func EnsureValidUTF8(s string) string {
	if utf8.ValidString(s) {
		return s
	}

	// Заменяем некорректные символы на пробелы
	result := make([]rune, 0, len(s))
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == utf8.RuneError && size == 1 {
			// Некорректный символ - заменяем на пробел
			result = append(result, ' ')
			i++
		} else {
			result = append(result, r)
			i += size
		}
	}
	return string(result)
}
