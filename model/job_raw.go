package model

import "time"

type JobRaw struct {
	ID         string
	Content    string
	SourceLink string
	DatePosted time.Time
	DateParsed time.Time
}
