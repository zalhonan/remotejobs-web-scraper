package jobs

import (
	"context"
	"fmt"

	"github.com/zalhonan/remotejobs-web-scraper/model"
)

func (r *repository) SaveJobs(ctx context.Context, jobs []model.JobRaw) (int, error) {
	for _, job := range jobs {
		fmt.Printf("Saving job by joblink %s\n", job.SourceLink)
	}

	return len(jobs), nil
}
