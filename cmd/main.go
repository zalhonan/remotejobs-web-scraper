package main

import (
	"context"
	"fmt"

	"github.com/zalhonan/remotejobs-web-scraper/internal/parser"
	"github.com/zalhonan/remotejobs-web-scraper/internal/parser/telegram"
	"github.com/zalhonan/remotejobs-web-scraper/internal/repository/jobs"
	"github.com/zalhonan/remotejobs-web-scraper/internal/service"
)

func main() {
	fmt.Println("Starting web scraping...")

	ctx := context.Background()
	db := "db" // TODO: connect real DB here

	repository := jobs.NewRepository(db)

	telegramParser := telegram.NewTelegramParser(ctx)

	parsers := []parser.Parser{telegramParser}

	service := service.NewService(repository, parsers, ctx)

	if err := service.CollectJobs(); err != nil {
		fmt.Printf("Error collecting jobs: %v\n", err)
		return
	}
}
