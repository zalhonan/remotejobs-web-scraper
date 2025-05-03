package jobs

import (
	"bufio"
	"os"
	"strings"

	"github.com/Masterminds/squirrel"
)

// SaveChannels загружает теги Telegram каналов из указанного файла в БД
func (r *repository) SaveChannels(filePath string) (int, error) {
	// Открываем файл для чтения
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	// Читаем файл в память
	var channelTags []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		tag := strings.TrimSpace(scanner.Text())
		if tag != "" {
			channelTags = append(channelTags, tag)
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	if len(channelTags) == 0 {
		return 0, nil
	}

	// Формируем запрос для массовой вставки с использованием squirrel
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	insertBuilder := psql.Insert("telegram_channels").
		Columns("tag", "last_post_id")

	// Добавляем все теги каналов в запрос, устанавливая last_post_id в NULL
	for _, tag := range channelTags {
		insertBuilder = insertBuilder.Values(tag, nil)
	}

	// Добавляем ON CONFLICT DO NOTHING для пропуска существующих записей
	query, args, err := insertBuilder.
		Suffix("ON CONFLICT (tag) DO NOTHING").
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return 0, err
	}

	// Выполняем запрос и получаем количество вставленных записей
	rows, err := r.db.Query(r.context, query, args...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	// Считаем количество вставленных записей
	count := 0
	for rows.Next() {
		count++
	}

	if err := rows.Err(); err != nil {
		return 0, err
	}

	return count, nil
}
