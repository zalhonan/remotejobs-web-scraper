package repository

import (
	"context"

	"github.com/zalhonan/remotejobs-web-scraper/model"
)

type JobsRepository interface {
	GetTelegramChannels(ctx context.Context) ([]model.TelegramChannel, error)
	SaveJobs(ctx context.Context, jobs []model.JobRaw) (int, error)
}
