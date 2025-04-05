package telegram

import "context"

type telegramParser struct {
	ctx context.Context
}

func NewTelegramParser(context context.Context) *telegramParser {
	return &telegramParser{
		ctx: context,
	}
}
