package jobs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zalhonan/remotejobs-web-scraper/internal/repository/test"
	"github.com/zalhonan/remotejobs-web-scraper/model"
	"go.uber.org/zap/zaptest"
)

// TestSaveJobs проверяет функцию сохранения вакансий согласно шаблону GIVEN-WHEN-THEN
func TestSaveJobs(t *testing.T) {
	// GIVEN: Создаем тестовый логгер
	logger := zaptest.NewLogger(t)

	t.Run("успешное сохранение вакансий", func(t *testing.T) {
		// GIVEN: Мокируем репозиторий и готовим тестовые данные
		mockRepo := test.NewMockRepository(logger)

		// Создаем список тестовых вакансий
		jobs := []model.JobRaw{
			test.CreateMockJob(1, "golang"),
			test.CreateMockJob(2, "java"),
			test.CreateMockJob(3, "javascript"),
		}

		// WHEN: Вызываем метод сохранения вакансий
		count, err := mockRepo.SaveJobs(jobs)

		// THEN: Проверяем результаты
		assert.NoError(t, err)
		assert.Equal(t, 3, count)
		assert.Equal(t, 3, mockRepo.SavedJobs)
	})

	t.Run("сохранение пустого списка вакансий", func(t *testing.T) {
		// GIVEN: Мокируем репозиторий и готовим пустой список
		mockRepo := test.NewMockRepository(logger)
		jobs := []model.JobRaw{}

		// WHEN: Вызываем метод сохранения пустого списка вакансий
		count, err := mockRepo.SaveJobs(jobs)

		// THEN: Проверяем, что метод корректно обрабатывает пустой список
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
		assert.Equal(t, 0, mockRepo.SavedJobs)
	})

	t.Run("ошибка при сохранении вакансий", func(t *testing.T) {
		// GIVEN: Мокируем репозиторий с ошибкой
		mockRepo := test.NewMockRepository(logger)
		mockRepo.ShouldError = true

		jobs := []model.JobRaw{
			test.CreateMockJob(1, "golang"),
		}

		// WHEN: Вызываем метод сохранения вакансий
		count, err := mockRepo.SaveJobs(jobs)

		// THEN: Проверяем, что возникла ошибка
		assert.Error(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("detectMainTechnology корректно определяет технологию", func(t *testing.T) {
		// GIVEN: Готовим тестовые данные для определения технологии
		mockRepo := test.NewMockRepository(logger)
		mockRepo.Technologies = []model.Technology{
			test.CreateMockTechnology(1, "golang", 0, "go", "golang"),
			test.CreateMockTechnology(2, "java", 1, "java", "джава"),
			test.CreateMockTechnology(3, "javascript", 2, "js", "javascript"),
		}

		// Тестируем метод напрямую через мок-объект
		job := test.CreateMockJob(1, "")
		job.Content = "Ищем опытного Go-разработчика со знанием golang"

		// WHEN: Вызываем метод получения технологий и определения технологии
		technologies, _ := mockRepo.GetTechnologies()

		// THEN: Проверяем, что технология определена правильно
		mainTechnology := mockRepo.DetectMainTechnology(job.Content, technologies)
		assert.Equal(t, "golang", mainTechnology)
	})

	t.Run("detectMainTechnology при нескольких совпадениях выбирает по приоритету", func(t *testing.T) {
		// GIVEN: Готовим тестовые данные с пересекающимися ключевыми словами
		mockRepo := test.NewMockRepository(logger)
		mockRepo.Technologies = []model.Technology{
			test.CreateMockTechnology(1, "golang", 0, "go", "golang"),
			test.CreateMockTechnology(2, "java", 1, "java", "джава"),
			test.CreateMockTechnology(3, "javascript", 2, "js", "javascript"),
		}

		// Текст содержит ключевые слова нескольких технологий
		job := test.CreateMockJob(1, "")
		job.Content = "Требуется разработчик javascript со знанием java"

		// WHEN: Вызываем метод получения технологий и определения технологии
		technologies, _ := mockRepo.GetTechnologies()

		// THEN: Проверяем, что выбрана технология с наивысшим приоритетом (наименьшим sort_order)
		mainTechnology := mockRepo.DetectMainTechnology(job.Content, technologies)
		assert.Equal(t, "java", mainTechnology)
	})

	t.Run("detectMainTechnology при отсутствии совпадений возвращает пустую строку", func(t *testing.T) {
		// GIVEN: Готовим тестовые данные без совпадений
		mockRepo := test.NewMockRepository(logger)
		mockRepo.Technologies = []model.Technology{
			test.CreateMockTechnology(1, "golang", 0, "go", "golang"),
			test.CreateMockTechnology(2, "java", 1, "java", "джава"),
		}

		job := test.CreateMockJob(1, "")
		job.Content = "Требуется разработчик C++ со знанием Python"

		// WHEN: Вызываем метод получения технологий и определения технологии
		technologies, _ := mockRepo.GetTechnologies()

		// THEN: Проверяем, что при отсутствии совпадений возвращается пустая строка
		mainTechnology := mockRepo.DetectMainTechnology(job.Content, technologies)
		assert.Equal(t, "", mainTechnology)
	})
}
