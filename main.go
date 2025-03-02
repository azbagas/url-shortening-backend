package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/azbagas/url-shortening-backend/config"
	"github.com/azbagas/url-shortening-backend/helper"
	"github.com/azbagas/url-shortening-backend/middleware"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

func NewServer(router *httprouter.Router) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf("localhost:%d", config.AppConfig.AppPort),
		Handler: middleware.LogMiddleware(middleware.AuthMiddleware(router)),
	}
}

func main() {
	// Load config
	config.LoadConfig()

	// Set logger
	file, err := os.OpenFile("log/application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	helper.PanicIfError(err)
	multiWriter := io.MultiWriter(os.Stdout, file)
	logrus.SetOutput(multiWriter)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	server := InitializedServer()
	logrus.Info("Server is running on " + server.Addr)
	err = server.ListenAndServe()
	helper.PanicIfError(err)
}
