package actors

import (
	"MovieService/internal/models"
	"context"
)

type ActorsRepo interface {
	ReadActors(context.Context) ([]models.Actor, error)
	ReadActor(context.Context, int) (*models.Actor, error)
	CreateActor(context.Context, *models.Actor) error
	UpdateActor(context.Context, *models.Actor) error
	DeleteActor(context.Context, int) error
}

type ActorsUsecase interface {
	GetActors(context.Context) ([]models.Actor, error)
	AddActor(context.Context, *models.Actor) error
	UpdateActor(context.Context, *models.Actor) error
	DeleteActor(context.Context, int) error
}
