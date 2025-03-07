package repository

import (
	"context"
	"database/sql"

	"github.com/azbagas/url-shortening-backend/model/domain"
)

type UrlRepository interface {
	Save(ctx context.Context, tx *sql.Tx, url domain.Url) domain.Url
	CountAll(ctx context.Context, tx *sql.Tx, userId int) int
	FindAll(ctx context.Context, tx *sql.Tx, userId int, page int, perPage int) []domain.Url
}
