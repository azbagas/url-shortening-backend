package test

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/azbagas/url-shortening-backend/app"
	"github.com/azbagas/url-shortening-backend/config"
	"github.com/azbagas/url-shortening-backend/controller"
	"github.com/azbagas/url-shortening-backend/helper"
	"github.com/azbagas/url-shortening-backend/middleware"

	"github.com/azbagas/url-shortening-backend/repository"
	"github.com/azbagas/url-shortening-backend/service"
	"github.com/go-playground/validator/v10"
)

type ResponseBody map[string]interface{}

func TestMain(m *testing.M) {
	// Load config
	os.Chdir("..")
	config.LoadConfig()

	code := m.Run()

	os.Exit(code)
}

func SetupTestDB() *sql.DB {
	db, err := sql.Open("pgx", "postgres://postgres:password@localhost:5432/url_shortening_backend")
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
	refreshTokenRepository := repository.NewRefreshTokenRepository()
	userService := service.NewUserService(userRepository, refreshTokenRepository, db, validate)
	userController := controller.NewUserController(userService)

	router := app.NewRouter(userController)
	return middleware.AuthMiddleware(router)
}

func TruncateTable(db *sql.DB, table string) {
	_, err := db.Exec("TRUNCATE TABLE " + table + " RESTART IDENTITY CASCADE")
	helper.PanicIfError(err)
}

func ReadResponseBody(response *http.Response, responseBody *ResponseBody) {
	body, _ := io.ReadAll(response.Body)
	json.Unmarshal(body, responseBody)
}

func SetContentTypeJson(request *http.Request) {
	request.Header.Set("Content-Type", "application/json")
}
