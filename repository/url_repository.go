package repository

import (
	"context"
	"database/sql"

	"github.com/azbagas/url-shortening-backend/model/domain"
)

type UrlRepository interface {
	Save(ctx context.Context, tx *sql.Tx, url domain.Url) domain.Url
	FindAll(ctx context.Context, tx *sql.Tx) []domain.Url
}
