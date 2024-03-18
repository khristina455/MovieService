package models

import "github.com/jackc/pgx/v5/pgtype"

type Movie struct {
	Id          int                 `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	ReleaseDate pgtype.Date         `json:"releaseDate"`
	Rating      int                 `json:"rating"`
	Actors      []ActorInMovieSlice `json:"actors"`
}

type MovieInActorSlice struct {
	Id          int         `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	ReleaseDate pgtype.Date `json:"releaseDate"`
	Rating      int         `json:"rating"`
}
