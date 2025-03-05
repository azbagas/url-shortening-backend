package test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRegisterUser(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter(db)
	defer db.Close()

	t.Run("Register success", func(t *testing.T) {
		TruncateTable(db, "users")

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

		response := recorder.Result()
		assert.Equal(t, 201, response.StatusCode)

		var responseBody ResponseBody
		ReadResponseBody(response, &responseBody)
		data := responseBody["data"].(map[string]interface{})
		assert.NotNil(t, data["id"])
		assert.Equal(t, "Koseki Bijou", data["name"])
		assert.Equal(t, "biboo@gmail.com", data["email"])
		assert.NotNil(t, data["photo"])
	})

	t.Run("Validation failed", func(t *testing.T) {
		TruncateTable(db, "users")

		requestBody := strings.NewReader(`{
			"name": "Koseki Bijou",
			"email": "invalid email",
			"password": "password",
			"passwordConfirmation": "password"
		}`)
		request := httptest.NewRequest(http.MethodPost, "/api/users", requestBody)
		SetContentTypeJson(request)
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

	t.Run("Email already registered", func(t *testing.T) {
		TruncateTable(db, "users")

		// First request
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

		// Second request
		requestBody = strings.NewReader(`{
			"name": "Koseki Bijou",
			"email": "biboo@gmail.com",
			"password": "password",
			"passwordConfirmation": "password"
		}`)
		request = httptest.NewRequest(http.MethodPost, "/api/users", requestBody)
		SetContentTypeJson(request)
		recorder = httptest.NewRecorder()
		router.ServeHTTP(recorder, request)

		// Capture the second response
		response := recorder.Result()
		assert.Equal(t, 409, response.StatusCode)

		var responseBody ResponseBody
		ReadResponseBody(response, &responseBody)
		assert.NotNil(t, responseBody["message"])
	})
}

func TestLoginUser(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter(db)
	defer db.Close()

	t.Run("Login success", func(t *testing.T) {
		TruncateTable(db, "users")

		// First request: Register user
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

		// Second request: Login user
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

		response := recorder.Result()
		assert.Equal(t, 200, response.StatusCode)

		var responseBody ResponseBody
		ReadResponseBody(response, &responseBody)
		data := responseBody["data"].(map[string]interface{})
		assert.NotNil(t, data["accessToken"])
		assert.NotNil(t, data["user"])
		user := data["user"].(map[string]interface{})
		assert.NotNil(t, user["id"])
		assert.Equal(t, "Koseki Bijou", user["name"])
		assert.Equal(t, "biboo@gmail.com", user["email"])
		assert.NotNil(t, user["photo"])
	})

	t.Run("Login failed: Validation error", func(t *testing.T) {
		TruncateTable(db, "users")

		// First request: Register user
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

		// Second request: Login user
		requestBody = strings.NewReader(`{
			"email": "invalidemail",
			"password": "password"
		}`)
		request = httptest.NewRequest(http.MethodPost, "/api/users/login", requestBody)
		SetContentTypeJson(request)
		// Set user agent
		request.Header.Set("User-Agent", "Mozilla/5.0")
		recorder = httptest.NewRecorder()
		router.ServeHTTP(recorder, request)

		response := recorder.Result()
		assert.Equal(t, 400, response.StatusCode)

		var responseBody ResponseBody
		ReadResponseBody(response, &responseBody)
		assert.NotNil(t, responseBody["message"])
		errors := responseBody["errors"].([]interface{})
		assert.NotEqual(t, 0, len(errors))
	})

	t.Run("Login failed: email or password is incorrect", func(t *testing.T) {
		TruncateTable(db, "users")

		// First request: Register user
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

		// Second request: Login user
		requestBody = strings.NewReader(`{
			"email": "wrongemail@gmail.com",
			"password": "password"
		}`)
		request = httptest.NewRequest(http.MethodPost, "/api/users/login", requestBody)
		SetContentTypeJson(request)
		// Set user agent
		request.Header.Set("User-Agent", "Mozilla/5.0")
		recorder = httptest.NewRecorder()
		router.ServeHTTP(recorder, request)

		response := recorder.Result()
		assert.Equal(t, 401, response.StatusCode)
		var responseBody ResponseBody
		ReadResponseBody(response, &responseBody)
		assert.NotNil(t, responseBody["message"])
	})
}

func TestGetCurrentUser(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter(db)
	defer db.Close()

	t.Run("Get current user success", func(t *testing.T) {
		TruncateTable(db, "users")

		// First request: Register user
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

		// Second request: Login user
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

		// Third request: Get current user
		request = httptest.NewRequest(http.MethodGet, "/api/users/current", nil)
		SetContentTypeJson(request)
		request.Header.Set("Authorization", "Bearer " + accessToken)
		recorder = httptest.NewRecorder()
		router.ServeHTTP(recorder, request)

		response := recorder.Result()
		assert.Equal(t, 200, response.StatusCode)

		ReadResponseBody(response, &responseBody)
		data := responseBody["data"].(map[string]interface{})
		assert.NotNil(t, data["id"])
		assert.Equal(t, "Koseki Bijou", data["name"])
		assert.Equal(t, "biboo@gmail.com", data["email"])
		assert.NotNil(t, data["photo"])
	})

	t.Run("Get current user failed", func(t *testing.T) {
		TruncateTable(db, "users")

		request := httptest.NewRequest(http.MethodGet, "/api/users/current", nil)
		SetContentTypeJson(request)
		request.Header.Set("Authorization", "Bearer " + "InvalidAccessToken")
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, request)

		response := recorder.Result()
		assert.Equal(t, 401, response.StatusCode)

		var responseBody ResponseBody
		ReadResponseBody(response, &responseBody)
		assert.NotNil(t, responseBody["message"])
	})
}

func TestGetNewAccessToken(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter(db)
	defer db.Close()

	t.Run("Get new access token success", func(t *testing.T) {
		TruncateTable(db, "users")

		// First request: Register user
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

		// Second request: Login user
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
		// Get first access token
		var responseBody ResponseBody
		ReadResponseBody(recorder.Result(), &responseBody)
		firstAccessToken := responseBody["data"].(map[string]interface{})["accessToken"].(string)
		
		// Extract refresh token from cookies
		refreshToken := ""
		for _, cookie := range recorder.Result().Cookies() {
			if cookie.Name == "refreshToken" {
				refreshToken = cookie.Value
				break
			}
		}
		assert.NotEmpty(t, refreshToken)

		// Sleep for a moment to make sure the new access token is different
		time.Sleep(500 * time.Millisecond)

		// Third request: Get new access token
		request = httptest.NewRequest(http.MethodPost, "/api/users/refresh", nil)
		SetContentTypeJson(request)
		request.AddCookie(&http.Cookie{Name: "refreshToken", Value: refreshToken})
		recorder = httptest.NewRecorder()
		router.ServeHTTP(recorder, request)

		response := recorder.Result()
		assert.Equal(t, 200, response.StatusCode)

		ReadResponseBody(response, &responseBody)
		data := responseBody["data"].(map[string]interface{})
		assert.NotNil(t, data["accessToken"])
		// Check if the new access token is different from the old one
		newAccessToken := data["accessToken"].(string)
		assert.NotEqual(t, firstAccessToken, newAccessToken)
	})

	t.Run("Failed get new access token", func(t *testing.T) {
		TruncateTable(db, "users")

		request := httptest.NewRequest(http.MethodPost, "/api/users/refresh", nil)
		SetContentTypeJson(request)
		request.AddCookie(&http.Cookie{Name: "refreshToken", Value: "InvalidRefreshToken"})
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, request)

		response := recorder.Result()
		assert.Equal(t, 401, response.StatusCode)

		var responseBody ResponseBody
		ReadResponseBody(response, &responseBody)
		assert.NotNil(t, responseBody["message"])
	})
}