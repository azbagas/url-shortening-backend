package web

type UrlFindAllRequest struct {
	// From query parameter
	Page    int `validate:"numeric,gt=0" json:"page"`
	PerPage int `validate:"oneof=5 10 25" json:"perPage"`
}
