//go:build wireinject
// +build wireinject

package main

import (
	"net/http"

	"github.com/azbagas/url-shortening-backend/app"
	"github.com/azbagas/url-shortening-backend/controller"
	"github.com/azbagas/url-shortening-backend/middleware"
	"github.com/azbagas/url-shortening-backend/repository"
	"github.com/azbagas/url-shortening-backend/service"
	"github.com/google/wire"
	"github.com/julienschmidt/httprouter"
)

var userSet = wire.NewSet(
	repository.NewUserRepository,
	service.NewUserService,
	controller.NewUserController,
)

func InitializedServer() *http.Server {
	wire.Build(
		app.NewConfig,
		app.NewDB,
		app.NewValidator,
		userSet,
		app.NewRouter,
		middleware.NewLogMiddleware,
		wire.Bind(new(http.Handler), new(*httprouter.Router)),
		NewServer,
	)
	return nil
}
