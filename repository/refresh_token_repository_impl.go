package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/azbagas/url-shortening-backend/helper"
	"github.com/azbagas/url-shortening-backend/model/domain"
)

type RefreshTokenRepositoryImpl struct {
}

func NewRefreshTokenRepository() RefreshTokenRepository {
	return &RefreshTokenRepositoryImpl{}
}

func (repository *RefreshTokenRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, refreshToken domain.RefreshToken) domain.RefreshToken {
	SQL := `INSERT INTO refresh_tokens (refresh_token, user_id, user_agent) VALUES ($1, $2, $3) RETURNING id`
	
	err := tx.QueryRowContext(ctx, SQL, refreshToken.RefreshToken, refreshToken.UserId, refreshToken.UserAgent).Scan(&refreshToken.Id)
	helper.PanicIfError(err)
	
	return refreshToken
}

func (repository *RefreshTokenRepositoryImpl) FindByToken(ctx context.Context, tx *sql.Tx, refreshToken string) (domain.RefreshToken, error) {
	SQL := `SELECT id, refresh_token, user_id, user_agent FROM refresh_tokens WHERE refresh_token = $1`

	rows, err := tx.QueryContext(ctx, SQL, refreshToken)
	helper.PanicIfError(err)
	defer rows.Close()

	refreshTokenDomain := domain.RefreshToken{}

	if rows.Next() {
		err = rows.Scan(&refreshTokenDomain.Id, &refreshTokenDomain.RefreshToken, &refreshTokenDomain.UserId, &refreshTokenDomain.UserAgent)
		helper.PanicIfError(err)

		return refreshTokenDomain, nil
	} else {
		return refreshTokenDomain, errors.New("Refresh token not found")
	}
}

func (repository *RefreshTokenRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, refreshToken domain.RefreshToken) {
	SQL := `DELETE FROM refresh_tokens WHERE refresh_token = $1`

	_, err := tx.ExecContext(ctx, SQL, refreshToken.RefreshToken)
	helper.PanicIfError(err)
}
