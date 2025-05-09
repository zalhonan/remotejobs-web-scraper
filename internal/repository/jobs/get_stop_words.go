package jobs

import (
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/zalhonan/remotejobs-web-scraper/model"
)

// GetStopWords возвращает список стоп-слов из базы данных
func (r *repository) GetStopWords() ([]model.StopWord, error) {
	op := "repository.jobs.GetStopWords"

	// Создаем билдер запросов SQL с указанием формата плейсхолдеров для PostgreSQL
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// Формируем SELECT запрос
	sql, args, err := psql.
		Select("id", "word").
		From("stop_words").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("%s: формирование SQL-запроса: %w", op, err)
	}

	// Выполняем запрос
	rows, err := r.db.Query(r.context, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: выполнение запроса: %w", op, err)
	}
	defer rows.Close()

	stopWords := make([]model.StopWord, 0)

	for rows.Next() {
		var stopWord model.StopWord

		err := rows.Scan(
			&stopWord.ID,
			&stopWord.Word,
		)

		if err != nil {
			return nil, fmt.Errorf("%s: сканирование строки: %w", op, err)
		}

		stopWords = append(stopWords, stopWord)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: итерация по результатам: %w", op, err)
	}

	return stopWords, nil
}
