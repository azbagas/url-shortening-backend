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

func (repository *UrlRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) []domain.Url {
	SQL := `SELECT id, user_id, url, short_code, created_at, updated_at FROM urls`
	
	rows, err := tx.QueryContext(ctx, SQL)
	helper.PanicIfError(err)
	
	defer rows.Close()
	
	var urls []domain.Url
	for rows.Next() {
		url := domain.Url{}
		err := rows.Scan(&url.Id, &url.UserId, &url.Url, &url.ShortCode, &url.CreatedAt, &url.UpdatedAt)
		helper.PanicIfError(err)
		
		urls = append(urls, url)
	}
	
	return urls
}