package test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func RegisterAndLoginUser(router http.Handler) string {
	// 1) Register user
	requestBody := strings.NewReader(`{
		"name": "Koseki Bijou",
		"email": "biboo@gmail.com",
		"password": "password",
		"passwordConfirmation": "password"
	}`)
	request := httptest.NewRequest(http.MethodPost, "/api/users", requestBody)
	SetContentTypeJson(request)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// 2) Login user
	requestBody = strings.NewReader(`{
		"email": "biboo@gmail.com",
		"password": "password"
	}`)
	request = httptest.NewRequest(http.MethodPost, "/api/users/login", requestBody)
	SetContentTypeJson(request)
	// Set user agent
	request.Header.Set("User-Agent", "Mozilla/5.0")
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	// Get the access token
	var responseBody ResponseBody
	ReadResponseBody(recorder.Result(), &responseBody)
	accessToken := responseBody["data"].(map[string]interface{})["accessToken"].(string)

	return accessToken
}

func TestShortenUrl(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter(db)
	defer db.Close()

	TruncateTable(db, "users")

	accessToken := RegisterAndLoginUser(router)

	t.Run("Shortening success", func(t *testing.T) {
		TruncateTable(db, "urls")

		// Shorten URL
		requestBody := strings.NewReader(`{
			"url": "https://azbagas.com"
		}`)
		request := httptest.NewRequest(http.MethodPost, "/api/shorten", requestBody)
		SetContentTypeJson(request)
		request.Header.Set("Authorization", "Bearer "+accessToken)
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, request)

		response := recorder.Result()
		assert.Equal(t, 201, response.StatusCode)

		var responseBody ResponseBody
		ReadResponseBody(response, &responseBody)
		data := responseBody["data"].(map[string]interface{})
		assert.NotNil(t, data["id"])
		assert.Equal(t, "https://azbagas.com", data["url"])
		assert.NotNil(t, data["shortCode"])
		assert.NotNil(t, data["createdAt"])
		assert.NotNil(t, data["updatedAt"])
	})

	t.Run("Shortening failed: Validation error", func(t *testing.T) {
		TruncateTable(db, "urls")

		// Shorten URL
		requestBody := strings.NewReader(`{
			"url": "invalid url"
		}`)
		request := httptest.NewRequest(http.MethodPost, "/api/shorten", requestBody)
		SetContentTypeJson(request)
		request.Header.Set("Authorization", "Bearer "+accessToken)
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, request)

		response := recorder.Result()
		assert.Equal(t, 400, response.StatusCode)

		var responseBody ResponseBody
		ReadResponseBody(response, &responseBody)
		assert.NotNil(t, responseBody["message"])
	})
}
