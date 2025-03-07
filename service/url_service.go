package service

import (
	"context"

	"github.com/azbagas/url-shortening-backend/model/web"
)

type UrlService interface {
	Shorten(ctx context.Context, request web.UrlShortenRequest, authUserId int) web.UrlResponse
	FindAll(ctx context.Context, authUserId int) []web.UrlResponse
}
