package app

import (
	"database/sql"
	"time"

	"github.com/azbagas/url-shortening-backend/config"
	"github.com/azbagas/url-shortening-backend/helper"
)

func NewDB() *sql.DB {
	db, err := sql.Open("pgx", config.AppConfig.DatabaseUrl)
	helper.PanicIfError(err)

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db
}
