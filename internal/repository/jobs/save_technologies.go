package jobs

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Masterminds/squirrel"
	"go.uber.org/zap"
)

// SaveTechnologies загружает технологии из CSV файла в БД
func (r *repository) SaveTechnologies(filePath string) (int, error) {
	op := "repository.jobs.SaveTechnologies"

	// Открываем CSV файл
	file, err := os.Open(filePath)
	if err != nil {
		return 0, fmt.Errorf("%s: открытие файла: %w", op, err)
	}
	defer file.Close()

	// Создаем CSV reader с гибкой настройкой
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 // Позволяет строкам иметь разное количество полей
	reader.TrimLeadingSpace = true

	// Читаем все строки из CSV
	rows, err := reader.ReadAll()
	if err != nil {
		return 0, fmt.Errorf("%s: чтение CSV: %w", op, err)
	}

	if len(rows) <= 1 {
		// Файл только с заголовками или пустой
		return 0, nil
	}

	// Удаляем строку с заголовками
	rows = rows[1:]

	// Создаем запрос для вставки с использованием squirrel
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	insertBuilder := psql.Insert("technologies").
		Columns("technology", "keywords", "sort_order")

	// Количество успешно добавленных технологий
	count := 0

	// Обрабатываем каждую строку CSV файла
	for _, row := range rows {
		// Пропускаем пустые строки или строки с недостаточным количеством полей
		if len(row) < 2 || row[0] == "" {
			continue
		}

		technology := strings.TrimSpace(row[0])

		// Парсим sort_order
		sortOrder := 0
		if len(row) > 1 && row[1] != "" {
			var err error
			sortOrder, err = strconv.Atoi(strings.TrimSpace(row[1]))
			if err != nil {
				r.logger.Warn("Некорректное значение sort_order, используем 0",
					zap.String("technology", technology),
					zap.String("sort_order", row[1]),
					zap.Error(err))
			}
		}

		// Проверяем диапазон sort_order
		if sortOrder < -100 {
			sortOrder = -100
		} else if sortOrder > 100 {
			sortOrder = 100
		}

		// Собираем ключевые слова, пропуская пустые
		var keywords []string
		for i := 2; i < len(row); i++ {
			keyword := strings.TrimSpace(row[i])
			if keyword != "" {
				keywords = append(keywords, keyword)
			}
		}

		// Если нет ключевых слов, используем название технологии
		if len(keywords) == 0 {
			keywords = append(keywords, technology)
		}

		// Формируем массив в формате PostgreSQL: {keyword1,keyword2,...}
		// Каждый элемент должен быть экранирован
		pgArray := "{"
		for i, keyword := range keywords {
			if i > 0 {
				pgArray += ","
			}
			// Экранируем каждое ключевое слово
			pgArray += strconv.Quote(keyword)
		}
		pgArray += "}"

		// Добавляем в запрос
		insertBuilder = insertBuilder.Values(
			technology,
			squirrel.Expr("?::text[]", pgArray),
			sortOrder,
		)

		// Для упрощения логирования, инкрементируем счетчик здесь
		count++
	}

	// Если нет технологий для вставки, завершаем
	if count == 0 {
		return 0, nil
	}

	// Добавляем ON CONFLICT для обновления существующих записей
	query, args, err := insertBuilder.
		Suffix("ON CONFLICT (technology) DO UPDATE SET keywords = EXCLUDED.keywords, sort_order = EXCLUDED.sort_order").
		ToSql()

	if err != nil {
		return 0, fmt.Errorf("%s: формирование SQL запроса: %w", op, err)
	}

	// Выполняем запрос
	_, err = r.db.Exec(r.context, query, args...)
	if err != nil {
		return 0, fmt.Errorf("%s: выполнение SQL запроса: %w", op, err)
	}

	// В pgx не удается получить количество затронутых строк, поэтому возвращаем посчитанное число
	return count, nil
}
