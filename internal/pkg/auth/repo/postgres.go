package repo

import (
	"MovieService/internal/models"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepo struct {
	db *pgxpool.Pool
}

func NewAuthRepo(db *pgxpool.Pool) *AuthRepo {
	return &AuthRepo{
		db: db,
	}
}

func (ar *AuthRepo) CreateUser(ctx context.Context, user *models.User) error {
	return nil
}

func (ar *AuthRepo) GetUserByLogin(context.Context, string) *models.User {
	return &models.User{}
}
