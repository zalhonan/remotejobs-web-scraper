package model

import "time"

type TelegramChannel struct {
	ID               int64
	Tag              string
	LastPostID       *int64
	DateChannelAdded time.Time
	PostsParsed      int64
	DateLastParsed   *time.Time
}
