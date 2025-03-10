package helper

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/azbagas/url-shortening-backend/model/domain"
	"github.com/azbagas/url-shortening-backend/model/web"
)

func ToUserResponse(user domain.User) web.UserResponse {
	var photoURL string
	if user.Photo.Valid {
		photoURL = user.Photo.String
	} else {
		photoURL = fmt.Sprintf("https://robohash.org/%s?set=set4", strings.ReplaceAll(user.Name, " ", "-"))
	}

	return web.UserResponse{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
		Photo: photoURL,
	}
}

func ToUrlResponse(url domain.Url) web.UrlResponse {
	return web.UrlResponse{
		Id:        url.Id,
		Url:       url.Url,
		ShortCode: url.ShortCode,
		CreatedAt: FormatToUTCString(url.CreatedAt),
		UpdatedAt: FormatToUTCString(url.UpdatedAt),
	}
}

func ToUrlResponses(urls []domain.Url) []web.UrlResponse {
	// If there is no data, return empty array
	if len(urls) == 0 {
		return []web.UrlResponse{}
	}

	var urlResponses []web.UrlResponse
	for _, url := range urls {
		urlResponses = append(urlResponses, ToUrlResponse(url))
	}

	return urlResponses
}

func ToPaginationResponse(page int, perPage int, totalData int) web.PaginationResponse {
	lastPage := math.Ceil(float64(totalData) / float64(perPage))

	from := (page-1)*perPage + 1

	to := page * perPage
	if page == int(lastPage) {
		to = totalData
	}

	if page > int(lastPage) {
		from = 0
		to = 0
	}

	return web.PaginationResponse{
		CurrentPage: page,
		LastPage:    int(lastPage),
		PerPage:     perPage,
		Total:       totalData,
		From:        from,
		To:          to,
	}
}

func ToUrlStatsResponse(url domain.Url, grandTotalAccessed int, lastAccessedAt time.Time, accessedDates []domain.UrlAccessTotalPerDate) web.UrlStatsResponse {
	return web.UrlStatsResponse{
		ShortUrl: ToUrlResponse(url),
		Stats:    ToUrlStatsDetailResponse(grandTotalAccessed, lastAccessedAt, accessedDates),
	}
}

func ToUrlStatsDetailResponse(grandTotalAccessed int, lastAccessedAt time.Time, urlStats []domain.UrlAccessTotalPerDate) web.UrlStatsDetailResponse {
	var lastAccessedAtResp string
	
	lastAccessedAtResp = FormatToUTCString(lastAccessedAt)

	return web.UrlStatsDetailResponse{
		GrandTotalAccessed: grandTotalAccessed,
		LastAccessedAt:     lastAccessedAtResp,
		AccessedDates:      ToUrlAccessTotalPerDateResponses(urlStats),
	}
}

func ToUrlAccessTotalPerDateResponses(urlStats []domain.UrlAccessTotalPerDate) []web.UrlAccessTotalPerDateResponse {
	// If there is no data, return empty array
	if len(urlStats) == 0 {
		return []web.UrlAccessTotalPerDateResponse{}
	}

	var urlAccessTotalPerDateResponses []web.UrlAccessTotalPerDateResponse
	for _, urlStat := range urlStats {
		urlAccessTotalPerDateResponses = append(urlAccessTotalPerDateResponses, web.UrlAccessTotalPerDateResponse{
			Date:          urlStat.Date,
			TotalAccessed: urlStat.TotalAccessed,
		})
	}

	return urlAccessTotalPerDateResponses
}
