package telegram

import (
	"fmt"
	"time"

	"github.com/gocolly/colly"
	"github.com/zalhonan/remotejobs-web-scraper/model"
)

func (p *telegramParser) parseChannel(tag string) (jobs []model.JobRaw, err error) {
	counter := 0

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"),
	)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request failed with status code:", r.StatusCode)
		fmt.Println("Error:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Page visited:", r.Request.URL)
	})

	// Обрабатываем каждый пост целиком
	c.OnHTML("div.tgme_widget_message", func(e *colly.HTMLElement) {

		messageTextDiv := e.DOM.Find("div.tgme_widget_message_text.js-message_text")
		htmlContent, _ := messageTextDiv.Html()

		// Получаем информацию из блока под сообщением
		infoBlock := e.DOM.Find("div.tgme_widget_message_info.short.js-message_info")

		// Извлекаем полную дату/время публикации из атрибута datetime тега time
		dateTime := ""
		timeElement := infoBlock.Find("a.tgme_widget_message_date time")
		if timeElement.Length() > 0 {
			dateTime, _ = timeElement.Attr("datetime")
		} else {
			// Если элемент time не найден, используем текст из meta как запасной вариант
			dateTime = infoBlock.Find("span.tgme_widget_message_meta").Text()
		}

		// Извлекаем ссылку на сообщение
		messageLink, _ := infoBlock.Find("a.tgme_widget_message_date").Attr("href")

		fmt.Printf("Ссылка на сообщение: %s\n", messageLink)

		counter++
		fmt.Printf("-----------Обработано сообщений: %d\n", counter)

		// Parse the dateTime string into a time.Time value
		parsedTime, err := time.Parse(time.RFC3339, dateTime)
		if err != nil {
			// If parsing fails, use current time as fallback
			parsedTime = time.Now()
			fmt.Printf("Error parsing date: %v\n", err)
		}

		jobs = append(jobs, model.JobRaw{
			Content:    htmlContent,
			SourceLink: messageLink,
			DatePosted: parsedTime,
			DateParsed: time.Now(),
		})
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished scraping", r.Request.URL)
	})

	channel := fmt.Sprintf("https://t.me/s/%s", tag)

	error := c.Visit(channel)

	if error != nil {
		fmt.Printf("Error visiting channel: %v\n", error)
		return nil, error
	}

	return jobs, nil

}
