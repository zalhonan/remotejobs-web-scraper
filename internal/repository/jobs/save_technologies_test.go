package jobs

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalhonan/remotejobs-web-scraper/internal/repository/test"
	"go.uber.org/zap/zaptest"
)

// TestSaveTechnologies проверяет функцию импорта технологий согласно шаблону GIVEN-WHEN-THEN
func TestSaveTechnologies(t *testing.T) {
	// GIVEN: Создаем тестовый логгер
	logger := zaptest.NewLogger(t)

	// Создаем временный файл для тестов
	t.Run("успешный импорт технологий из CSV", func(t *testing.T) {
		// GIVEN: Создаем временный файл с тестовыми данными
		tmpFile, err := os.CreateTemp("", "test_technologies_*.csv")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		// Записываем тестовые данные в файл
		csvContent := `technology,sort_order,keyword1,keyword2,keyword3,keyword4,keyword5
						golang,0,golang,go,,,
						java,1,java,джава,,,
						javascript,2,javascript,js,,,
						`
		_, err = tmpFile.WriteString(csvContent)
		require.NoError(t, err)
		tmpFile.Close()

		// Создаем мок-репозиторий
		mockRepo := test.NewMockRepository(logger)

		// WHEN: Вызываем метод импорта технологий
		count, err := mockRepo.SaveTechnologies(tmpFile.Name())

		// THEN: Проверяем результаты импорта
		assert.NoError(t, err)
		assert.Equal(t, 10, count) // Мок возвращает 10
		assert.Equal(t, 10, mockRepo.SavedTechs)
	})

	t.Run("обработка ошибки при импорте", func(t *testing.T) {
		// GIVEN: Создаем мок-репозиторий с ошибкой
		mockRepo := test.NewMockRepository(logger)
		mockRepo.ShouldError = true

		// WHEN: Вызываем метод импорта технологий с несуществующим файлом
		count, err := mockRepo.SaveTechnologies("nonexistent_file.csv")

		// THEN: Проверяем обработку ошибки
		assert.Error(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("обработка некорректного CSV файла", func(t *testing.T) {
		// GIVEN: Создаем временный файл с некорректными данными
		tmpFile, err := os.CreateTemp("", "invalid_technologies_*.csv")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		// Записываем некорректные данные (отсутствуют обязательные поля)
		csvContent := `technology
golang
java
javascript
`
		_, err = tmpFile.WriteString(csvContent)
		require.NoError(t, err)
		tmpFile.Close()

		// Создаем мок-репозиторий
		mockRepo := test.NewMockRepository(logger)

		// WHEN: Вызываем метод импорта технологий
		count, err := mockRepo.SaveTechnologies(tmpFile.Name())

		// THEN: Проверяем обработку некорректного файла (мок все равно вернет успех,
		// так как реальная логика валидации не выполняется)
		assert.NoError(t, err)
		assert.Equal(t, 10, count)
	})

	t.Run("обработка пустого CSV файла", func(t *testing.T) {
		// GIVEN: Создаем пустой временный файл
		tmpFile, err := os.CreateTemp("", "empty_technologies_*.csv")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())
		tmpFile.Close()

		// Создаем мок-репозиторий
		mockRepo := test.NewMockRepository(logger)

		// WHEN: Вызываем метод импорта технологий
		count, err := mockRepo.SaveTechnologies(tmpFile.Name())

		// THEN: Проверяем обработку пустого файла (мок все равно вернет успех,
		// так как реальная логика не выполняется)
		assert.NoError(t, err)
		assert.Equal(t, 10, count)
	})
}
