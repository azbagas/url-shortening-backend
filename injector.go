//go:build wireinject
// +build wireinject

package main

import (
	"net/http"

	"github.com/azbagas/url-shortening-backend/app"
	"github.com/azbagas/url-shortening-backend/controller"
	"github.com/azbagas/url-shortening-backend/repository"
	"github.com/azbagas/url-shortening-backend/service"
	"github.com/google/wire"
)

var userSet = wire.NewSet(
	repository.NewUserRepository,
	repository.NewRefreshTokenRepository,
	service.NewUserService,
	controller.NewUserController,
)

var urlSet = wire.NewSet(
	repository.NewUrlRepository,
	service.NewUrlService,
	controller.NewUrlController,
)

func InitializedServer() *http.Server {
	wire.Build(
		app.NewDB,
		app.NewValidator,
		userSet,
		urlSet,
		app.NewRouter,
		NewServer,
	)
	return nil
}
