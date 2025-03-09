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

	// Get query params
	var errors []web.ValidationErrorFieldMessage

	page, errResponse := helper.ConvertQueryParamToInt("page", request.URL.Query().Get("page"), "1")
	if errResponse != nil {
		errors = append(errors, *errResponse)
	}

	perPage, errResponse := helper.ConvertQueryParamToInt("perPage", request.URL.Query().Get("perPage"), "5")
	if errResponse != nil {
		errors = append(errors, *errResponse)
	}

	if len(errors) > 0 {
		helper.WriteToResponseBody(writer, http.StatusBadRequest, web.ValidationErrorResponse{Message: "Validation error", Errors: errors})
		return
	}
	// End of get query params

	urlFindAllRequest := web.UrlFindAllRequest{
		Page:    page,
		PerPage: perPage,
	}

	urlResponses, paginationResponse := controller.UrlService.FindAll(request.Context(), urlFindAllRequest, authUserId)
	dataResponse := web.DataWithPaginationResponse{
		Data: urlResponses,
		Metadata: paginationResponse,
	}

	helper.WriteToResponseBody(writer, http.StatusOK, dataResponse)
}

func (controller *UrlControllerImpl) FindByShortCode(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	shortCode := params.ByName("shortCode")

	urlResponse := controller.UrlService.FindByShortCode(request.Context(), shortCode)
	dataResponse := web.DataResponse{
		Data: urlResponse,
	}

	helper.WriteToResponseBody(writer, http.StatusOK, dataResponse)
}