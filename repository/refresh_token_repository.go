package repository

import (
	"context"
	"database/sql"

	"github.com/azbagas/url-shortening-backend/model/domain"
)

type RefreshTokenRepository interface {
	Save(ctx context.Context, tx *sql.Tx, refreshToken domain.RefreshToken) domain.RefreshToken
	FindByToken(ctx context.Context, tx *sql.Tx, refreshToken string) (domain.RefreshToken, error)
	Delete(ctx context.Context, tx *sql.Tx, refreshToken domain.RefreshToken)
}
