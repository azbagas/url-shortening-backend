package middleware

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/azbagas/url-shortening-backend/helper"
	"github.com/azbagas/url-shortening-backend/model/web"
	"github.com/azbagas/url-shortening-backend/token"
)

var publicRoutes = map[string]map[string]bool{
	"/api/users": {
		"POST": true,
	},
	"/api/users/login": {
		"POST": true,
	},
	"/api/users/refresh": {
		"POST": true,
	},
	"/api/ping": {
		"GET": true,
	},
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		if methods, ok := publicRoutes[request.URL.Path]; ok {
			if _, ok := methods[request.Method]; ok {
				next.ServeHTTP(writer, request)
				return
			}
		}

		// Verify access token
		authorizationHeader := request.Header.Get("Authorization")
		if !strings.Contains(authorizationHeader, "Bearer") {
			messageResponse := web.MessageResponse{
				Message: "Invalid token",
			}
		
			helper.WriteToResponseBody(writer, http.StatusUnauthorized, messageResponse)
			return
		}

		accessToken := strings.Replace(authorizationHeader, "Bearer ", "", -1)

		claims, err := token.VerifyAccessToken(accessToken)
		if err != nil {
			messageResponse := web.MessageResponse{
				Message: "Invalid token",
			}
		
			helper.WriteToResponseBody(writer, http.StatusUnauthorized, messageResponse)
			return
		}

		authUserIdStr := claims["sub"].(string)
		authUserId, err := strconv.Atoi(authUserIdStr)
		if err != nil {
			messageResponse := web.MessageResponse{
				Message: "Oops! Something went wrong",
			}
		
			helper.WriteToResponseBody(writer, http.StatusUnauthorized, messageResponse)
			return
		}

		// Set authUserId to context
		ctx := context.WithValue(context.Background(), "authUserId", authUserId)
		request = request.WithContext(ctx)

		next.ServeHTTP(writer, request)
	})
}
