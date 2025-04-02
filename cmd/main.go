package main

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
	"golang.org/x/net/html"
)

func main() {
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

		processedText := strings.ReplaceAll(htmlContent, "<br>", "\n")
		processedText = strings.ReplaceAll(processedText, "<br/>", "\n")
		processedText = stripHTMLTags(processedText)

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

		fmt.Printf("\n==== СООБЩЕНИЕ ====\n")
		// fmt.Printf("Текст сообщения:\n%s\n", processedText) // process text
		fmt.Printf("Текст сообщения:\n%s\n", htmlContent)
		fmt.Printf("-------------------\n")
		fmt.Printf("Дата и время публикации: %s\n", dateTime)
		fmt.Printf("Ссылка на сообщение: %s\n", messageLink)
		fmt.Printf("==== КОНЕЦ СООБЩЕНИЯ ====\n\n")

		counter++
		fmt.Printf("-----------Обработано сообщений: %d\n", counter)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished scraping", r.Request.URL)
	})

	c.Visit("https://t.me/s/rabota_razrabotchikj")
	// c.Visit("https://t.me/s/javadevjob")
}

func stripHTMLTags(s string) string {
	doc, err := html.Parse(strings.NewReader(s))
	if err != nil {
		return s
	}

	var textBuilder strings.Builder
	var extractText func(*html.Node)
	extractText = func(n *html.Node) {
		if n.Type == html.TextNode {
			textBuilder.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractText(c)
		}
	}
	extractText(doc)

	return textBuilder.String()
}
