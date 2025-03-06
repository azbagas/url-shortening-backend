package web

type UrlShortenRequest struct {
	Url string `validate:"required,url" json:"url"`
}
