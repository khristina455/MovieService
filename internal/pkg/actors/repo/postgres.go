package repo

import (
	"MovieService/internal/models"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	readActors        = "SELECT id, name, surname, gender, birth_date FROM actor;"
	readActor         = "SELECT name, surname, gender, birth_date FROM actor WHERE id=$1;"
	createActor       = "INSERT INTO actor (name, surname, gender, birth_date) VALUES ($1, $2, $3, $4);"
	updateActor       = "UPDATE actor SET name=$1, surname=$2, gender=$3, birth_date=$4 WHERE id=$5;"
	deleteActor       = "DELETE FROM actor WHERE id=$1;"
	readMoviesOfActor = "SELECT m.id, m.name, m.description, m.release_date, m.rating FROM movie AS m JOIN movie_actor AS ma ON ma.movie_id = m.id WHERE ma.actor_id=$1"
)

type ActorsRepo struct {
	db *pgxpool.Pool
}

func NewActorsRepo(db *pgxpool.Pool) *ActorsRepo {
	return &ActorsRepo{
		db: db,
	}
}

func (ar *ActorsRepo) ReadActors(ctx context.Context) ([]models.Actor, error) {
	actorSlice := make([]models.Actor, 0)
	rows, err := ar.db.Query(ctx, readActors)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []models.Actor{}, err
		}
		err = fmt.Errorf("error happened in db.QueryContext: %w", err)

		return []models.Actor{}, err
	}
	actor := models.Actor{}
	for rows.Next() {
		err = rows.Scan(
			&actor.Id,
			&actor.Name,
			&actor.Surname,
			&actor.Gender,
			&actor.BirthDate,
		)
		if err != nil {
			err = fmt.Errorf("error happened in rows.Scan: %w", err)

			return []models.Actor{}, err
		}

		movieRows, err := ar.db.Query(ctx, readMoviesOfActor, actor.Id)
		if err != nil {
			err = fmt.Errorf("error happened in rows.Scan: %w", err)

			return []models.Actor{}, err
		}

		movieSlice := make([]models.MovieInActorSlice, 0)
		movie := models.MovieInActorSlice{}
		for movieRows.Next() {
			err = movieRows.Scan(
				&movie.Id,
				&movie.Name,
				&movie.Description,
				&movie.ReleaseDate,
				&movie.Rating,
			)

			if err != nil {
				err = fmt.Errorf("error happened in rows.Scan: %w", err)
				return []models.Actor{}, err
			}

			movieSlice = append(movieSlice, movie)
		}
		actor.Movies = movieSlice

		actorSlice = append(actorSlice, actor)
	}
	defer rows.Close()

	return actorSlice, nil
}

func (ar *ActorsRepo) ReadActor(ctx context.Context, id int) (*models.Actor, error) {
	a := &models.Actor{Id: id}
	if err := ar.db.QueryRow(ctx, readActor, id).
		Scan(&a.Name, &a.Surname, &a.Gender, &a.BirthDate); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &models.Actor{}, err
		}
		err = fmt.Errorf("error happened in row.Scan: %w", err)

		return &models.Actor{}, err
	}
	return a, nil
}

func (ar *ActorsRepo) CreateActor(ctx context.Context, actor *models.Actor) error {
	_, err := ar.db.Exec(ctx, createActor,
		actor.Name, actor.Surname, actor.Gender, actor.BirthDate)

	if err != nil {
		err = fmt.Errorf("error happened in db.Exec: %w", err)

		return err
	}

	return nil
}

func (ar *ActorsRepo) UpdateActor(ctx context.Context, actor *models.Actor) error {
	_, err := ar.db.Exec(ctx, updateActor, actor.Name, actor.Surname, actor.Gender, actor.BirthDate, actor.Id)
	if err != nil {
		err = fmt.Errorf("error happened in db.Exec: %w", err)

		return err
	}

	return nil
}

func (ar *ActorsRepo) DeleteActor(ctx context.Context, id int) error {
	_, err := ar.db.Exec(ctx, deleteActor, id)
	if err != nil {
		err = fmt.Errorf("error happened in db.Exec: %w", err)

		return err
	}

	return nil
}
