package jobs

import (
	"context"
	"time"

	"github.com/zalhonan/remotejobs-web-scraper/model"
)

func (r *repository) GetTelegramChannels(ctx context.Context) ([]model.TelegramChannel, error) {

	channelsTags := []string{
		"java_c_net_golang_jobs",
		// "java_rabota",
		// "rabota_razrabotchikh",
		// "rabota_razrabotchikq",
		// "job_javadevs",
		// "rabota_razrabotchika",
		// "rabotac_razrabotchik",
		// "rabota_razrabotchikj",
		// "jvmjobs",
		// "Java_workit",
		// "javadevjob",
	}

	channels := make([]model.TelegramChannel, 0, len(channelsTags))

	for _, tag := range channelsTags {
		channels = append(channels, model.TelegramChannel{
			Tag:            tag,
			DateLastParsed: time.Now(),
		})
	}

	return channels, nil
}
