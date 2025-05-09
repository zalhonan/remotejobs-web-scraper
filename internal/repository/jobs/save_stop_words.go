package jobs

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Masterminds/squirrel"
	"go.uber.org/zap"
)

// SaveStopWords загружает стоп-слова из указанного файла в БД
func (r *repository) SaveStopWords(filePath string) (int, error) {
	op := "repository.jobs.SaveStopWords"

	// Открываем файл для чтения
	file, err := os.Open(filePath)
	if err != nil {
		return 0, fmt.Errorf("%s: открытие файла: %w", op, err)
	}
	defer file.Close()

	// Читаем файл в память
	var stopWords []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if word != "" {
			stopWords = append(stopWords, word)
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("%s: чтение файла: %w", op, err)
	}

	if len(stopWords) == 0 {
		return 0, nil
	}

	// Формируем запрос для массовой вставки с использованием squirrel
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	insertBuilder := psql.Insert("stop_words").
		Columns("word")

	// Добавляем все стоп-слова в запрос
	for _, word := range stopWords {
		insertBuilder = insertBuilder.Values(word)
	}

	// Добавляем ON CONFLICT DO NOTHING для пропуска существующих записей
	query, args, err := insertBuilder.
		Suffix("ON CONFLICT (word) DO NOTHING").
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return 0, fmt.Errorf("%s: формирование запроса: %w", op, err)
	}

	// Выполняем запрос и получаем количество вставленных записей
	rows, err := r.db.Query(r.context, query, args...)
	if err != nil {
		return 0, fmt.Errorf("%s: выполнение запроса: %w", op, err)
	}
	defer rows.Close()

	// Считаем количество вставленных записей
	count := 0
	for rows.Next() {
		count++
	}

	if err := rows.Err(); err != nil {
		return 0, fmt.Errorf("%s: итерация по результатам: %w", op, err)
	}

	r.logger.Info("Стоп-слова успешно добавлены", zap.Int("count", count))

	return count, nil
}
