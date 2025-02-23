package app

import (
	"database/sql"
	"time"

	"github.com/azbagas/url-shortening-backend/helper"
	"github.com/spf13/viper"
)

func NewDB(config *viper.Viper) *sql.DB {
	db, err := sql.Open("pgx", config.GetString("DATABASE_URL"))
	helper.PanicIfError(err)

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db
}
