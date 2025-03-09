package service

import (
	"context"

	"github.com/azbagas/url-shortening-backend/model/web"
)

type UrlService interface {
	Shorten(ctx context.Context, request web.UrlShortenRequest, authUserId int) web.UrlResponse
	FindAll(ctx context.Context, request web.UrlFindAllRequest, authUserId int) ([]web.UrlResponse, web.PaginationResponse)
	FindByShortCode(ctx context.Context, shortCode string) web.UrlResponse
}
