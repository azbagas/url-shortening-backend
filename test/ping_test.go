package test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/azbagas/url-shortening-backend/app"
	"github.com/azbagas/url-shortening-backend/middleware"
	"github.com/stretchr/testify/assert"
)

func setupRouter() http.Handler {
	router := app.NewRouter()
	return middleware.NewLogMiddleware(router)
}

func TestPingSuccess(t *testing.T) {
	router := setupRouter()

	request := httptest.NewRequest(http.MethodGet, "/api/ping", nil)
	request.Header.Add("Content-Type", "text/html")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	bodyString := string(body)
	assert.Equal(t, "pong", bodyString)
}
