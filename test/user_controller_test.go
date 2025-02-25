package test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterUser(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter(db)

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

		assert.Equal(t, "Validation error", responseBody["message"])

		errors := responseBody["errors"].([]interface{})
		assert.NotEqual(t, 0, len(errors))
	})

	t.Run("Email already registered", func(t *testing.T) {
		TruncateTable(db, "users")

		// First request
		requestBody1 := strings.NewReader(`{
			"name": "Koseki Bijou",
			"email": "biboo@gmail.com",
			"password": "password",
			"passwordConfirmation": "password"
		}`)
		request1 := httptest.NewRequest(http.MethodPost, "/api/users", requestBody1)
		SetContentTypeJson(request1)
		recorder1 := httptest.NewRecorder()
		router.ServeHTTP(recorder1, request1)

		// Second request
		requestBody2 := strings.NewReader(`{
			"name": "Koseki Bijou",
			"email": "biboo@gmail.com",
			"password": "password",
			"passwordConfirmation": "password"
		}`)
		request2 := httptest.NewRequest(http.MethodPost, "/api/users", requestBody2)
		SetContentTypeJson(request2)
		recorder2 := httptest.NewRecorder()
		router.ServeHTTP(recorder2, request2)

		// Capture the second response
		response := recorder2.Result()
		assert.Equal(t, 409, response.StatusCode)

		var responseBody ResponseBody
		ReadResponseBody(response, &responseBody)

		assert.Equal(t, "Email already registered", responseBody["message"])
	})
}
