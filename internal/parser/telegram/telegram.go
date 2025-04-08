package telegram

import (
	"context"

	"github.com/zalhonan/remotejobs-web-scraper/internal/repository"
	"go.uber.org/zap"
)

type telegramParser struct {
	repository repository.JobsRepository
	logger     *zap.Logger
	ctx        context.Context
}

func NewTelegramParser(
	repository repository.JobsRepository,
	logger *zap.Logger,
	context context.Context,
) *telegramParser {
	return &telegramParser{
		repository: repository,
		logger:     logger,
		ctx:        context,
	}
}
