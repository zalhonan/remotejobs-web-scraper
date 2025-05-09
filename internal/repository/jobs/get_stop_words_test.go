package jobs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zalhonan/remotejobs-web-scraper/internal/repository/test"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

// Тест основан на шаблоне GIVEN-WHEN-THEN
func TestGetStopWords(t *testing.T) {
	// GIVEN: Создаем тестовый логгер
	logger := zaptest.NewLogger(t, zaptest.Level(zap.InfoLevel))

	t.Run("успешное получение стоп-слов", func(t *testing.T) {
		// GIVEN: Мокируем репозиторий со стоп-словами
		mockRepo := test.NewMockRepository(logger)

		// WHEN: Вызываем метод получения стоп-слов
		stopWords, err := mockRepo.GetStopWords()

		// THEN: Проверяем, что стоп-слова получены успешно
		assert.NoError(t, err)
		assert.Len(t, stopWords, 3)
		assert.Equal(t, "стремитесь", stopWords[0].Word)
		assert.Equal(t, "адвокат", stopWords[1].Word)
		assert.Equal(t, "реклама", stopWords[2].Word)
	})

	t.Run("ошибка при запросе", func(t *testing.T) {
		// GIVEN: Мокируем репозиторий с ошибкой
		mockRepo := test.NewMockRepository(logger)
		mockRepo.ShouldError = true

		// WHEN: Вызываем метод получения стоп-слов
		stopWords, err := mockRepo.GetStopWords()

		// THEN: Проверяем, что возникла ошибка
		assert.Error(t, err)
		assert.Nil(t, stopWords)

		// Восстанавливаем флаг для следующих тестов
		mockRepo.ShouldError = false
	})
}
