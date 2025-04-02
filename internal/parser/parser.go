package parser

import (
	"context"

	"github.com/zalhonan/remotejobs-web-scraper/model"
)

type Parser interface {
	ParseJobs(ctx context.Context) (jobs []model.JobRaw, err error)
}
