package jobs

import (
	"context"
	"strings"
	"time"

	"github.com/zalhonan/remotejobs-web-scraper/model"
)

func (r *repository) GetChannels(ctx context.Context) ([]model.TelegramChannel, error) {

	channelsTags := []string{
		"java_c_net_golang_jobs",
		// "https://t.me/java_rabota", // not subscribed
		// "https://t.me/rabota_razrabotchikh",
		// "https://t.me/rabota_razrabotchikq",
		// "https://t.me/job_javadevs",
		// "https://t.me/rabota_razrabotchika",
		// "https://t.me/rabotac_razrabotchik",
		// "https://t.me/rabota_razrabotchikj",
		// "https://t.me/jvmjobs",
		// "https://t.me/Java_workit",
		// "https://t.me/javadevjob",
	}

	channels := make([]model.TelegramChannel, 0, len(channelsTags))

	for _, tag := range channelsTags {
		channelTag := strings.ReplaceAll(tag, "https://t.me/", "")
		channels = append(channels, model.TelegramChannel{
			Name:      channelTag,
			DateAdded: time.Now(),
		})
	}

	return channels, nil
}
