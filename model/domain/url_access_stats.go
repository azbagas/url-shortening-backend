package domain

import "time"

type UrlAccessStats struct {
	Id         int
	UrlId      int
	AccessedAt time.Time
}

type UrlAccessTotalPerDate struct {
	Date          string
	TotalAccessed int
}
