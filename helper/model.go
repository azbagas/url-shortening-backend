package helper

import (
	"fmt"
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
