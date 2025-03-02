package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type UserController interface {
	Register(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Login(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	GetCurrentUser(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	RefreshToken(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Logout(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
