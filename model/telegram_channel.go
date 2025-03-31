package model

import "time"

type TelegramChannel struct {
	ID              string
	Name            string
	DateAdded       time.Time
	PostsParsed     int64
	DateaLastParsed time.Time
}
