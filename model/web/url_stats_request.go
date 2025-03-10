package web

type UrlStatsRequest struct {
	ShortCode string `validate:"required" json:"shortCode"`
	Timezone  string `validate:"timezone" json:"timezone"`
	TimeRange string `validate:"required,oneof=7d 30d 90d" json:"timeRange"`
}
