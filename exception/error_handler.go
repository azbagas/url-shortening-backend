package exception

import (
	"net/http"
	"strings"

	"github.com/azbagas/url-shortening-backend/helper"
	"github.com/azbagas/url-shortening-backend/model/web"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

func ErrorHandler(writer http.ResponseWriter, request *http.Request, err interface{}) {

	if notFoundError(writer, request, err) {
		return
	}

	if validationErrors(writer, request, err) {
		return
	}

	if unauthorizedError(writer, request, err) {
		return
	}

	if forbiddenError(writer, request, err) {
		return
	}

	if conflictError(writer, request, err) {
		return
	}

	internalServerError(writer, request, err)
}

func validationErrors(writer http.ResponseWriter, _ *http.Request, err interface{}) bool {
	exception, ok := err.(validator.ValidationErrors)

	if !ok {
		return false
	}

	var errors []web.ValidationErrorFieldMessage
	for _, e := range exception {
		field := e.Field()
		msg := helper.GetValidationErrorMessage(e)

		// WARNING: It assume that all the json field name is camelCase
		// Example:
		// In struct: FirstName, JSON: firstName
		field = strings.ToLower(field[:1]) + field[1:]

		errorMessage := web.ValidationErrorFieldMessage{
			Field:   field,
			Message: msg,
		}

		errors = append(errors, errorMessage)
	}

	ValidationErrorResponse := web.ValidationErrorResponse{
		Message: "Validation error",
		Errors:  errors,
	}

	helper.WriteToResponseBody(writer, http.StatusBadRequest, ValidationErrorResponse)
	return true
}

func notFoundError(writer http.ResponseWriter, _ *http.Request, err interface{}) bool {
	exception, ok := err.(NotFoundError)

	if !ok {
		return false
	}

	messageResponse := web.MessageResponse{
		Message: exception.Error,
	}

	helper.WriteToResponseBody(writer, http.StatusNotFound, messageResponse)
	return true
}

func unauthorizedError(writer http.ResponseWriter, _ *http.Request, err interface{}) bool {
	exception, ok := err.(UnauthorizedError)

	if !ok {
		return false
	}

	messageResponse := web.MessageResponse{
		Message: exception.Error,
	}

	helper.WriteToResponseBody(writer, http.StatusUnauthorized, messageResponse)
	return true
}

func forbiddenError(writer http.ResponseWriter, _ *http.Request, err interface{}) bool {
	exception, ok := err.(ForbiddenError)

	if !ok {
		return false
	}

	messageResponse := web.MessageResponse{
		Message: exception.Error,
	}

	helper.WriteToResponseBody(writer, http.StatusForbidden, messageResponse)
	return true
}

func conflictError(writer http.ResponseWriter, _ *http.Request, err interface{}) bool {
	exception, ok := err.(ConflictError)

	if !ok {
		return false
	}

	messageResponse := web.MessageResponse{
		Message: exception.Error,
	}

	helper.WriteToResponseBody(writer, http.StatusConflict, messageResponse)
	return true
}

func internalServerError(writer http.ResponseWriter, _ *http.Request, err interface{}) {
	messageResponse := web.MessageResponse{
		Message: "Oops! Something went wrong",
	}

	logrus.Error(err)

	helper.WriteToResponseBody(writer, http.StatusInternalServerError, messageResponse)
}
