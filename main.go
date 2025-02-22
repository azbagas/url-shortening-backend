package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/azbagas/url-shortening-backend/helper"
	"github.com/azbagas/url-shortening-backend/middleware"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewServer(config *viper.Viper, logMiddleware *middleware.LogMiddleware) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf("localhost:%d", config.GetInt("APP_PORT")),
		Handler: logMiddleware,
	}
}

func main() {
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
