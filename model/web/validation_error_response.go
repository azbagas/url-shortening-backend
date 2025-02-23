package web

type ValidationErrorResponse struct {
	Message string                        `json:"message"`
	Errors  []ValidationErrorFieldMessage `json:"errors"`
}

type ValidationErrorFieldMessage struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
