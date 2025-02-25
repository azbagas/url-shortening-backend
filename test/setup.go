package test

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/azbagas/url-shortening-backend/app"
	"github.com/azbagas/url-shortening-backend/controller"
	"github.com/azbagas/url-shortening-backend/helper"

	// "github.com/azbagas/url-shortening-backend/middleware"
	"github.com/azbagas/url-shortening-backend/repository"
	"github.com/azbagas/url-shortening-backend/service"
	"github.com/go-playground/validator/v10"
)

type ResponseBody map[string]interface{}

func SetupTestDB() *sql.DB {
	db, err := sql.Open("pgx", "postgres://postgres:password@localhost:5432/url_shortening_backend?sslmode=disable")
	helper.PanicIfError(err)

	// Set connection pooling options
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db
}

func SetupRouter(db *sql.DB) http.Handler {
	validate := validator.New()

	userRepository := repository.NewUserRepository()
	userService := service.NewUserService(userRepository, db, validate)
	userController := controller.NewUserController(userService)

	router := app.NewRouter(userController)
	return router
	// return middleware.NewLogMiddleware(router)
}

func TruncateTable(db *sql.DB, table string) {
	_, err := db.Exec("TRUNCATE TABLE " + table + " RESTART IDENTITY")
	helper.PanicIfError(err)
}

func ReadResponseBody(response *http.Response, responseBody *ResponseBody) {
	body, _ := io.ReadAll(response.Body)
	json.Unmarshal(body, responseBody)
}

func SetContentTypeJson(request *http.Request) {
	request.Header.Set("Content-Type", "application/json")
}
