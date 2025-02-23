package test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterUserSuccess(t *testing.T) {
	db := SetupTestDB()
	TruncateTable(db, "users")
	router := SetupRouter(db)

	requestBody := strings.NewReader(`{
		"name": "Koseki Bijou",
		"email": "biboo@gmail.com",
		"password": "password",
		"passwordConfirmation": "password"
	}`)
	request := httptest.NewRequest(http.MethodPost, "/api/users", requestBody)
	request.Header.Add("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 201, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	data := responseBody["data"].(map[string]interface{})
	assert.NotNil(t, data["id"])
	assert.Equal(t, "Koseki Bijou", data["name"])
	assert.Equal(t, "biboo@gmail.com", data["email"])
	assert.NotNil(t, data["photo"])
}

func TestRegisterUserValidationFailed(t *testing.T) {
	db := SetupTestDB()
	TruncateTable(db, "users")
	router := SetupRouter(db)

	requestBody := strings.NewReader(`{
		"name": "Koseki Bijou",
		"email": "invalid email",
		"password": "password",
		"passwordConfirmation": "password"
	}`)
	request := httptest.NewRequest(http.MethodPost, "/api/users", requestBody)
	request.Header.Add("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 400, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, "Validation error", responseBody["message"])

	errors := responseBody["errors"].([]interface{})
	assert.NotEqual(t, 0, len(errors))
}

func TestRegisterUserEmailAlreadyRegistered(t *testing.T) {
	TestRegisterUserSuccess(t)

	db := SetupTestDB()
	router := SetupRouter(db)

	requestBody := strings.NewReader(`{
		"name": "Koseki Bijou",
		"email": "biboo@gmail.com",
		"password": "password",
		"passwordConfirmation": "password"
	}`)
	request := httptest.NewRequest(http.MethodPost, "/api/users", requestBody)
	request.Header.Add("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 409, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, "Email already registered", responseBody["message"])
}
