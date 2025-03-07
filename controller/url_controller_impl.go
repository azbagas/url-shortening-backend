package controller

import (
	"net/http"

	"github.com/azbagas/url-shortening-backend/helper"
	"github.com/azbagas/url-shortening-backend/model/web"
	"github.com/azbagas/url-shortening-backend/service"
	"github.com/julienschmidt/httprouter"
)

type UrlControllerImpl struct {
	UrlService service.UrlService
}

func NewUrlController(urlService service.UrlService) UrlController {
	return &UrlControllerImpl{
		UrlService: urlService,
	}
}

func (controller *UrlControllerImpl) Shorten(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	authUserId := request.Context().Value("authUserId").(int)
	
	urlShortenRequest := web.UrlShortenRequest{}
	helper.ReadFromRequestBody(request, &urlShortenRequest)

	urlResponse := controller.UrlService.Shorten(request.Context(), urlShortenRequest, authUserId)
	dataResponse := web.DataResponse{
		Data: urlResponse,
	}

	helper.WriteToResponseBody(writer, http.StatusCreated, dataResponse)
}

func (controller *UrlControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	authUserId := request.Context().Value("authUserId").(int)

	urlResponses := controller.UrlService.FindAll(request.Context(), authUserId)
	dataResponse := web.DataResponse{
		Data: urlResponses,
	}

	helper.WriteToResponseBody(writer, http.StatusOK, dataResponse)
}