package model

import "time"

type TelegramChannel struct {
	ID               string
	Name             string
	DateChannelAdded time.Time
	PostsParsed      int64
	DateLastParsed   time.Time
	LastPostID       int64
}
