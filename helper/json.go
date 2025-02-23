package helper

import (
	"encoding/json"
	"net/http"
)

func ReadFromRequestBody(request *http.Request, result interface{}) {
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(result)
	PanicIfError(err)
}

func WriteToResponseBody(writer http.ResponseWriter, httpCode int, response interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(httpCode)

	encoder := json.NewEncoder(writer)
	err := encoder.Encode(response)
	PanicIfError(err)
}