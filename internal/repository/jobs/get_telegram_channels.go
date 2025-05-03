package jobs

import (
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/zalhonan/remotejobs-web-scraper/model"
)

func (r *repository) GetTelegramChannels() ([]model.TelegramChannel, error) {
	op := "repository.jobs.GetTelegramChannels"

	// Создаем билдер запросов SQL с указанием формата плейсхолдеров для PostgreSQL
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// Формируем SELECT запрос
	sql, args, err := psql.
		Select("id", "tag", "last_post_id", "date_channel_added", "posts_parsed", "date_last_parsed").
		From("telegram_channels").
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

	channels := make([]model.TelegramChannel, 0)

	for rows.Next() {
		var channel model.TelegramChannel

		err := rows.Scan(
			&channel.ID,
			&channel.Tag,
			&channel.LastPostID,
			&channel.DateChannelAdded,
			&channel.PostsParsed,
			&channel.DateLastParsed,
		)

		if err != nil {
			return nil, fmt.Errorf("%s: сканирование строки: %w", op, err)
		}

		channels = append(channels, channel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: итерация по результатам: %w", op, err)
	}

	if len(channels) == 0 {
		r.logger.Info("Не найдено ни одного Telegram-канала в базе данных")
	}

	return channels, nil
}
