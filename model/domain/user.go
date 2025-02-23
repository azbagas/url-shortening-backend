package domain

import (
	"database/sql"
	"time"
)

type User struct {
	Id        int
	Name      string
	Email     string
	Password  string
	Photo     sql.NullString
	CreatedAt time.Time
	UpdatedAt time.Time
}
