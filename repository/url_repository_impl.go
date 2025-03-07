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

func (repository *UrlRepositoryImpl) CountAll(ctx context.Context, tx *sql.Tx, userId int) int {
	SQL := `SELECT COUNT(id) FROM urls WHERE user_id = $1`
	
	var count int
	err := tx.QueryRowContext(ctx, SQL, userId).Scan(&count)
	helper.PanicIfError(err)
	
	return count
}

func (repository *UrlRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, userId int, page int, perPage int) []domain.Url {
	SQL := `SELECT id, user_id, url, short_code, created_at, updated_at 
					FROM urls
					WHERE user_id = $1
					ORDER BY created_at DESC
					LIMIT $2 OFFSET $3`
	
	rows, err := tx.QueryContext(ctx, SQL, userId, perPage, (page-1)*perPage)
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