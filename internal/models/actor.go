package models

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Actor struct {
	Id        int                 `json:"id"`
	Name      string              `json:"name"`
	Surname   string              `json:"surname"`
	Gender    string              `json:"gender"`
	BirthDate pgtype.Date         `json:"birthDate"`
	Movies    []MovieInActorSlice `json:"movies"`
}

type ActorInMovieSlice struct {
	Id        int         `json:"id"`
	Name      string      `json:"name"`
	Surname   string      `json:"surname"`
	Gender    string      `json:"gender"`
	BirthDate pgtype.Date `json:"birthDate"`
}
