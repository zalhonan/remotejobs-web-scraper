package telegram

import (
	"fmt"

	"github.com/zalhonan/remotejobs-web-scraper/model"
)

func (p *telegramParser) ParseJobs() (jobs []model.JobRaw, err error) {
	channels, error := p.repository.GetTelegramChannels(p.ctx)

	if error != nil {
		fmt.Printf("error getting telegram channels: %v", error)
		return nil, error
	}

	for _, channel := range channels {
		parsedJobs, error := p.parseChannel(channel.Tag)

		if error != nil {
			fmt.Printf("error parsing channel: %v", error)
		}

		jobs = append(jobs, parsedJobs...)

	}

	return jobs, nil
}
