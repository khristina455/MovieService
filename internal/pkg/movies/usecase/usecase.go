package usecase

import (
	"MovieService/internal/models"
	"MovieService/internal/pkg/movies"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
)

type MoviesUsecase struct {
	repo movies.MoviesRepo
}

func NewMoviesUsecase(repo movies.MoviesRepo) *MoviesUsecase {
	return &MoviesUsecase{
		repo: repo,
	}
}

func (mu MoviesUsecase) GetMovies(ctx context.Context, sortType string) ([]models.Movie, error) {
	m, err := mu.repo.ReadMovies(ctx, sortType)
	return m, err
}

func (mu MoviesUsecase) AddMovie(ctx context.Context, movie *models.Movie) error {
	movieId, err := mu.repo.CreateMovie(ctx, movie)
	if err != nil {
		return err
	}

	fmt.Println(movieId)
	for _, actor := range movie.Actors {
		err = mu.repo.AddActorToMovie(ctx, movieId, actor.Id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (mu MoviesUsecase) UpdateMovie(ctx context.Context, movie *models.Movie) error {
	m, err := mu.repo.ReadMovie(ctx, movie.Id)
	if movie.Name != "" {
		m.Name = movie.Name
	}

	if movie.Description != "" {
		m.Description = movie.Description
	}

	if movie.Rating != 0 {
		m.Rating = movie.Rating
	}

	if movie.ReleaseDate != (pgtype.Date{}) {
		m.ReleaseDate = movie.ReleaseDate
	}

	err = mu.repo.UpdateMovie(ctx, m)
	if err != nil {
		return err
	}
	return nil
}

func (mu MoviesUsecase) DeleteMovie(ctx context.Context, id int) error {
	err := mu.repo.DeleteMovie(ctx, id)
	return err
}

func (mu MoviesUsecase) GetMoviesByMovieName(ctx context.Context, s string) ([]models.Movie, error) {
	m, err := mu.repo.ReadMoviesByMovieName(ctx, s)
	return m, err
}

func (mu MoviesUsecase) GetMoviesByActorName(ctx context.Context, s string) ([]models.Movie, error) {
	m, err := mu.repo.ReadMoviesByActorName(ctx, s)
	return m, err
}

func (mu MoviesUsecase) AddActorToMovie(ctx context.Context, movieId int, actorId int) error {
	err := mu.repo.AddActorToMovie(ctx, movieId, actorId)
	return err
}

func (mu MoviesUsecase) DeleteActorFromMovie(ctx context.Context, movieId int, actorId int) error {
	err := mu.repo.DeleteActorFromMovie(ctx, movieId, actorId)
	return err
}
