package telegram

import "github.com/zalhonan/remotejobs-web-scraper/internal/repository"

type telegramParser struct {
	repository repository.JobsRepository
}

func NewTelegramParser(repository repository.JobsRepository) *telegramParser {
	return &telegramParser{
		repository: repository,
	}
}
