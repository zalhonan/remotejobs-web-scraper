package main

import (
	"context"

	"github.com/zalhonan/remotejobs-web-scraper/internal/db"
	"github.com/zalhonan/remotejobs-web-scraper/internal/logger"
	"github.com/zalhonan/remotejobs-web-scraper/internal/parser"
	"github.com/zalhonan/remotejobs-web-scraper/internal/parser/telegram"
	"github.com/zalhonan/remotejobs-web-scraper/internal/repository/jobs"
	"github.com/zalhonan/remotejobs-web-scraper/internal/service"
	"go.uber.org/zap"
)

func main() {
	logger, err := logger.InitLogger()
	if err != nil {
		panic("Cannot init logger: " + err.Error())
	}
	defer logger.Sync()

	logger.Info("Starting remote jobs scraper",
		zap.String("version", "1.0.0"),
	)

	ctx := context.Background()

	// Инициализация соединения с базой данных
	database, err := db.InitDB(ctx, logger)
	if err != nil {
		logger.Fatal("Не удалось инициализировать базу данных", zap.Error(err))
	}
	defer database.Close()

	repository := jobs.NewRepository(database, logger, ctx)

	// Загрузка необходимых данных в базу данных
	if err := db.PopulateDatabase(ctx, repository, logger); err != nil {
		logger.Error("Ошибка при загрузке данных в базу", zap.Error(err))
	}

	telegramParser := telegram.NewTelegramParser(repository, logger, ctx)

	parsers := []parser.Parser{telegramParser}

	service := service.NewService(repository, parsers, logger, ctx)

	if err := service.CollectJobs(); err != nil {
		logger.Error("Ошибка сбора вакансий",
			zap.Error(err),
		)
	}

	logger.Info("Вакансии успешно собраны")

	if err := repository.UpdateTechnologiesCount(); err != nil {
		logger.Error("Ошибка обновления count в technologies", zap.Error(err))
	} else {
		logger.Info("Таблица technologies обновлена: count пересчитан")
	}
}
