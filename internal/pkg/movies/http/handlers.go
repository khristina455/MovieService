package http

import (
	"MovieService/internal/models"
	"MovieService/internal/pkg/movies"
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
	allMoviesRe      = regexp.MustCompile(`^\/api\/movies((\?.*)|(\/*))$`)
	moviesBySearchRe = regexp.MustCompile(`^\/api\/movies\/search\?.*$`)
	addMovieRe       = regexp.MustCompile(`^\/api\/movies[\/]*$`)
	updateMovieRe    = regexp.MustCompile(`^\/api\/movies\/(\d+)$`)
	deleteMovieRe    = regexp.MustCompile(`^\/api\/movies\/(\d+)$`)
)

type MoviesHandler struct {
	log *slog.Logger
	uc  movies.MoviesUsecase
}

func NewMoviesHandler(log *slog.Logger, uc movies.MoviesUsecase) MoviesHandler {
	return MoviesHandler{
		log: log,
		uc:  uc,
	}
}

// TODO разделить поиск по названию и актеру
func (mh *MoviesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	fmt.Println(r.URL.RequestURI())
	switch {
	case r.Method == http.MethodGet && allMoviesRe.MatchString(r.URL.RequestURI()):
		mh.GetMovies(w, r)
		return
	case r.Method == http.MethodGet && moviesBySearchRe.MatchString(r.URL.RequestURI()):
		mh.GetMoviesBySearch(w, r)
		return
	case r.Method == http.MethodPost && addMovieRe.MatchString(r.URL.RequestURI()):
		mh.AddMovie(w, r)
		return
	case r.Method == http.MethodPut && updateMovieRe.MatchString(r.URL.RequestURI()):
		mh.UpdateMovie(w, r)
		return
	case r.Method == http.MethodDelete && deleteMovieRe.MatchString(r.URL.RequestURI()):
		mh.DeleteMovie(w, r)
		return
	default:
		return
	}
}

func (mh *MoviesHandler) GetMovies(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get movies")
	sort := r.URL.Query().Get("sorting")

	if sort == "" {
		sort = movies.RATING_DESC
	}

	movies, err := mh.uc.GetMovies(r.Context(), sort)
	if err != nil {
		fmt.Println(err)
		resp.JSONStatus(w, http.StatusInternalServerError)
		return
	}

	resp.JSON(w, http.StatusOK, movies)
}

func (mh *MoviesHandler) GetMoviesBySearch(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get search movies")
	movieName := r.URL.Query().Get("movie_name")
	actorName := r.URL.Query().Get("actor_name")

	if movieName != "" && actorName == "" {
		movies, err := mh.uc.GetMoviesByMovieName(r.Context(), movieName)
		if err != nil {
			resp.JSONStatus(w, http.StatusInternalServerError)
			return
		}
		resp.JSON(w, http.StatusOK, movies)
	} else if movieName == "" && actorName != "" {
		movies, err := mh.uc.GetMoviesByActorName(r.Context(), actorName)
		fmt.Println(err)
		if err != nil {
			resp.JSONStatus(w, http.StatusInternalServerError)
			return
		}
		resp.JSON(w, http.StatusOK, movies)
	} else {
		resp.JSONStatus(w, http.StatusBadRequest)
	}
}

func (mh *MoviesHandler) AddMovie(w http.ResponseWriter, r *http.Request) {
	fmt.Println("add movie")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}
	defer r.Body.Close()

	m := &models.Movie{}
	err = json.Unmarshal(body, m)

	err = mh.uc.AddMovie(r.Context(), m)
	if err != nil {
		fmt.Println(err)
		resp.JSONStatus(w, http.StatusInternalServerError)
		return
	}

	resp.JSONStatus(w, http.StatusOK)
}

func (mh *MoviesHandler) UpdateMovie(w http.ResponseWriter, r *http.Request) {
	fmt.Println("update movies")
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

	m := &models.Movie{}
	err = json.Unmarshal(body, m)
	m.Id = id

	err = mh.uc.UpdateMovie(r.Context(), m)

	if err != nil {
		resp.JSONStatus(w, http.StatusInternalServerError)
		return
	}

	resp.JSONStatus(w, http.StatusOK)
}

func (mh *MoviesHandler) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	fmt.Println("delete movies")
	idStr := filepath.Base(r.URL.Path)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		resp.JSONStatus(w, http.StatusBadRequest)
		return
	}

	err = mh.uc.DeleteMovie(r.Context(), id)

	if err != nil {
		resp.JSONStatus(w, http.StatusInternalServerError)
		return
	}

	resp.JSONStatus(w, http.StatusOK)
}
