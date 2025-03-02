package middleware

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		startTime := time.Now()
		requestId := uuid.New().String()

		logrus.WithFields(logrus.Fields{
			"requestId": requestId,
			"method":    request.Method,
			"url":       request.URL.Path,
			"clientIP":  request.RemoteAddr,
			"userAgent": request.UserAgent(),
		}).Info("Incoming request")

		next.ServeHTTP(writer, request)

		logrus.WithFields(logrus.Fields{
			"requestId": requestId,
			"latency":   time.Since(startTime).Milliseconds(),
		}).Info("Response sent")
	})
}
