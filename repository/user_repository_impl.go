package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/azbagas/url-shortening-backend/helper"
	"github.com/azbagas/url-shortening-backend/model/domain"
)

type UserRepositoryImpl struct {
}

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{}
}

func (repository *UserRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, user domain.User) domain.User {
	SQL := "INSERT INTO users(name, email, password, photo, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"

	err := tx.QueryRowContext(ctx, SQL, user.Name, user.Email, user.Password, user.Photo, user.CreatedAt, user.UpdatedAt).Scan(&user.Id)
	helper.PanicIfError(err)

	return user
}

func (repository *UserRepositoryImpl) FindByEmail(ctx context.Context, tx *sql.Tx, email string) (domain.User, error) {
	SQL := "SELECT id, name, email, password, photo, created_at, updated_at FROM users WHERE email = $1"
	
	rows, err := tx.QueryContext(ctx, SQL, email)
	helper.PanicIfError(err)
	defer rows.Close()

	user := domain.User{}

	if rows.Next() {
		err = rows.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.Photo, &user.CreatedAt, &user.UpdatedAt)
		helper.PanicIfError(err)

		return user, nil
	} else {
		return user, errors.New("User not found")
	}
}

func (repository *UserRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (domain.User, error) {
	SQL := "SELECT id, name, email, password, photo, created_at, updated_at FROM users WHERE id = $1"
	
	rows, err := tx.QueryContext(ctx, SQL, id)
	helper.PanicIfError(err)
	defer rows.Close()

	user := domain.User{}

	if rows.Next() {
		err = rows.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.Photo, &user.CreatedAt, &user.UpdatedAt)
		helper.PanicIfError(err)

		return user, nil
	} else {
		return user, errors.New("User not found")
	}
}
