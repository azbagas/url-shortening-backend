package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/azbagas/url-shortening-backend/exception"
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
		UrlRepository: urlRepository,
		DB:            DB,
		Validate:      validate,
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

func (service *UrlServiceImpl) FindByShortCode(ctx context.Context, shortCode string) web.UrlResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	url, err := service.UrlRepository.FindByShortCode(ctx, tx, shortCode)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	// Increment access count
	service.UrlRepository.IncrementAccessCount(ctx, tx, url.Id)

	return helper.ToUrlResponse(url)
}

func (service *UrlServiceImpl) Update(ctx context.Context, request web.UrlUpdateRequest, authUserId int) web.UrlResponse {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	url, err := service.UrlRepository.FindByShortCode(ctx, tx, request.ShortCode)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	if url.UserId != authUserId {
		panic(exception.NewForbiddenError("You don't have permission to update this url"))
	}

	url.Url = request.Url
	url.UpdatedAt = time.Now().UTC()

	url = service.UrlRepository.Update(ctx, tx, url)

	return helper.ToUrlResponse(url)
}

func (service *UrlServiceImpl) Delete(ctx context.Context, shortCode string, authUserId int) {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	url, err := service.UrlRepository.FindByShortCode(ctx, tx, shortCode)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	if url.UserId != authUserId {
		panic(exception.NewForbiddenError("You don't have permission to delete this url"))
	}

	service.UrlRepository.Delete(ctx, tx, url)
}

func (service *UrlServiceImpl) GetStats(ctx context.Context, request web.UrlStatsRequest, authUserId int) web.UrlStatsResponse {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	url, err := service.UrlRepository.FindByShortCode(ctx, tx, request.ShortCode)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	if url.UserId != authUserId {
		panic(exception.NewForbiddenError("You don't have permission to get stats for this url"))
	}

	// Create a variable to store user current time with request timezone
	userCurrentTime := helper.GetCurrentTime(request.Timezone)

	var startTime time.Time
	switch request.TimeRange {
	case "7d":
		startTime = userCurrentTime.AddDate(0, 0, -6)
	case "30d":
		startTime = userCurrentTime.AddDate(0, 0, -29)
	case "90d":
		startTime = userCurrentTime.AddDate(0, 0, -89)
	}

	var dateMap []domain.UrlAccessTotalPerDate

	for t := startTime; !t.After(userCurrentTime); t = t.AddDate(0, 0, 1) {
		dateStr := t.Format("2006-01-02")

		urlAccess := domain.UrlAccessTotalPerDate{
			Date:          dateStr,
			TotalAccessed: 0,
		}

		dateMap = append(dateMap, urlAccess)
	}

	startDate := startTime.Format("2006-01-02")
	endDate := userCurrentTime.Format("2006-01-02")

	accessStats := service.UrlRepository.TotalAccessedPerDate(ctx, tx, url.Id, request.Timezone, startDate, endDate)
	grandTotalAccessed := service.UrlRepository.GrandTotalAccessed(ctx, tx, url.Id)
	lastAccessed, _ := service.UrlRepository.LastAccessed(ctx, tx, url.Id)

	dateMapIndex := make(map[string]*domain.UrlAccessTotalPerDate)
	for i := range dateMap {
		dateMapIndex[dateMap[i].Date] = &dateMap[i] // Store pointer for direct updates
	}

	for _, stats := range accessStats {
		if entry, exists := dateMapIndex[stats.Date]; exists {
			entry.TotalAccessed = stats.TotalAccessed
		}
	}

	return helper.ToUrlStatsResponse(url, grandTotalAccessed, lastAccessed, dateMap)
}
