package jobs

import (
	"fmt"

	"github.com/zalhonan/remotejobs-web-scraper/model"
)

func (r *repository) SaveJobs(jobs []model.JobRaw) (int, error) {
	for _, job := range jobs {
		fmt.Printf("Saving job by joblink %s\n", job.SourceLink)
	}

	return len(jobs), nil
}
