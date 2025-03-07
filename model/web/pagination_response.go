package web

type PaginationResponse struct {
	CurrentPage int `json:"currentPage"`
	LastPage    int `json:"lastPage"`
	PerPage     int `json:"perPage"`
	Total       int `json:"total"`
	From        int `json:"from"`
	To          int `json:"to"`
}
