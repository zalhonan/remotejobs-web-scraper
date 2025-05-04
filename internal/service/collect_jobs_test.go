package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zalhonan/remotejobs-web-scraper/internal/parser"
	"github.com/zalhonan/remotejobs-web-scraper/internal/repository/test"
	"github.com/zalhonan/remotejobs-web-scraper/model"
	"go.uber.org/zap/zaptest"
)

// TestCollectJobs проверяет функцию сбора вакансий согласно шаблону GIVEN-WHEN-THEN
func TestCollectJobs(t *testing.T) {
	// GIVEN: Создаем тестовый логгер и контекст
	logger := zaptest.NewLogger(t)
	ctx := context.Background()

	t.Run("успешный сбор вакансий из одного парсера", func(t *testing.T) {
		// GIVEN: Мокируем репозиторий и парсеры
		mockRepo := test.NewMockRepository(logger)

		// Создаем мок-парсер с тестовыми вакансиями
		mockParser := test.NewMockParser(logger)
		mockParser.Jobs = []model.JobRaw{
			test.CreateMockJob(1, "golang"),
			test.CreateMockJob(2, "java"),
		}

		// Создаем тестовый сервис
		service := NewService(mockRepo, []parser.Parser{mockParser}, logger, ctx)

		// WHEN: Вызываем метод сбора вакансий
		err := service.CollectJobs()

		// THEN: Проверяем результаты
		assert.NoError(t, err)
		assert.Equal(t, 2, mockRepo.SavedJobs) // Должно быть сохранено 2 вакансии
	})

	t.Run("успешный сбор вакансий из нескольких парсеров", func(t *testing.T) {
		// GIVEN: Мокируем репозиторий и несколько парсеров
		mockRepo := test.NewMockRepository(logger)

		// Создаем первый мок-парсер с вакансиями
		mockParser1 := test.NewMockParser(logger)
		mockParser1.Jobs = []model.JobRaw{
			test.CreateMockJob(1, "golang"),
			test.CreateMockJob(2, "java"),
		}
		mockParser1.ParserName = "Parser1"

		// Создаем второй мок-парсер с вакансиями
		mockParser2 := test.NewMockParser(logger)
		mockParser2.Jobs = []model.JobRaw{
			test.CreateMockJob(3, "javascript"),
		}
		mockParser2.ParserName = "Parser2"

		// Создаем тестовый сервис с двумя парсерами
		service := NewService(mockRepo, []parser.Parser{mockParser1, mockParser2}, logger, ctx)

		// WHEN: Вызываем метод сбора вакансий
		err := service.CollectJobs()

		// THEN: Проверяем результаты
		assert.NoError(t, err)
		// Сначала сохраняется 2 вакансии, затем 1. mockRepo.SavedJobs содержит только последнее значение.
		assert.Equal(t, 1, mockRepo.SavedJobs)
	})

	t.Run("обработка ошибки парсера", func(t *testing.T) {
		// GIVEN: Мокируем репозиторий и парсер с ошибкой
		mockRepo := test.NewMockRepository(logger)

		// Создаем мок-парсер, который вернет ошибку
		mockParser := test.NewMockParser(logger)
		mockParser.ShouldError = true

		// Создаем тестовый сервис
		service := NewService(mockRepo, []parser.Parser{mockParser}, logger, ctx)

		// WHEN: Вызываем метод сбора вакансий
		err := service.CollectJobs()

		// THEN: Метод должен продолжить выполнение, несмотря на ошибку парсера
		assert.NoError(t, err)
		assert.Equal(t, 0, mockRepo.SavedJobs) // Ничего не должно быть сохранено
	})

	t.Run("обработка ошибки сохранения", func(t *testing.T) {
		// GIVEN: Мокируем репозиторий с ошибкой сохранения и парсер
		mockRepo := test.NewMockRepository(logger)
		mockRepo.ShouldError = true

		// Создаем мок-парсер с вакансиями
		mockParser := test.NewMockParser(logger)
		mockParser.Jobs = []model.JobRaw{
			test.CreateMockJob(1, "golang"),
		}

		// Создаем тестовый сервис
		service := NewService(mockRepo, []parser.Parser{mockParser}, logger, ctx)

		// WHEN: Вызываем метод сбора вакансий
		err := service.CollectJobs()

		// THEN: Метод должен продолжить выполнение, несмотря на ошибку сохранения
		assert.NoError(t, err)
		assert.Equal(t, 0, mockRepo.SavedJobs) // Ничего не должно быть сохранено из-за ошибки
	})

	t.Run("обработка пустого списка вакансий", func(t *testing.T) {
		// GIVEN: Мокируем репозиторий и парсер с пустым списком вакансий
		mockRepo := test.NewMockRepository(logger)

		// Создаем мок-парсер с пустым списком вакансий
		mockParser := test.NewMockParser(logger)
		// Jobs по умолчанию инициализируется как пустой слайс

		// Создаем тестовый сервис
		service := NewService(mockRepo, []parser.Parser{mockParser}, logger, ctx)

		// WHEN: Вызываем метод сбора вакансий
		err := service.CollectJobs()

		// THEN: Метод должен корректно обработать пустой список
		assert.NoError(t, err)
		assert.Equal(t, 0, mockRepo.SavedJobs) // Ничего не должно быть сохранено
	})
}
