package repository

import (
	"github.com/zalhonan/remotejobs-web-scraper/model"
)

type JobsRepository interface {
	GetTelegramChannels() ([]model.TelegramChannel, error)
	SaveJobs(jobs []model.JobRaw) (int, error)
	SaveChannels(jobsList string) (int, error)
	SaveTechnologies(technologiesFile string) (int, error)
	SaveStopWords(stopWordsFile string) (int, error)
	GetTechnologies() ([]model.Technology, error)
	GetStopWords() ([]model.StopWord, error)
	UpdateTechnologiesCount() error
}
