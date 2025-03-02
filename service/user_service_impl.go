package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/azbagas/url-shortening-backend/exception"
	"github.com/azbagas/url-shortening-backend/helper"
	"github.com/azbagas/url-shortening-backend/model/domain"
	"github.com/azbagas/url-shortening-backend/model/web"
	"github.com/azbagas/url-shortening-backend/repository"
	"github.com/azbagas/url-shortening-backend/token"
	"github.com/go-playground/validator/v10"
)

type UserServiceImpl struct {
	UserRepository         repository.UserRepository
	RefreshTokenRepository repository.RefreshTokenRepository
	DB                     *sql.DB
	Validate               *validator.Validate
}

func NewUserService(userRepository repository.UserRepository, refreshTokenRepository repository.RefreshTokenRepository, DB *sql.DB, validate *validator.Validate) UserService {
	return &UserServiceImpl{
		UserRepository:         userRepository,
		RefreshTokenRepository: refreshTokenRepository,
		DB:                     DB,
		Validate:               validate,
	}
}

func (service *UserServiceImpl) Register(ctx context.Context, request web.UserRegisterRequest) web.UserResponse {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	_, err = service.UserRepository.FindByEmail(ctx, tx, request.Email)
	if err == nil {
		// Email already registered
		panic(exception.NewConflictError("Email already registered"))
	}

	hashedPassword, _ := helper.HashPassword(request.Password)

	user := domain.User{
		Name:      request.Name,
		Email:     request.Email,
		Password:  hashedPassword,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	user = service.UserRepository.Save(ctx, tx, user)

	return helper.ToUserResponse(user)
}

func (service *UserServiceImpl) Login(ctx context.Context, request web.UserLoginRequest) web.UserLoginResponseWithRefreshToken {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	user, err := service.UserRepository.FindByEmail(ctx, tx, request.Email)
	if err != nil {
		panic(exception.NewUnauthorizedError("Email or password is incorrect"))
	}

	hashedPassword, _ := helper.HashPassword(request.Password)
	if !helper.CheckPasswordHash(request.Password, hashedPassword) {
		panic(exception.NewUnauthorizedError("Email or password is incorrect"))
	}

	// Create access token
	accessToken := token.CreateAccessToken(user.Id)

	// Create refresh token
	refreshToken := token.CreateRefreshToken(user.Id)

	// Save refresh token to database
	refreshTokenDomain := domain.RefreshToken{
		RefreshToken: refreshToken,
		UserId:       user.Id,
		UserAgent:    request.UserAgent,
	}
	service.RefreshTokenRepository.Save(ctx, tx, refreshTokenDomain)

	return web.UserLoginResponseWithRefreshToken{
		User:         helper.ToUserResponse(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}

func (service *UserServiceImpl) GetCurrentUser(ctx context.Context, authUserId int) web.UserResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	user, err := service.UserRepository.FindById(ctx, tx, authUserId)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	return helper.ToUserResponse(user)
}

func (service *UserServiceImpl) RefreshToken(ctx context.Context, refreshTokenRequest string) (web.NewAccessTokenResponse, error) {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	refreshTokenDomain, err := service.RefreshTokenRepository.FindByToken(ctx, tx, refreshTokenRequest)
	if err != nil {
		return web.NewAccessTokenResponse{}, errors.New("Refresh token is invalid. Please login again.")
	}

	// Verify refresh token
	_, err = token.VerifyRefreshToken(refreshTokenRequest)
	if err != nil {
		service.RefreshTokenRepository.Delete(ctx, tx, refreshTokenDomain)
		return web.NewAccessTokenResponse{}, errors.New(err.Error())
	}

	// Create new access token
	accessToken := token.CreateAccessToken(refreshTokenDomain.UserId)

	return web.NewAccessTokenResponse{
		AccessToken: accessToken,
	}, nil
}

func (service *UserServiceImpl) Logout(ctx context.Context, refreshTokenRequest string) error {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	refreshTokenDomain, err := service.RefreshTokenRepository.FindByToken(ctx, tx, refreshTokenRequest)
	if err != nil {
		return errors.New("Invalid refresh token.")
	}

	// Delete refresh token
	service.RefreshTokenRepository.Delete(ctx, tx, refreshTokenDomain)

	return nil
}
