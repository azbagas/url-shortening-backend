package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type UrlController interface {
	Shorten(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindByShortCode(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
