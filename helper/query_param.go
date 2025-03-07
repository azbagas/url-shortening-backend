package helper

import (
	"fmt"
	"strconv"

	"github.com/azbagas/url-shortening-backend/model/web"
)

func ConvertQueryParamToInt(key string, value string, defaultValue string) (int, *web.ValidationErrorFieldMessage) {
	if value == "" {
		value = defaultValue
	}

	result, err := strconv.Atoi(value)
	if err != nil {
		response := &web.ValidationErrorFieldMessage{
			Field:   key,
			Message: fmt.Sprintf("%s must be a number", key),
		}

		return 0, response
	}

	return result, nil
}
