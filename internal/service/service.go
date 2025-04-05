package service

import (
	"context"

	"github.com/zalhonan/remotejobs-web-scraper/internal/parser"
	"github.com/zalhonan/remotejobs-web-scraper/internal/repository"
)

type service struct {
	repository repository.JobsRepository
	parsers    []parser.Parser
	context    context.Context
}

func NewService(
	repository repository.JobsRepository,
	parsers []parser.Parser,
	ctx context.Context,
) *service {
	return &service{
		repository: repository,
		parsers:    parsers,
		context:    ctx,
	}
}
