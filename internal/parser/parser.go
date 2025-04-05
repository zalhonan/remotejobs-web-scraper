package parser

import (
	"github.com/zalhonan/remotejobs-web-scraper/model"
)

type Parser interface {
	ParseJobs() (jobs []model.JobRaw, err error)
	Name() string
}
