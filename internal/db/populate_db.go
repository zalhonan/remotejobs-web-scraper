package db

import (
	"context"

	"github.com/zalhonan/remotejobs-web-scraper/internal/repository"
	"go.uber.org/zap"
)

// PopulateDatabase загружает все необходимые данные в базу данных
func PopulateDatabase(ctx context.Context, repo repository.JobsRepository, logger *zap.Logger) error {
	// Импорт Telegram каналов
	telegramChannelsPath := "../data/telegram_channels.txt"
	if err := populateTelegramChannels(ctx, repo, telegramChannelsPath, logger); err != nil {
		return err
	}

	// Импорт технологий
	technologiesPath := "../data/technologies.csv"
	if err := populateTechnologies(ctx, repo, technologiesPath, logger); err != nil {
		return err
	}

	return nil
}

// populateTelegramChannels импортирует Telegram каналы из файла в базу данных
func populateTelegramChannels(ctx context.Context, repo repository.JobsRepository, filePath string, logger *zap.Logger) error {
	logger.Info("Начинаем импорт Telegram каналов из файла", zap.String("filePath", filePath))

	channels, err := repo.SaveChannels(filePath)
	if err != nil {
		logger.Error("Ошибка сохранения каналов", zap.Error(err))
		return err
	}

	logger.Info("Каналы успешно сохранены", zap.Int("count", channels))
	return nil
}

// populateTechnologies импортирует технологии из CSV файла в базу данных
func populateTechnologies(ctx context.Context, repo repository.JobsRepository, filePath string, logger *zap.Logger) error {
	logger.Info("Начинаем импорт технологий из файла", zap.String("filePath", filePath))

	technologies, err := repo.SaveTechnologies(filePath)
	if err != nil {
		logger.Error("Ошибка сохранения технологий", zap.Error(err))
		return err
	}

	logger.Info("Технологии успешно сохранены", zap.Int("count", technologies))
	return nil
}
