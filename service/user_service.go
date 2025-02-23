package service

import (
	"context"

	"github.com/azbagas/url-shortening-backend/model/web"
)

type UserService interface {
	Register(ctx context.Context, request web.UserRegisterRequest) web.UserResponse
}
