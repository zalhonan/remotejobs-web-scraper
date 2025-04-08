package telegram

import (
	"fmt"

	"github.com/zalhonan/remotejobs-web-scraper/model"
	"go.uber.org/zap"
)

func (p *telegramParser) ParseJobs() (jobs []model.JobRaw, err error) {
	op := "internal.parser.telegram.ParseJobs"

	channels, error := p.repository.GetTelegramChannels(p.ctx)

	if error != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for _, channel := range channels {
		parsedJobs, error := p.parseChannel(channel.Tag)

		if error != nil {
			p.logger.Warn(
				"Error parsing jobs from channel",
				zap.String("Channel", channel.Tag),
				zap.Error(err),
			)
		}

		jobs = append(jobs, parsedJobs...)

	}

	return jobs, nil
}
