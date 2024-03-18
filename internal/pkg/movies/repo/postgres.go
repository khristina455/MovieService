package repo

import (
	"MovieService/internal/models"
	"MovieService/internal/pkg/movies"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	readMovies         = "SELECT id, name, description, release_date, rating FROM movie "
	readeMovie         = "SELECT name, description, release_date, rating FROM movie WHERE id=$1;"
	readMoviesBySearch = "SELECT id, name, description, release_date, rating FROM movie WHERE name LIKE $1"
	createMovie        = "INSERT INTO movie (name, description, release_date, rating) VALUES ($1, $2, $3, $4) RETURNING id;"
	updateMovie        = "UPDATE movie SET name=$1, description=$2, release_date=$3, rating=$4 WHERE id=$5;"
	deleteMovie        = "DELETE FROM movie WHERE id=$1;"
	readActorsOfMovie  = "SELECT a.id, a.name, a.surname, a.gender, a.birth_date FROM actor AS a JOIN movie_actor AS ma ON ma.actor_id = a.id WHERE ma.movie_id=$1"
	createActorMovie   = "INSERT INTO movie_actor (movie_id, actor_id) VALUES ($1, $2)"
)

type MoviesRepo struct {
	db *pgxpool.Pool
}

func NewMoviesRepo(db *pgxpool.Pool) *MoviesRepo {
	return &MoviesRepo{
		db: db,
	}
}

func (mr *MoviesRepo) ReadMovies(ctx context.Context, sortType string) ([]models.Movie, error) {
	var endExpr string
	switch sortType {
	case movies.NAME_ASC:
		endExpr = "ORDER BY name;"
		break
	case movies.NAME_DESC:
		endExpr = "ORDER BY name DESC;"
		break
	case movies.DATE_ASC:
		endExpr = "ORDER BY release_date;"
		break
	case movies.DATE_DESC:
		endExpr = "ORDER BY release_date DESC;"
		break
	case movies.RATING_ASC:
		endExpr = "ORDER BY rating;"
		break
	case movies.RATING_DESC:
		endExpr = "ORDER BY rating DESC;"
		break
	default:
		return make([]models.Movie, 0), nil
	}

	movieSlice := make([]models.Movie, 0)
	rows, err := mr.db.Query(ctx, readMovies+endExpr)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []models.Movie{}, err
		}
		err = fmt.Errorf("error happened in db.QueryContext: %w", err)

		return []models.Movie{}, err
	}
	movie := models.Movie{}
	for rows.Next() {
		err = rows.Scan(
			&movie.Id,
			&movie.Name,
			&movie.Description,
			&movie.ReleaseDate,
			&movie.Rating,
		)
		if err != nil {
			err = fmt.Errorf("error happened in rows.Scan: %w", err)

			return []models.Movie{}, err
		}

		actorRows, err := mr.db.Query(ctx, readActorsOfMovie, movie.Id)
		if err != nil {
			err = fmt.Errorf("error happened in rows.Scan: %w", err)

			return []models.Movie{}, err
		}

		actorSlice := make([]models.ActorInMovieSlice, 0)
		actor := models.ActorInMovieSlice{}
		for actorRows.Next() {
			err = actorRows.Scan(
				&actor.Id,
				&actor.Name,
				&actor.Surname,
				&actor.Gender,
				&actor.BirthDate,
			)

			if err != nil {
				err = fmt.Errorf("error happened in rows.Scan: %w", err)
				return []models.Movie{}, err
			}

			actorSlice = append(actorSlice, actor)
		}
		movie.Actors = actorSlice

		movieSlice = append(movieSlice, movie)
	}
	defer rows.Close()

	return movieSlice, nil
}

func (mr *MoviesRepo) ReadMovie(ctx context.Context, id int) (*models.Movie, error) {
	m := &models.Movie{Id: id}
	if err := mr.db.QueryRow(ctx, readeMovie, id).
		Scan(&m.Name, &m.Description, &m.ReleaseDate, &m.Rating); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &models.Movie{}, err
		}
		err = fmt.Errorf("error happened in row.Scan: %w", err)

		return &models.Movie{}, err
	}
	return m, nil
}

func (mr *MoviesRepo) CreateMovie(ctx context.Context, movie *models.Movie) (int, error) {
	var id int
	err := mr.db.QueryRow(ctx, createMovie,
		movie.Name, movie.Description, movie.ReleaseDate, movie.Rating).Scan(&id)

	if err != nil {
		err = fmt.Errorf("error happened in scan.Scan: %w", err)

		return 0, err
	}

	return id, nil
}

func (mr *MoviesRepo) UpdateMovie(ctx context.Context, movie *models.Movie) error {
	_, err := mr.db.Exec(ctx, updateMovie, movie.Name, movie.Description, movie.ReleaseDate, movie.Rating, movie.Id)
	if err != nil {
		err = fmt.Errorf("error happened in db.Exec: %w", err)

		return err
	}

	return nil
}

func (mr *MoviesRepo) DeleteMovie(ctx context.Context, id int) error {
	_, err := mr.db.Exec(ctx, deleteMovie, id)
	if err != nil {
		err = fmt.Errorf("error happened in db.Exec: %w", err)

		return err
	}

	return nil
}

func (mr *MoviesRepo) ReadMoviesBySearch(ctx context.Context, movieName string, actorName string) ([]models.Movie, error) {
	movieSlice := make([]models.Movie, 0)
	movieName = "%" + movieName + "%"
	rows, err := mr.db.Query(ctx, readMoviesBySearch, movieName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []models.Movie{}, err
		}
		err = fmt.Errorf("error happened in db.QueryContext: %w", err)

		return []models.Movie{}, err
	}
	movie := models.Movie{}
	for rows.Next() {
		err = rows.Scan(
			&movie.Id,
			&movie.Name,
			&movie.Description,
			&movie.ReleaseDate,
			&movie.Rating,
		)
		if err != nil {
			err = fmt.Errorf("error happened in rows.Scan: %w", err)

			return []models.Movie{}, err
		}
		movieSlice = append(movieSlice, movie)
	}
	defer rows.Close()

	return movieSlice, nil
}

func (mr *MoviesRepo) AddActorToMovie(ctx context.Context, movieId int, actorId int) error {
	_, err := mr.db.Exec(ctx, createActorMovie, movieId, actorId)

	if err != nil {
		err = fmt.Errorf("error happened in db.Exec: %w", err)

		return err
	}

	return nil
}
