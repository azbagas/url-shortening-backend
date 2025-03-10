package web

type UrlStatsResponse struct {
	ShortUrl UrlResponse            `json:"shortUrl"`
	Stats    UrlStatsDetailResponse `json:"stats"`
}

type UrlStatsDetailResponse struct {
	GrandTotalAccessed int                             `json:"grandTotalAccessed"`
	LastAccessedAt     string                          `json:"lastAccessedAt"`
	AccessedDates      []UrlAccessTotalPerDateResponse `json:"accessedDates"`
}

type UrlAccessTotalPerDateResponse struct {
	Date          string `json:"date"`
	TotalAccessed int    `json:"totalAccessed"`
}
