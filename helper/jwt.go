package helper

import (
	"net/http"
	"strconv"
	"time"

	"github.com/azbagas/url-shortening-backend/config"
	"github.com/golang-jwt/jwt/v5"
)

func CreateAccessToken(userId int) string {
	accessTokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": config.AppConfig.AppName,
		"sub": strconv.Itoa(userId),
		"exp": time.Now().Add(time.Minute * 15).Unix(), // 15 minutes
	})
	accessToken, err := accessTokenClaims.SignedString([]byte(config.AppConfig.AccessTokenSecret))
	PanicIfError(err)

	return accessToken
}

func CreateRefreshToken(userId int) string {
	refreshTokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": config.AppConfig.AppName,
		"sub": strconv.Itoa(userId),
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(), // 30 days
	})
	refreshToken, err := refreshTokenClaims.SignedString([]byte(config.AppConfig.RefreshTokenSecret))
	PanicIfError(err)

	return refreshToken
}

func VerifyAccessToken(accessToken string) (jwt.MapClaims, error) {
	token, err := jwt.ParseWithClaims(accessToken, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.AccessTokenSecret), nil
	})
	
	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, err
	}

	return claims, nil
}

func VerifyRefreshToken(refreshToken string) (jwt.MapClaims, error) {
	token, err := jwt.ParseWithClaims(refreshToken, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.RefreshTokenSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, err
	}

	return claims, nil
}

func CreateRefreshTokenCookie(refreshToken string) http.Cookie {
	cookie := http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/api/users",
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(24 * time.Hour * 30), // 30 days
	}

	return cookie
}

func DeleteRefreshTokenCookie() http.Cookie {
	cookie := http.Cookie{
		Name:     "refreshToken",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		Path:     "/api/users",
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	}

	return cookie
}
