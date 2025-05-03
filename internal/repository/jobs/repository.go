package jobs

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type repository struct {
	db      *pgxpool.Pool
	logger  *zap.Logger
	context context.Context
}

func NewRepository(db *pgxpool.Pool, logger *zap.Logger, ctx context.Context) *repository {
	return &repository{
		db:      db,
		logger:  logger,
		context: ctx,
	}
}
