package middleware

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	Handler http.Handler
}

func NewLogMiddleware(handler http.Handler) *LogMiddleware {
	return &LogMiddleware{Handler: handler}
}

func (middleware *LogMiddleware) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	startTime := time.Now()
	requestId := uuid.New().String()

	logrus.WithFields(logrus.Fields{
		"requestId": requestId,
		"method":    request.Method,
		"url":       request.URL.Path,
		"clientIP":  request.RemoteAddr,
		"userAgent": request.UserAgent(),
	}).Info("Incoming request")

	middleware.Handler.ServeHTTP(writer, request)

	logrus.WithFields(logrus.Fields{
		"requestId": requestId,
		"latency":   time.Since(startTime).Milliseconds(),
	}).Info("Response sent")
}
