package usecase

import (
	"MovieService/internal/models"
	"MovieService/internal/pkg/auth"
	"context"
	"fmt"
)

type AuthUsecase struct {
	repo auth.AuthRepo
}

func NewAuthUsecase(repo auth.AuthRepo) *AuthUsecase {
	return &AuthUsecase{
		repo: repo,
	}
}

func (au *AuthUsecase) SignIn(ctx context.Context, user *models.User) error {
	u, err := au.repo.GetUserByLogin(ctx, user.Login)
	if err != nil {
		return err
	}

	if u.Password == user.Password {
		return nil
	}

	return fmt.Errorf("forbidden")
}

func (au *AuthUsecase) SignUp(ctx context.Context, user *models.User) (int, error) {
	id, err := au.repo.CreateUser(ctx, user)
	return id, err
}
