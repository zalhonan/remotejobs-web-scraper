package main

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
	"golang.org/x/net/html"
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

		htmlContent, _ := e.DOM.Html()

		processedText := strings.ReplaceAll(htmlContent, "<br>", "\n")
		processedText = strings.ReplaceAll(processedText, "<br/>", "\n")
		processedText = stripHTMLTags(processedText)

		fmt.Printf("Message text:\n%s\n", processedText)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished scraping", r.Request.URL)
	})

	c.Visit("https://t.me/s/rabota_razrabotchikj")
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
