package jobs

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/zalhonan/remotejobs-web-scraper/model"
	"go.uber.org/zap"
)

var channelTagRegexp = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

// detectMainTechnology определяет основную технологию вакансии на основе ключевых слов
func (r *repository) detectMainTechnology(content string, technologies []model.Technology) string {
	// Преобразуем контент в нижний регистр для регистронезависимого поиска
	contentLower := strings.ToLower(content)

	// Для каждой технологии проверяем наличие ключевых слов
	for _, tech := range technologies {
		for _, keyword := range tech.Keywords {
			if strings.Contains(contentLower, strings.ToLower(keyword)) {
				return tech.Technology
			}
		}
	}

	return ""
}

func (r *repository) SaveJobs(jobs []model.JobRaw) (int, error) {
	op := "repository.jobs.SaveJobs"

	if len(jobs) == 0 {
		return 0, nil
	}

	// Получаем список технологий, отсортированный по приоритету
	technologies, err := r.GetTechnologies()
	if err != nil {
		r.logger.Warn("Не удалось получить список технологий, вакансии будут сохранены без определения технологии",
			zap.Error(err))
	}

	// Создаем билдер запросов с соответствующим форматом плейсхолдеров
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// Получаем каналы и их last_post_id из БД
	channelsQuery, channelsArgs, err := psql.
		Select("tag", "last_post_id").
		From("telegram_channels").
		ToSql()

	if err != nil {
		return 0, fmt.Errorf("%s: формирование запроса каналов: %w", op, err)
	}

	rows, err := r.db.Query(r.context, channelsQuery, channelsArgs...)
	if err != nil {
		return 0, fmt.Errorf("%s: выполнение запроса каналов: %w", op, err)
	}
	defer rows.Close()

	// Создаем мапу для хранения информации о каналах
	channels := make(map[string]int64)
	for rows.Next() {
		var tag string
		var lastPostID *int64

		if err := rows.Scan(&tag, &lastPostID); err != nil {
			return 0, fmt.Errorf("%s: сканирование строки каналов: %w", op, err)
		}

		// Если lastPostID == nil, устанавливаем его в 0
		var lastID int64 = 0
		if lastPostID != nil {
			lastID = *lastPostID
		}

		channels[tag] = lastID
	}

	if err := rows.Err(); err != nil {
		return 0, fmt.Errorf("%s: итерация по результатам каналов: %w", op, err)
	}

	// Группируем вакансии по каналам
	jobsByChannel := make(map[string][]model.JobRaw)
	for _, job := range jobs {
		// Определяем тег канала из ссылки на сообщение
		// Пример: https://t.me/java_c_net_golang_jobs/1234
		parts := strings.Split(job.SourceLink, "/")
		if len(parts) < 2 {
			r.logger.Warn("Некорректная ссылка на сообщение", zap.String("link", job.SourceLink))
			continue
		}

		// Получаем тег канала и ID сообщения
		tagParts := strings.Split(parts[len(parts)-2], "@")
		tag := tagParts[len(tagParts)-1]

		// Валидация тега канала
		if !channelTagRegexp.MatchString(tag) {
			r.logger.Warn("Некорректный тег канала (валидация)", zap.String("tag", tag))
			continue
		}

		// Пропускаем, если канал не найден в БД
		if _, exists := channels[tag]; !exists {
			r.logger.Warn("Канал не найден в БД", zap.String("tag", tag))
			continue
		}

		// Определяем основную технологию вакансии
		if len(technologies) > 0 {
			job.MainTechnology = r.detectMainTechnology(job.Content, technologies)
		}

		jobsByChannel[tag] = append(jobsByChannel[tag], job)
	}

	// Для каждого канала сохраняем вакансии и обновляем информацию о канале
	totalSaved := 0

	for tag, channelJobs := range jobsByChannel {
		// Получаем текущий last_post_id канала
		lastPostID := channels[tag]
		newLastPostID := lastPostID
		newJobsCount := 0

		// Начинаем транзакцию
		tx, err := r.db.Begin(r.context)
		if err != nil {
			return totalSaved, fmt.Errorf("%s: начало транзакции: %w", op, err)
		}

		// Для всех вакансий из канала
		for _, job := range channelJobs {
			// Извлекаем ID поста из ссылки
			postIDStr := strings.Split(job.SourceLink, "/")[len(strings.Split(job.SourceLink, "/"))-1]
			postID, err := strconv.ParseInt(postIDStr, 10, 64)

			if err != nil {
				r.logger.Warn("Не удалось получить ID поста из ссылки",
					zap.String("link", job.SourceLink),
					zap.Error(err))
				continue
			}

			// Пропускаем посты, которые уже были обработаны
			if postID <= lastPostID {
				continue
			}

			// Обновляем наибольший ID поста
			if postID > newLastPostID {
				newLastPostID = postID
			}

			// Формируем INSERT запрос для вакансии
			insertQuery, insertArgs, err := psql.
				Insert("jobs_raw").
				Columns("content", "source_link", "main_technology", "date_posted", "date_parsed").
				Values(job.Content, job.SourceLink, job.MainTechnology, job.DatePosted, job.DateParsed).
				ToSql()

			if err != nil {
				r.logger.Warn("Ошибка формирования запроса вставки вакансии",
					zap.String("link", job.SourceLink),
					zap.Error(err))
				continue
			}

			// Выполняем INSERT запрос
			_, err = tx.Exec(r.context, insertQuery, insertArgs...)
			if err != nil {
				r.logger.Warn("Ошибка выполнения запроса вставки вакансии",
					zap.String("link", job.SourceLink),
					zap.Error(err))
				continue
			}

			newJobsCount++
		}

		// Если были добавлены новые вакансии, обновляем информацию о канале
		if newJobsCount > 0 {
			now := time.Now()

			// Формируем UPDATE запрос для канала
			updateQuery, updateArgs, err := psql.
				Update("telegram_channels").
				Set("last_post_id", newLastPostID).
				Set("posts_parsed", squirrel.Expr("posts_parsed + ?", newJobsCount)).
				Set("date_last_parsed", now).
				Where(squirrel.Eq{"tag": tag}).
				ToSql()

			if err != nil {
				tx.Rollback(r.context)
				return totalSaved, fmt.Errorf("%s: формирование запроса обновления канала: %w", op, err)
			}

			// Выполняем UPDATE запрос
			_, err = tx.Exec(r.context, updateQuery, updateArgs...)
			if err != nil {
				tx.Rollback(r.context)
				return totalSaved, fmt.Errorf("%s: выполнение запроса обновления канала: %w", op, err)
			}

			totalSaved += newJobsCount
		}

		// Фиксируем транзакцию
		if err := tx.Commit(r.context); err != nil {
			return totalSaved, fmt.Errorf("%s: завершение транзакции: %w", op, err)
		}
	}

	return totalSaved, nil
}
