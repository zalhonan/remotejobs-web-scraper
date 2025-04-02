package telegram

import (
	"strconv"
	"time"

	"github.com/zalhonan/remotejobs-web-scraper/model"
)

func (p *telegramParser) ParseJobs() (jobs []model.JobRaw, err error) {
	ans := []model.JobRaw{}
	for i := 0; i < 10; i++ {
		ans = append(ans, model.JobRaw{
			ID:         strconv.Itoa(i),
			Content:    "Telegram job",
			SourceLink: "https://t.me/remotejobs",
			DatePosted: time.Now(),
			DateParsed: time.Now(),
		})
	}
	return ans, nil
}
