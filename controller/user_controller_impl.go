package controller

import (
	"net/http"

	"github.com/azbagas/url-shortening-backend/exception"
	"github.com/azbagas/url-shortening-backend/helper"
	"github.com/azbagas/url-shortening-backend/model/web"
	"github.com/azbagas/url-shortening-backend/service"
	"github.com/azbagas/url-shortening-backend/token"
	"github.com/julienschmidt/httprouter"
)

type UserControllerImpl struct {
	UserService service.UserService
}

func NewUserController(userService service.UserService) UserController {
	return &UserControllerImpl{
		UserService: userService,
	}
}

func (controller *UserControllerImpl) Register(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userRegisterRequest := web.UserRegisterRequest{}
	helper.ReadFromRequestBody(request, &userRegisterRequest)

	userResponse := controller.UserService.Register(request.Context(), userRegisterRequest)
	dataResponse := web.DataResponse{
		Data: userResponse,
	}

	helper.WriteToResponseBody(writer, http.StatusCreated, dataResponse)
}

func (controller *UserControllerImpl) Login(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	UserLoginRequest := web.UserLoginRequest{}
	helper.ReadFromRequestBody(request, &UserLoginRequest)
	UserLoginRequest.UserAgent = request.UserAgent()

	serviceLoginResponse := controller.UserService.Login(request.Context(), UserLoginRequest)

	// Set refresh token to cookie
	cookie := token.CreateRefreshTokenCookie(serviceLoginResponse.RefreshToken)
	http.SetCookie(writer, &cookie)

	loginResponse := web.UserLoginResponse{
		User:        serviceLoginResponse.User,
		AccessToken: serviceLoginResponse.AccessToken,
	}

	dataResponse := web.DataResponse{
		Data: loginResponse,
	}

	helper.WriteToResponseBody(writer, http.StatusOK, dataResponse)
}

func (controller *UserControllerImpl) GetCurrentUser(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	authUserId := request.Context().Value("authUserId").(int)

	userResponse := controller.UserService.GetCurrentUser(request.Context(), authUserId)
	dataResponse := web.DataResponse{
		Data: userResponse,
	}

	helper.WriteToResponseBody(writer, http.StatusOK, dataResponse)
}

func (controller *UserControllerImpl) RefreshToken(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// Get refresh token from cookie
	cookie, err := request.Cookie("refreshToken")
	if err != nil {
		panic(exception.NewUnauthorizedError("Refresh token is invalid. Please login again."))
	}

	refreshToken := cookie.Value
	newAccessTokenResponse, err := controller.UserService.RefreshToken(request.Context(), refreshToken)

	if err != nil {
		// Either refresh token is invalid or expired
		message := web.MessageResponse{
			Message: err.Error(),
		}

		// Delete refresh token cookie
		cookie := token.DeleteRefreshTokenCookie()
		http.SetCookie(writer, &cookie)

		helper.WriteToResponseBody(writer, http.StatusUnauthorized, message)
		return
	}

	dataResponse := web.DataResponse{
		Data: newAccessTokenResponse,
	}

	helper.WriteToResponseBody(writer, http.StatusOK, dataResponse)
}

func (controller *UserControllerImpl) Logout(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// Get refresh token from cookie
	cookie, err := request.Cookie("refreshToken")
	if err != nil {
		panic(exception.NewUnauthorizedError("Invalid refresh token."))
	}

	refreshToken := cookie.Value
	err = controller.UserService.Logout(request.Context(), refreshToken)

	if err != nil {
		message := web.MessageResponse{
			Message: err.Error(),
		}

		// Delete refresh token cookie
		cookieData := token.DeleteRefreshTokenCookie()
		http.SetCookie(writer, &cookieData)

		helper.WriteToResponseBody(writer, http.StatusUnauthorized, message)
		return
	}

	// Delete refresh token cookie
	cookieData := token.DeleteRefreshTokenCookie()
	http.SetCookie(writer, &cookieData)

	// Send empty access token
	dataResponse := web.DataResponse{
		Data: web.NewAccessTokenResponse{
			AccessToken: "",
		},
	}

	helper.WriteToResponseBody(writer, http.StatusOK, dataResponse)
}
