package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/azbagas/url-shortening-backend/model/domain"
)

type UrlRepository interface {
	Save(ctx context.Context, tx *sql.Tx, url domain.Url) domain.Url
	CountAll(ctx context.Context, tx *sql.Tx, userId int) int
	FindAll(ctx context.Context, tx *sql.Tx, userId int, page int, perPage int) []domain.Url
	FindByShortCode(ctx context.Context, tx *sql.Tx, shortCode string) (domain.Url, error)
	IncrementAccessCount(ctx context.Context, tx *sql.Tx, urlId int)
	Update(ctx context.Context, tx *sql.Tx, url domain.Url) domain.Url
	Delete(ctx context.Context, tx *sql.Tx, url domain.Url)
	TotalAccessedPerDate(ctx context.Context, tx *sql.Tx, urlId int, timezone string, startDate string, endDate string) []domain.UrlAccessTotalPerDate
	GrandTotalAccessed(ctx context.Context, tx *sql.Tx, urlId int) int
	LastAccessed(ctx context.Context, tx *sql.Tx, urlId int) (time.Time, error)
}
