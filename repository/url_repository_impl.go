package repository

import (
	"context"
	"database/sql"
	"errors"

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

func (repository *UrlRepositoryImpl) FindByShortCode(ctx context.Context, tx *sql.Tx, shortCode string) (domain.Url, error) {
	SQL := `SELECT id, user_id, url, short_code, created_at, updated_at 
					FROM urls
					WHERE short_code = $1`

	rows, err := tx.QueryContext(ctx, SQL, shortCode)
	helper.PanicIfError(err)
	defer rows.Close()

	url := domain.Url{}

	if rows.Next() {
		err = rows.Scan(&url.Id, &url.UserId, &url.Url, &url.ShortCode, &url.CreatedAt, &url.UpdatedAt)
		helper.PanicIfError(err)

		return url, nil
	} else {
		return url, errors.New("Url not found")
	}
}

func (repository *UrlRepositoryImpl) IncrementAccessCount(ctx context.Context, tx *sql.Tx, urlId int) {
	SQL := `INSERT INTO url_access_stats (url_id) VALUES ($1)`

	_, err := tx.ExecContext(ctx, SQL, urlId)
	helper.PanicIfError(err)
}

func (repository *UrlRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, url domain.Url) domain.Url {
	SQL := `UPDATE urls SET url = $1, updated_at = $2 WHERE short_code = $3`

	_, err := tx.ExecContext(ctx, SQL, url.Url, url.UpdatedAt, url.ShortCode)
	helper.PanicIfError(err)

	return url
}