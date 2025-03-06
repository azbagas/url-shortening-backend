package domain

import (
	"time"
)

type Url struct {
	Id        int
	UserId    int
	Url       string
	ShortCode string
	CreatedAt time.Time
	UpdatedAt time.Time
}
