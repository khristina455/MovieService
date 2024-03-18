package movies

import (
	"MovieService/internal/models"
	"context"
)

const (
	NAME_ASC    = "name_asc"
	NAME_DESC   = "name_desc"
	RATING_DESC = "rating_desc"
	RATING_ASC  = "rating_asc"
	DATE_ASC    = "date_asc"
	DATE_DESC   = "date_desc"
)

type MoviesRepo interface {
	ReadMovies(context.Context, string) ([]models.Movie, error)
	ReadMovie(context.Context, int) (*models.Movie, error)
	CreateMovie(context.Context, *models.Movie) (int, error)
	UpdateMovie(context.Context, *models.Movie) error
	DeleteMovie(context.Context, int) error
	ReadMoviesBySearch(context.Context, string, string) ([]models.Movie, error)
	AddActorToMovie(context.Context, int, int) error
}

type MoviesUsecase interface {
	GetMovies(context.Context, string) ([]models.Movie, error)
	AddMovie(context.Context, *models.Movie) error
	UpdateMovie(context.Context, *models.Movie) error
	DeleteMovie(context.Context, int) error
	GetMoviesBySearch(context.Context, string, string) ([]models.Movie, error)
}
