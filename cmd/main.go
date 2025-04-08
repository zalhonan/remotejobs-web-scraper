package main

import (
	"context"

	"github.com/zalhonan/remotejobs-web-scraper/internal/logger"
	"github.com/zalhonan/remotejobs-web-scraper/internal/parser"
	"github.com/zalhonan/remotejobs-web-scraper/internal/parser/telegram"
	"github.com/zalhonan/remotejobs-web-scraper/internal/repository/jobs"
	"github.com/zalhonan/remotejobs-web-scraper/internal/service"
	"go.uber.org/zap"
)

func main() {
	logger, err := logger.NewLogger()
	if err != nil {
		panic("Cannot init logger: " + err.Error())
	}
	defer logger.Sync()

	logger.Info("Starting remote jobs scraper",
		zap.String("version", "1.0.0"),
	)

	ctx := context.Background()
	db := "db" // TODO: connect real DB here

	repository := jobs.NewRepository(db)

	telegramParser := telegram.NewTelegramParser(repository, ctx)

	parsers := []parser.Parser{telegramParser}

	service := service.NewService(repository, parsers, logger, ctx)

	if err := service.CollectJobs(); err != nil {
		logger.Error("Error collecting jobs",
			zap.Error(err),
		)
	}

	logger.Info("Jobs collected successfully")
}
