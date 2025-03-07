package app

import (
	"fmt"
	"net/http"

	"github.com/azbagas/url-shortening-backend/controller"
	"github.com/azbagas/url-shortening-backend/exception"
	"github.com/julienschmidt/httprouter"
)

func NewRouter(userController controller.UserController, urlController controller.UrlController) *httprouter.Router {
	router := httprouter.New()

	router.GET("/api/ping", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		fmt.Fprint(writer, "pong")
	})

	router.POST("/api/users", userController.Register)
	router.POST("/api/users/login", userController.Login)
	router.GET("/api/users/current", userController.GetCurrentUser)
	router.POST("/api/users/refresh", userController.RefreshToken)
	router.DELETE("/api/users/logout", userController.Logout)

	router.POST("/api/shorten", urlController.Shorten)
	router.GET("/api/shorten", urlController.FindAll)

	router.PanicHandler = exception.ErrorHandler

	return router
}
