package service

import (
	"context"

	"github.com/azbagas/url-shortening-backend/model/web"
)

type UserService interface {
	Register(ctx context.Context, request web.UserRegisterRequest) web.UserResponse
	Login(ctx context.Context, request web.UserLoginRequest) web.UserLoginResponseWithRefreshToken
	GetCurrentUser(ctx context.Context, authUserId int) web.UserResponse
	RefreshToken(ctx context.Context, refreshTokenRequest string) (web.NewAccessTokenResponse, error)
	Logout(ctx context.Context, refreshTokenRequest string) error
}
