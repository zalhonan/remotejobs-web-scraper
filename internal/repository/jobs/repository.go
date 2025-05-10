package jobs

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type repository struct {
	db      *pgxpool.Pool
	logger  *zap.Logger
	context context.Context
}

func NewRepository(db *pgxpool.Pool, logger *zap.Logger, ctx context.Context) *repository {
	return &repository{
		db:      db,
		logger:  logger,
		context: ctx,
	}
}

func (r *repository) UpdateTechnologiesCount() error {
	op := "repository.jobs.UpdateTechnologiesCount"

	// Формируем запрос с использованием squirrel
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// Получаем данные для обновления (технологии и их количество)
	subQuery, subArgs, err := psql.
		Select("main_technology", "COUNT(*) as cnt").
		From("jobs_raw").
		Where(squirrel.And{
			squirrel.NotEq{"main_technology": nil},
			squirrel.NotEq{"main_technology": ""},
		}).
		GroupBy("main_technology").
		ToSql()

	if err != nil {
		return fmt.Errorf("%s: формирование SQL подзапроса: %w", op, err)
	}

	// Строим запрос для обновления счётчиков в технологиях
	query := fmt.Sprintf(`
		UPDATE technologies t
		SET count = sub.cnt
		FROM (%s) AS sub
		WHERE t.technology = sub.main_technology
	`, subQuery)

	// Выполняем обновление
	_, err = r.db.Exec(r.context, query, subArgs...)
	if err != nil {
		return fmt.Errorf("%s: выполнение запроса обновления счётчиков: %w", op, err)
	}

	// Сбрасываем count у технологий, которые отсутствуют в jobs_raw
	resetQuery := `
		UPDATE technologies t
		SET count = 0
		WHERE NOT EXISTS (
			SELECT 1 
			FROM jobs_raw j 
			WHERE j.main_technology = t.technology 
			AND j.main_technology IS NOT NULL 
			AND j.main_technology != ''
		)
	`

	// Выполняем сброс
	_, err = r.db.Exec(r.context, resetQuery)
	if err != nil {
		return fmt.Errorf("%s: выполнение запроса сброса счётчиков: %w", op, err)
	}

	return nil
}
