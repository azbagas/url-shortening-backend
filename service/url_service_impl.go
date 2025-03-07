package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/azbagas/url-shortening-backend/helper"
	"github.com/azbagas/url-shortening-backend/model/domain"
	"github.com/azbagas/url-shortening-backend/model/web"
	"github.com/azbagas/url-shortening-backend/repository"
	"github.com/go-playground/validator/v10"
)

type UrlServiceImpl struct {
	UrlRepository repository.UrlRepository
	DB            *sql.DB
	Validate      *validator.Validate
}

func NewUrlService(urlRepository repository.UrlRepository, DB *sql.DB, validate *validator.Validate) UrlService {
	return &UrlServiceImpl{
		UrlRepository:          urlRepository,
		DB:                     DB,
		Validate:               validate,
	}
}

func (service *UrlServiceImpl) Shorten(ctx context.Context, request web.UrlShortenRequest, authUserId int) web.UrlResponse {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	// Generate random code
	randomCode, err := helper.GenerateRandomString(6)
	helper.PanicIfError(err)

	url := domain.Url{
		UserId:    authUserId,
		Url:       request.Url,
		ShortCode: randomCode,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	url = service.UrlRepository.Save(ctx, tx, url)

	return helper.ToUrlResponse(url)
}

func (service *UrlServiceImpl) FindAll(ctx context.Context, request web.UrlFindAllRequest, authUserId int) ([]web.UrlResponse, web.PaginationResponse) {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)
	
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	countUrls := service.UrlRepository.CountAll(ctx, tx, authUserId)
	urls := service.UrlRepository.FindAll(ctx, tx, authUserId, request.Page, request.PerPage)

	return helper.ToUrlResponses(urls), helper.ToPaginationResponse(request.Page, request.PerPage, countUrls)
}
