package jobs

import (
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/zalhonan/remotejobs-web-scraper/model"
)

// GetTechnologies возвращает список технологий из базы данных, отсортированный по приоритету (sort_order)
func (r *repository) GetTechnologies() ([]model.Technology, error) {
	op := "repository.jobs.GetTechnologies"

	// Создаем билдер запросов SQL с указанием формата плейсхолдеров для PostgreSQL
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// Формируем SELECT запрос с сортировкой по sort_order
	sql, args, err := psql.
		Select("id", "technology", "keywords", "sort_order").
		From("technologies").
		OrderBy("sort_order ASC").
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

	technologies := make([]model.Technology, 0)

	for rows.Next() {
		var tech model.Technology
		var keywordsArray []string

		err := rows.Scan(
			&tech.ID,
			&tech.Technology,
			&keywordsArray,
			&tech.SortOrder,
		)

		if err != nil {
			return nil, fmt.Errorf("%s: сканирование строки: %w", op, err)
		}

		tech.Keywords = keywordsArray
		technologies = append(technologies, tech)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: итерация по результатам: %w", op, err)
	}

	return technologies, nil
}
