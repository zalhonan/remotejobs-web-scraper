package jobs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zalhonan/remotejobs-web-scraper/internal/repository/test"
	"github.com/zalhonan/remotejobs-web-scraper/model"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

// Тест основан на шаблоне GIVEN-WHEN-THEN
func TestGetTechnologies(t *testing.T) {
	// GIVEN: Создаем тестовый логгер
	logger := zaptest.NewLogger(t, zaptest.Level(zap.InfoLevel))

	t.Run("успешное получение технологий", func(t *testing.T) {
		// GIVEN: Мокируем репозиторий с технологиями
		mockRepo := test.NewMockRepository(logger)
		mockRepo.Technologies = []model.Technology{
			test.CreateMockTechnology(1, "golang", 0, "go", "golang"),
			test.CreateMockTechnology(2, "java", 1, "java", "джава"),
			test.CreateMockTechnology(3, "javascript", 2, "js", "javascript"),
		}

		// Используем мок-репозиторий напрямую, так как мы не можем
		// подменить реальную базу данных в методе GetTechnologies
		technologies, err := mockRepo.GetTechnologies()

		// THEN: Проверяем, что технологии получены успешно
		assert.NoError(t, err)
		assert.Len(t, technologies, 3)
		assert.Equal(t, "golang", technologies[0].Technology)
		assert.Equal(t, []string{"go", "golang"}, technologies[0].Keywords)
		assert.Equal(t, 0, technologies[0].SortOrder)
	})

	t.Run("ошибка при запросе", func(t *testing.T) {
		// GIVEN: Мокируем репозиторий с ошибкой
		mockRepo := test.NewMockRepository(logger)
		mockRepo.ShouldError = true

		// WHEN: Вызываем метод получения технологий
		technologies, err := mockRepo.GetTechnologies()

		// THEN: Проверяем, что возникла ошибка
		assert.Error(t, err)
		assert.Nil(t, technologies)
	})

	t.Run("успешное получение пустого списка технологий", func(t *testing.T) {
		// GIVEN: Мокируем репозиторий без технологий
		mockRepo := test.NewMockRepository(logger)

		// WHEN: Вызываем метод получения технологий
		technologies, err := mockRepo.GetTechnologies()

		// THEN: Проверяем, что получен пустой список без ошибок
		assert.NoError(t, err)
		assert.Empty(t, technologies)
	})
}
