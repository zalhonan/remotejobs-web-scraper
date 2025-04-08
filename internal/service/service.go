package service

import (
	"context"

	"github.com/zalhonan/remotejobs-web-scraper/internal/parser"
	"github.com/zalhonan/remotejobs-web-scraper/internal/repository"
	"go.uber.org/zap"
)

type service struct {
	repository repository.JobsRepository
	parsers    []parser.Parser
	logger     *zap.Logger
	context    context.Context
}

func NewService(
	repository repository.JobsRepository,
	parsers []parser.Parser,
	logger *zap.Logger,
	ctx context.Context,
) *service {
	return &service{
		repository: repository,
		parsers:    parsers,
		logger:     logger,
		context:    ctx,
	}
}
