package web

type UrlUpdateRequest struct {
	ShortCode string `validate:"required" json:"shortCode"`
	Url       string `validate:"required,url" json:"url"`
}
