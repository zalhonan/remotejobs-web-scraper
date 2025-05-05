package model

import "time"

type JobRaw struct {
	ID             int64
	Content        string
	Title          string
	ContentPure    string
	SourceLink     string
	MainTechnology string
	DatePosted     time.Time
	DateParsed     time.Time
}
