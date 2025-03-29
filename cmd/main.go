package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

func main() {
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

	c.OnHTML("div.tgme_widget_message_text.js-message_text", func(e *colly.HTMLElement) {

		fmt.Printf("Message text: %s\n", e.Text)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished scraping", r.Request.URL)
	})

	// c.Visit("https://www.scrapingcourse.com/ecommerce")
	c.Visit("https://t.me/s/rabota_razrabotchikj")
}
