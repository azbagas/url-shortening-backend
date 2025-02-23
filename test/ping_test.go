package test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPingSuccess(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter(db)

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
