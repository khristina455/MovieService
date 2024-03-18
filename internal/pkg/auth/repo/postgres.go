package repo

import (
	"MovieService/internal/models"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	createUser = `INSERT INTO "user" (login, password, is_admin) VALUES ($1, $2, $3) RETURNING id;`
)

type AuthRepo struct {
	db *pgxpool.Pool
}

func NewAuthRepo(db *pgxpool.Pool) *AuthRepo {
	return &AuthRepo{
		db: db,
	}
}

func (ar *AuthRepo) CreateUser(ctx context.Context, user *models.User) (int, error) {
	var id int
	err := ar.db.QueryRow(ctx, createUser,
		user.Login, user.Password, user.IsAdmin).Scan(&id)

	if err != nil {
		err = fmt.Errorf("error happened in scan.Scan: %w", err)

		return 0, err
	}

	return id, nil
}

func (ar *AuthRepo) GetUserByLogin(context.Context, string) *models.User {
	return &models.User{}
}
