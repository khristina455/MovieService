package usecase

import (
	"MovieService/internal/models"
	"MovieService/internal/pkg/auth"
	"context"
)

type AuthUsecase struct {
	repo auth.AuthRepo
}

func NewAuthUsecase(repo auth.AuthRepo) *AuthUsecase {
	return &AuthUsecase{
		repo: repo,
	}
}

func (ac *AuthUsecase) SignIn(ctx context.Context, user *models.User) error {
	return nil
}

func (ac *AuthUsecase) SignUp(ctx context.Context, user *models.User) error {
	return nil
}
