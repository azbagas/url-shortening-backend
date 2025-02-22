//go:build wireinject
// +build wireinject

package main

import (
	"net/http"

	"github.com/azbagas/url-shortening-backend/app"
	"github.com/azbagas/url-shortening-backend/middleware"
	"github.com/google/wire"
	"github.com/julienschmidt/httprouter"
)

func InitializedServer() *http.Server {
	wire.Build(
		app.NewConfig,
		app.NewRouter,
		middleware.NewLogMiddleware,
		wire.Bind(new(http.Handler), new(*httprouter.Router)),
		NewServer,
	)
	return nil
}
