package telegram

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gocolly/colly"
	"github.com/zalhonan/remotejobs-web-scraper/internal/utils"
	"github.com/zalhonan/remotejobs-web-scraper/model"
	"go.uber.org/zap"
)

// sanitizeUTF8 удаляет или заменяет некорректные UTF-8 символы в строке
func sanitizeUTF8(s string) string {
	if utf8.ValidString(s) {
		return s
	}

	// Заменяем некорректные символы на пробелы
	result := make([]rune, 0, len(s))
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == utf8.RuneError && size == 1 {
			// Некорректный символ - заменяем на пробел
			result = append(result, ' ')
			i++
		} else {
			result = append(result, r)
			i += size
		}
	}
	return string(result)
}

// cleanContent удаляет все HTML теги из строки, сохраняя абзацы и переводы строк
func cleanContent(html string) string {
	// Заменяем теги <br>, <p>, <div>, <h1>-<h6>, <li> на перевод строки
	re := regexp.MustCompile(`<br\s*/?>|</p>|</div>|</h[1-6]>|</li>`)
	html = re.ReplaceAllString(html, "\n")

	// Заменяем открывающие теги параграфов, заголовков и div на перевод строки
	re = regexp.MustCompile(`<p[^>]*>|<div[^>]*>|<h[1-6][^>]*>|<li[^>]*>`)
	html = re.ReplaceAllString(html, "\n")

	// Удаляем все оставшиеся HTML теги
	re = regexp.MustCompile("<[^>]*>")
	text := re.ReplaceAllString(html, "")

	// Заменяем HTML-специальные символы на обычные
	text = strings.ReplaceAll(text, "&amp;", "&")
	text = strings.ReplaceAll(text, "&lt;", "<")
	text = strings.ReplaceAll(text, "&gt;", ">")
	text = strings.ReplaceAll(text, "&quot;", "\"")
	text = strings.ReplaceAll(text, "&#39;", "'")
	text = strings.ReplaceAll(text, "&nbsp;", " ")

	// Нормализуем пробелы в каждой строке (но не удаляем переводы строк)
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		// Заменяем множественные пробелы на один
		re = regexp.MustCompile(`\s+`)
		lines[i] = strings.TrimSpace(re.ReplaceAllString(line, " "))
	}

	// Соединяем строки обратно с переводами строк
	text = strings.Join(lines, "\n")

	// Нормализуем переводы строк - не более двух подряд (одна пустая строка между параграфами)
	re = regexp.MustCompile(`\n{3,}`)
	text = re.ReplaceAllString(text, "\n\n")

	return strings.TrimSpace(text)
}

// extractFirstTagContent извлекает текст из первого HTML тега в строке
func extractFirstTagContent(html string) string {
	// Регулярное выражение для поиска первого HTML-тега с содержимым
	re := regexp.MustCompile(`<[^>]+>(.*?)</[^>]+>`)
	matches := re.FindStringSubmatch(html)

	if len(matches) >= 2 {
		// Очищаем содержимое от возможных вложенных тегов
		content := matches[1]
		return cleanContent(content)
	}

	// Если тег не найден, возвращаем пустую строку
	return ""
}

func (p *telegramParser) parseChannel(tag string) (jobs []model.JobRaw, err error) {
	op := "internal.parser.telegram.parseChannel"

	counter := 0

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"),
	)

	c.OnRequest(func(r *colly.Request) {
		p.logger.Info(
			"Visiting URL",
			zap.String("URL", r.URL.String()),
		)
	})

	c.OnError(func(r *colly.Response, err error) {
		p.logger.Warn(
			"Request failed",
			zap.Int("Status code:", r.StatusCode),
			zap.Error(err),
		)
	})

	c.OnResponse(func(r *colly.Response) {
		p.logger.Info(
			"Page visited",
			zap.String("URL", r.Request.URL.String()),
		)
	})

	// Обрабатываем каждый пост целиком
	c.OnHTML("div.tgme_widget_message", func(e *colly.HTMLElement) {

		messageTextDiv := e.DOM.Find("div.tgme_widget_message_text.js-message_text")
		htmlContent, _ := messageTextDiv.Html()
		htmlContent = sanitizeUTF8(htmlContent)

		// Получаем чистый текст без HTML-тегов используя нашу функцию
		contentPure := cleanContent(htmlContent)
		contentPure = sanitizeUTF8(contentPure)

		// Извлекаем текст из первого HTML-тега для заголовка
		title := utils.EnsureValidUTF8(extractFirstTagContent(htmlContent))

		// Если первый тег слишком длинный, то скорее всего это реклама и мы её пропускаем
		if len(title) > 70 {
			title = ""
		} else if title == "" || len(title) < 5 {
			// Если в первом теге нет текста или текст слишком короткий, используем запасной вариант
			// Используем первые 100 символов контента или весь текст, если он короче
			if len(contentPure) > 100 {
				title = strings.TrimSpace(contentPure[:100]) + "..."
			} else {
				title = contentPure
			}
		}

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

		counter++

		// Parse the dateTime string into a time.Time value
		parsedTime, err := time.Parse(time.RFC3339, dateTime)
		if err != nil {
			// If parsing fails, use current time as fallback
			parsedTime = time.Now()
		}

		jobs = append(jobs, model.JobRaw{
			Content:     htmlContent,
			Title:       title,
			ContentPure: contentPure,
			SourceLink:  messageLink,
			DatePosted:  parsedTime,
			DateParsed:  time.Now(),
		})
	})

	c.OnScraped(func(r *colly.Response) {
		p.logger.Info(
			"Page visited",
			zap.String("URL", r.Request.URL.String()),
		)
	})

	channel := fmt.Sprintf("https://t.me/s/%s", tag)

	error := c.Visit(channel)

	if error != nil {
		return nil, fmt.Errorf("%s: %w", op, error)
	}

	p.logger.Info(
		"Messages parsed",
		zap.String("URL", channel),
		zap.Int("Processed", counter),
	)

	return jobs, nil

}
