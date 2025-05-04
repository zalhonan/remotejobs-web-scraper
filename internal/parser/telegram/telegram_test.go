package telegram

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zalhonan/remotejobs-web-scraper/internal/repository/test"
	"go.uber.org/zap/zaptest"
)

// TestTelegramParser проверяет корректность инициализации парсера
func TestTelegramParser(t *testing.T) {
	// GIVEN: Создаем тестовый логгер и контекст
	logger := zaptest.NewLogger(t)
	ctx := context.Background()
	mockRepo := test.NewMockRepository(logger)

	// WHEN: Создаем новый парсер
	parser := NewTelegramParser(mockRepo, logger, ctx)

	// THEN: Проверяем, что парсер корректно инициализирован
	assert.NotNil(t, parser)
	assert.Equal(t, "Telegram", parser.Name())
}

// TestEmptyChannelsList проверяет обработку пустого списка каналов
func TestEmptyChannelsList(t *testing.T) {
	// GIVEN: Создаем тестовый логгер, контекст и пустой репозиторий
	logger := zaptest.NewLogger(t)
	ctx := context.Background()
	mockRepo := test.NewMockRepository(logger)

	// Создаем парсер
	parser := NewTelegramParser(mockRepo, logger, ctx)

	// WHEN: Вызываем метод парсинга
	jobs, err := parser.ParseJobs()

	// THEN: Проверяем, что получен пустой список без ошибок
	assert.NoError(t, err)
	assert.Empty(t, jobs)
}

// TestRepositoryError проверяет обработку ошибок репозитория
func TestRepositoryError(t *testing.T) {
	// GIVEN: Создаем тестовый логгер, контекст и репозиторий с ошибкой
	logger := zaptest.NewLogger(t)
	ctx := context.Background()
	mockRepo := test.NewMockRepository(logger)
	mockRepo.ShouldError = true

	// Создаем парсер
	parser := NewTelegramParser(mockRepo, logger, ctx)

	// WHEN: Вызываем метод парсинга
	jobs, err := parser.ParseJobs()

	// THEN: Проверяем, что возникла ошибка
	assert.Error(t, err)
	assert.Nil(t, jobs)
}
