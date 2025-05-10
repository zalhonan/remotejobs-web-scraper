package test

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/zalhonan/remotejobs-web-scraper/model"
	"go.uber.org/zap"
)

// MockRepository реализует интерфейс repository.JobsRepository для тестирования
type MockRepository struct {
	TelegramChannels []model.TelegramChannel
	Technologies     []model.Technology
	SavedJobs        int
	SavedChannels    int
	SavedTechs       int
	ShouldError      bool
	Logger           *zap.Logger
}

// NewMockRepository создает новый мок-репозиторий для тестирования
func NewMockRepository(logger *zap.Logger) *MockRepository {
	return &MockRepository{
		TelegramChannels: []model.TelegramChannel{},
		Technologies:     []model.Technology{},
		SavedJobs:        0,
		SavedChannels:    0,
		SavedTechs:       0,
		ShouldError:      false,
		Logger:           logger,
	}
}

// GetTelegramChannels возвращает моковые каналы Telegram
func (m *MockRepository) GetTelegramChannels() ([]model.TelegramChannel, error) {
	if m.ShouldError {
		return nil, errors.New("mock error getting telegram channels")
	}
	return m.TelegramChannels, nil
}

// SaveJobs имитирует сохранение вакансий
func (m *MockRepository) SaveJobs(jobs []model.JobRaw) (int, error) {
	if m.ShouldError {
		return 0, errors.New("mock error saving jobs")
	}
	m.SavedJobs = len(jobs)
	return m.SavedJobs, nil
}

// SaveChannels имитирует сохранение каналов
func (m *MockRepository) SaveChannels(channelsFile string) (int, error) {
	if m.ShouldError {
		return 0, errors.New("mock error saving channels")
	}
	m.SavedChannels = 5 // Имитация сохранения 5 каналов
	return m.SavedChannels, nil
}

// SaveTechnologies имитирует сохранение технологий
func (m *MockRepository) SaveTechnologies(technologiesFile string) (int, error) {
	if m.ShouldError {
		return 0, errors.New("mock error saving technologies")
	}
	m.SavedTechs = 10 // Имитация сохранения 10 технологий
	return m.SavedTechs, nil
}

// GetTechnologies возвращает моковые технологии
func (m *MockRepository) GetTechnologies() ([]model.Technology, error) {
	if m.ShouldError {
		return nil, errors.New("mock error getting technologies")
	}
	if len(m.Technologies) > 0 {
		return m.Technologies, nil
	}
	return []model.Technology{
		{ID: 1, Technology: "Go", Keywords: []string{"golang", "go lang"}, SortOrder: 10},
		{ID: 2, Technology: "Python", Keywords: []string{"python", "python dev"}, SortOrder: 20},
		{ID: 3, Technology: "Java", Keywords: []string{"java", "java programming"}, SortOrder: 30},
	}, nil
}

// GetStopWords возвращает моковые стоп-слова
func (m *MockRepository) GetStopWords() ([]model.StopWord, error) {
	if m.ShouldError {
		return nil, errors.New("mock error getting stop words")
	}
	return []model.StopWord{
		{ID: 1, Word: "стремитесь"},
		{ID: 2, Word: "адвокат"},
		{ID: 3, Word: "реклама"},
	}, nil
}

// SaveStopWords имитирует сохранение стоп-слов
func (m *MockRepository) SaveStopWords(stopWordsFile string) (int, error) {
	return 3, nil
}

// UpdateTechnologiesCount имитирует обновление счетчика вакансий для каждой технологии
func (m *MockRepository) UpdateTechnologiesCount() error {
	if m.ShouldError {
		return errors.New("mock error updating technologies count")
	}
	return nil
}

// DetectMainTechnology определяет основную технологию вакансии на основе ключевых слов
func (m *MockRepository) DetectMainTechnology(content string, technologies []model.Technology) string {
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

// MockParser реализует интерфейс parser.Parser для тестирования
type MockParser struct {
	Jobs        []model.JobRaw
	ShouldError bool
	ParserName  string
	Logger      *zap.Logger
}

// NewMockParser создает новый мок-парсер для тестирования
func NewMockParser(logger *zap.Logger) *MockParser {
	return &MockParser{
		Jobs:        []model.JobRaw{},
		ShouldError: false,
		ParserName:  "MockParser",
		Logger:      logger,
	}
}

// ParseJobs имитирует парсинг вакансий
func (p *MockParser) ParseJobs() ([]model.JobRaw, error) {
	if p.ShouldError {
		return nil, errors.New("mock error parsing jobs")
	}
	return p.Jobs, nil
}

// Name возвращает имя парсера
func (p *MockParser) Name() string {
	return p.ParserName
}

// CreateMockJob создает тестовую вакансию
func CreateMockJob(id int64, technology string) model.JobRaw {
	content := fmt.Sprintf("Тестовая вакансия %d по технологии %s. Это описание вакансии.", id, technology)
	return model.JobRaw{
		ID:             id,
		Content:        fmt.Sprintf("<p>%s</p>", content),
		Title:          fmt.Sprintf("Тестовая вакансия %d по технологии %s.", id, technology),
		ContentPure:    content,
		SourceLink:     fmt.Sprintf("https://t.me/test_channel/%d", id),
		MainTechnology: technology,
		DatePosted:     time.Now().Add(-24 * time.Hour),
		DateParsed:     time.Now(),
	}
}

// CreateMockTechnology создает тестовую технологию
func CreateMockTechnology(id int64, name string, sortOrder int, keywords ...string) model.Technology {
	return model.Technology{
		ID:         id,
		Technology: name,
		Keywords:   keywords,
		SortOrder:  sortOrder,
	}
}

// CreateMockTelegramChannel создает тестовый Telegram-канал
func CreateMockTelegramChannel(id int64, tag string, lastPostID int64, postsParsed int64) model.TelegramChannel {
	now := time.Now().Add(-24 * time.Hour)

	return model.TelegramChannel{
		ID:               id,
		Tag:              tag,
		LastPostID:       &lastPostID,
		DateChannelAdded: time.Now().Add(-30 * 24 * time.Hour),
		PostsParsed:      postsParsed,
		DateLastParsed:   &now,
	}
}
