package helper

import (
	"fmt"
	"math"
	"strings"

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
