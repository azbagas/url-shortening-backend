package repository

import (
	"context"
	"database/sql"

	"github.com/azbagas/url-shortening-backend/model/domain"
)

type UserRepository interface {
	Save(ctx context.Context, tx *sql.Tx, user domain.User) domain.User
	FindByEmail(ctx context.Context, tx *sql.Tx, email string) (domain.User, error)
	FindById(ctx context.Context, tx *sql.Tx, userId int) (domain.User, error)
}