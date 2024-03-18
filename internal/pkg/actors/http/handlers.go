package http

import (
	"MovieService/internal/models"
	"MovieService/internal/pkg/actors"
	"MovieService/internal/pkg/middleware"
	resp "MovieService/internal/pkg/utils/responser"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"
)

var (
	allActorsRe   = regexp.MustCompile(`^\/api\/actors[\/]*$`)
	addActorRe    = regexp.MustCompile(`^\/api\/actors[\/]*$`)
	updateActorRe = regexp.MustCompile(`^\/api\/actors\/([0-9]*)$`)
	deleteActorRe = regexp.MustCompile(`^\/api\/actors\/(\d+)$`)
)

type ActorsHandler struct {
	log *slog.Logger
	uc  actors.ActorsUsecase
}

func NewActorsHandler(log *slog.Logger, uc actors.ActorsUsecase) ActorsHandler {
	return ActorsHandler{
		log: log,
		uc:  uc,
	}
}

// TODO Добавить ручку для добавления фильма
func (ah *ActorsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	fmt.Println(r.URL.Path)
	switch {
	case r.Method == http.MethodGet && allActorsRe.MatchString(r.URL.Path):
		middleware.RoleCheck(w, r, ah.GetActors, []models.Role{models.Client})
		return
	case r.Method == http.MethodPost && addActorRe.MatchString(r.URL.Path):
		middleware.RoleCheck(w, r, ah.AddActor, []models.Role{models.Admin})
		return
	case r.Method == http.MethodPut && updateActorRe.MatchString(r.URL.Path):
		middleware.RoleCheck(w, r, ah.UpdateActor, []models.Role{models.Admin})
		return
	case r.Method == http.MethodDelete && deleteActorRe.MatchString(r.URL.Path):
		middleware.RoleCheck(w, r, ah.DeleteActor, []models.Role{models.Admin})
		return
	default:
		return
	}
}

func (ah *ActorsHandler) GetActors(w http.ResponseWriter, r *http.Request) {
	actors, err := ah.uc.GetActors(r.Context())
	fmt.Println("get actors")
	if err != nil {
		fmt.Println(err)
		resp.JSONStatus(w, http.StatusInternalServerError)
		return
	}

	resp.JSON(w, http.StatusOK, actors)
}

// TODO:добавить обработку пустых полей
func (ah *ActorsHandler) AddActor(w http.ResponseWriter, r *http.Request) {
	fmt.Println("add actor")
	body, err := io.ReadAll(r.Body)

	if err != nil {
		resp.JSONStatus(w, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	a := &models.Actor{}
	err = json.Unmarshal(body, a)

	err = ah.uc.AddActor(r.Context(), a)
	if err != nil {
		resp.JSONStatus(w, http.StatusInternalServerError)
		return
	}

	resp.JSONStatus(w, http.StatusOK)
}

func (ah *ActorsHandler) UpdateActor(w http.ResponseWriter, r *http.Request) {
	fmt.Println("update actors")
	idStr := filepath.Base(r.URL.Path)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		resp.JSONStatus(w, http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}
	defer r.Body.Close()

	a := &models.Actor{}
	err = json.Unmarshal(body, a)
	a.Id = id
	fmt.Println(id)

	err = ah.uc.UpdateActor(r.Context(), a)

	if err != nil {
		fmt.Println(err)
		resp.JSONStatus(w, http.StatusInternalServerError)
		return
	}

	resp.JSONStatus(w, http.StatusOK)
}

func (ah *ActorsHandler) DeleteActor(w http.ResponseWriter, r *http.Request) {
	fmt.Println("delete actors")
	idStr := filepath.Base(r.URL.Path)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		resp.JSONStatus(w, http.StatusBadRequest)
		return
	}

	err = ah.uc.DeleteActor(r.Context(), id)

	if err != nil {
		resp.JSONStatus(w, http.StatusInternalServerError)
		return
	}

	resp.JSONStatus(w, http.StatusOK)
}
