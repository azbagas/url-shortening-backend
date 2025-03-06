package repository

import (
	"context"
	"database/sql"

	"github.com/azbagas/url-shortening-backend/helper"
	"github.com/azbagas/url-shortening-backend/model/domain"
)

type UrlRepositoryImpl struct {
}

func NewUrlRepository() UrlRepository {
	return &UrlRepositoryImpl{}
}

func (repository *UrlRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, url domain.Url) domain.Url {
	SQL := `INSERT INTO urls (user_id, url, short_code, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	
	err := tx.QueryRowContext(ctx, SQL, url.UserId, url.Url, url.ShortCode, url.UpdatedAt, url.UpdatedAt).Scan(&url.Id)
	helper.PanicIfError(err)
	
	return url
}