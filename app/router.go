package app

import (
	"fmt"
	"net/http"

	"github.com/azbagas/url-shortening-backend/controller"
	"github.com/azbagas/url-shortening-backend/exception"
	"github.com/julienschmidt/httprouter"
)

func NewRouter(userController controller.UserController) *httprouter.Router {
	router := httprouter.New()

	router.GET("/api/ping", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		fmt.Fprint(writer, "pong")
	})
	router.POST("/api/users", userController.Register)

	router.PanicHandler = exception.ErrorHandler

	return router
}
