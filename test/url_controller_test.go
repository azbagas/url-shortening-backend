package test

import (
	"fmt"
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

func createShortenedUrl(router http.Handler, accessToken string, url string) {
	// Shorten URL
	requestBody := strings.NewReader(fmt.Sprintf(`{
		"url": "%s"
	}`, url))
	request := httptest.NewRequest(http.MethodPost, "/api/shorten", requestBody)
	SetContentTypeJson(request)
	request.Header.Set("Authorization", "Bearer "+accessToken)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
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

func TestFindAll(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter(db)
	defer db.Close()

	TruncateTable(db, "users")
	TruncateTable(db, "urls")

	accessToken := RegisterAndLoginUser(router)

	longUrls := []string{
		"https://azbagas.com",
		"https://google.com",
		"https://www.youtube.com/",
		"https://www.facebook.com/",
		"https://www.instagram.com/",
		"https://www.linkedin.com/",
		"https://www.github.com/",
	}

	// Shorten URL
	for _, url := range longUrls {
		createShortenedUrl(router, accessToken, url)
	}

	t.Run("Find all success", func(t *testing.T) {
		// Find all URLs
		request := httptest.NewRequest(http.MethodGet, "/api/shorten", nil)
		SetContentTypeJson(request)
		request.Header.Set("Authorization", "Bearer "+accessToken)
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, request)

		response := recorder.Result()
		assert.Equal(t, 200, response.StatusCode)

		var responseBody ResponseBody
		ReadResponseBody(response, &responseBody)
		
		data := responseBody["data"].([]interface{})
		assert.Equal(t, 5, len(data))

		for _, url := range data {
			urlData := url.(map[string]interface{})
			assert.NotNil(t, urlData["id"])
			assert.NotNil(t, urlData["url"])
			assert.NotNil(t, urlData["shortCode"])
			assert.NotNil(t, urlData["createdAt"])
			assert.NotNil(t, urlData["updatedAt"])
		}

		metadata := responseBody["metadata"].(map[string]interface{})
		assert.Equal(t, 1, int(metadata["currentPage"].(float64)))
		assert.Equal(t, 2, int(metadata["lastPage"].(float64)))
		assert.Equal(t, 5, int(metadata["perPage"].(float64)))
		assert.Equal(t, 7, int(metadata["total"].(float64)))
		assert.Equal(t, 1, int(metadata["from"].(float64)))
		assert.Equal(t, 5, int(metadata["to"].(float64)))
	})

	t.Run("Page query parameter success", func(t *testing.T) {
		// Find all URLs
		request := httptest.NewRequest(http.MethodGet, "/api/shorten", nil)
		SetContentTypeJson(request)
		request.Header.Set("Authorization", "Bearer "+accessToken)
		// Set query params
		q := request.URL.Query()
		q.Add("page", "1")
		q.Add("perPage", "10")
		request.URL.RawQuery = q.Encode()
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, request)

		response := recorder.Result()
		assert.Equal(t, 200, response.StatusCode)

		var responseBody ResponseBody
		ReadResponseBody(response, &responseBody)
		
		data := responseBody["data"].([]interface{})
		assert.Equal(t, 7, len(data))

		for _, url := range data {
			urlData := url.(map[string]interface{})
			assert.NotNil(t, urlData["id"])
			assert.NotNil(t, urlData["url"])
			assert.NotNil(t, urlData["shortCode"])
			assert.NotNil(t, urlData["createdAt"])
			assert.NotNil(t, urlData["updatedAt"])
		}

		metadata := responseBody["metadata"].(map[string]interface{})
		assert.Equal(t, 1, int(metadata["currentPage"].(float64)))
		assert.Equal(t, 1, int(metadata["lastPage"].(float64)))
		assert.Equal(t, 10, int(metadata["perPage"].(float64)))
		assert.Equal(t, 7, int(metadata["total"].(float64)))
		assert.Equal(t, 1, int(metadata["from"].(float64)))
		assert.Equal(t, 7, int(metadata["to"].(float64)))
	})

	t.Run("Page query parameter validation failed", func(t *testing.T) {
		// Find all URLs
		request := httptest.NewRequest(http.MethodGet, "/api/shorten", nil)
		SetContentTypeJson(request)
		request.Header.Set("Authorization", "Bearer "+accessToken)
		// Set query params
		q := request.URL.Query()
		q.Add("page", "1")
		q.Add("perPage", "4") // perPage is not 5, 10, 25
		request.URL.RawQuery = q.Encode()
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, request)

		response := recorder.Result()
		assert.Equal(t, 400, response.StatusCode)

		var responseBody ResponseBody
		ReadResponseBody(response, &responseBody)
		assert.NotNil(t, responseBody["message"])
		errors := responseBody["errors"].([]interface{})
		assert.NotEqual(t, 0, len(errors))
	})
}
