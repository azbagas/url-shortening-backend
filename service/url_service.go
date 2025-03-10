package service

import (
	"context"

	"github.com/azbagas/url-shortening-backend/model/web"
)

type UrlService interface {
	Shorten(ctx context.Context, request web.UrlShortenRequest, authUserId int) web.UrlResponse
	FindAll(ctx context.Context, request web.UrlFindAllRequest, authUserId int) ([]web.UrlResponse, web.PaginationResponse)
	FindByShortCode(ctx context.Context, shortCode string) web.UrlResponse
	Update(ctx context.Context, request web.UrlUpdateRequest, authUserId int) web.UrlResponse
	Delete(ctx context.Context, shortCode string, authUserId int)
	GetStats(ctx context.Context, request web.UrlStatsRequest, authUserId int) web.UrlStatsResponse
}
