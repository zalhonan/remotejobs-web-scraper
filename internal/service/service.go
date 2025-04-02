package service

import (
	"context"

	"github.com/zalhonan/remotejobs-web-scraper/internal/repository"
)

type service struct {
	repository repository.JobsRepository
	context    context.Context
}

func NewService(
	repository repository.JobsRepository,
	ctx context.Context,
) *service {
	return &service{
		repository: repository,
		context:    ctx,
	}
}
