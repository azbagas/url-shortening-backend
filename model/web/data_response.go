package web

type DataResponse struct {
	Data interface{} `json:"data"`
}

type DataWithPaginationResponse struct {
	Data     interface{}        `json:"data"`
	Metadata PaginationResponse `json:"metadata"`
}
