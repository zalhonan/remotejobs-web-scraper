package telegram

import (
	"context"

	"github.com/zalhonan/remotejobs-web-scraper/internal/repository"
)

type telegramParser struct {
	repository repository.JobsRepository
	ctx        context.Context
}

func NewTelegramParser(repository repository.JobsRepository, context context.Context) *telegramParser {
	return &telegramParser{
		repository: repository,
		ctx:        context,
	}
}
