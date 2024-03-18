package usecase

import (
	"MovieService/internal/models"
	"MovieService/internal/pkg/actors"
	"context"
	"github.com/jackc/pgx/v5/pgtype"
)

type ActorsUsecase struct {
	repo actors.ActorsRepo
}

func NewActorsUsecase(repo actors.ActorsRepo) *ActorsUsecase {
	return &ActorsUsecase{
		repo: repo,
	}
}

func (au *ActorsUsecase) GetActors(ctx context.Context) ([]models.Actor, error) {
	actors, err := au.repo.ReadActors(ctx)
	if err != nil {
		return make([]models.Actor, 0), err
	}
	return actors, nil
}

func (au *ActorsUsecase) AddActor(ctx context.Context, actor *models.Actor) error {
	err := au.repo.CreateActor(ctx, actor)
	if err != nil {
		return err
	}
	return nil
}

func (au *ActorsUsecase) UpdateActor(ctx context.Context, actor *models.Actor) error {
	a, err := au.repo.ReadActor(ctx, actor.Id)
	if actor.Name != "" {
		a.Name = actor.Name
	}

	if actor.Surname != "" {
		a.Surname = actor.Surname
	}

	if actor.Gender != "" {
		a.Gender = actor.Gender
	}

	if actor.BirthDate != (pgtype.Date{}) {
		a.BirthDate = actor.BirthDate
	}

	err = au.repo.UpdateActor(ctx, a)
	if err != nil {
		return err
	}
	return nil
}

func (au *ActorsUsecase) DeleteActor(ctx context.Context, id int) error {
	err := au.repo.DeleteActor(ctx, id)
	return err
}
