package auth

import (
	"MovieService/internal/models"
	"context"
)

type AuthRepo interface {
	CreateUser(context.Context, *models.User) (int, error)
	GetUserByLogin(context.Context, string) *models.User
}

type AuthUsecase interface {
	SignIn(context.Context, *models.User) error
	SignUp(context.Context, *models.User) (int, error)
}
