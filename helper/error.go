package helper

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func PanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func GetValidationErrorMessage(e validator.FieldError) string {
	// WARNING: It assume that all the json field name is camelCase
	// Example:
	// In struct: FirstName, JSON: firstName
	field := strings.ToLower(e.Field()[:1]) + e.Field()[1:]

	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", field, e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters long", field, e.Param())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, e.Param())
	case "lt":
		return fmt.Sprintf("%s must be less than %s", field, e.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, e.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, e.Param())
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters long", field, e.Param())
	case "numeric":
		return fmt.Sprintf("%s must be a number", field)
	case "oneof":
		return fmt.Sprintf("%s must be one of %s", field, e.Param())
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}
