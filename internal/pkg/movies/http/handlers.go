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
	"strings"
)

var (
	allMoviesRe            = regexp.MustCompile(`^\/api\/movies((\?.*)|(\/*))$`)
	moviesBySearchRe       = regexp.MustCompile(`^\/api\/movies\/search\?.*$`)
	addMovieRe             = regexp.MustCompile(`^\/api\/movies[\/]*$`)
	updateMovieRe          = regexp.MustCompile(`^\/api\/movies\/(\d+)$`)
	deleteMovieRe          = regexp.MustCompile(`^\/api\/movies\/(\d+)$`)
	deleteActorFromMovieRe = regexp.MustCompile(`^\/api\/movies\/(\d+)\/actors[\/]*$`)
	addActorToMovieRe      = regexp.MustCompile(`^\/api\/movies\/(\d+)\/actors[\/]*$`)
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
	case r.Method == http.MethodPost && addActorToMovieRe.MatchString(r.URL.Path):
		mh.AddActorToMovie(w, r)
		return
	case r.Method == http.MethodDelete && deleteActorFromMovieRe.MatchString(r.URL.Path):
		mh.DeleteActorFromMovie(w, r)
	default:
		resp.JSONStatus(w, http.StatusNotFound)
	}
}

// GetMovies godoc
// @Summary      Get list of movies
// @Description  Retrieves a list of movies based on the provided parameters
// @Tags         Movies
// @Produce      json
// @Param        sorting   query    string  false  "Query string to sort movies"
// @Success      200  {array}  models.Movie
// @Failure      500
// @Router       /api/movies [get]
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

// GetMoviesBySearch godoc
// @Summary      Get list of movies
// @Description  Retrieves a list of movies based on the provided parameters
// @Tags         Movies
// @Produce      json
// @Param        movie_name   query    string  false  "Name of movie to filter movies"
// @Param        actor_name   query    string  false  "Name of actor to filter movies"
// @Success      200  {array}  models.Movie
// @Failure      400
// @Failure      500
// @Router       /api/movies/search [get]
func (mh *MoviesHandler) GetMoviesBySearch(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get search movies")
	movieName := r.URL.Query().Get("movie_name")
	actorName := r.URL.Query().Get("actor_name")

	if movieName != "" && actorName != "" {
		resp.JSONStatus(w, http.StatusBadRequest)
	} else if movieName == "" && actorName != "" {
		movies, err := mh.uc.GetMoviesByActorName(r.Context(), actorName)
		fmt.Println(err)
		if err != nil {
			resp.JSONStatus(w, http.StatusInternalServerError)
			return
		}
		resp.JSON(w, http.StatusOK, movies)
	} else {
		movies, err := mh.uc.GetMoviesByMovieName(r.Context(), movieName)
		if err != nil {
			resp.JSONStatus(w, http.StatusInternalServerError)
			return
		}
		resp.JSON(w, http.StatusOK, movies)
	}
}

// AddMovie godoc
// @Summary      Add a new movie
// @Description  Add a new movie with name, description, release date, rating
// @Tags         Movies
// @Accept       json
// @Param        movie  body  models.Movie  true  "Movie information"
// @Success      200
// @Failure      400
// @Failure      500
// @Router       /api/movies [post]
func (mh *MoviesHandler) AddMovie(w http.ResponseWriter, r *http.Request) {
	fmt.Println("add movie")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}
	defer r.Body.Close()

	m := &models.Movie{}
	err = json.Unmarshal(body, m)
	if err != nil {
		resp.JSONStatus(w, http.StatusBadRequest)
	}

	err = mh.uc.AddMovie(r.Context(), m)
	if err != nil {
		fmt.Println(err)
		resp.JSONStatus(w, http.StatusInternalServerError)
		return
	}

	resp.JSONStatus(w, http.StatusOK)
}

// UpdateMovie godoc
// @Summary      Update movie by ID
// @Description  Updates a movie with the given ID
// @Tags         Movies
// @Accept       json
// @Param        id  path  int  true  "Movie ID"
// @Param        movie  body  models.Movie  true  "Movie information to update"
// @Success      200
// @Failure      400
// @Failure      500
// @Router       /api/movie/{id} [put]
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

// DeleteMovie godoc
// @Summary      Delete movie by ID
// @Description  Deletes a movie with the given ID
// @Tags         Movies
// @Accept       json
// @Param        id  path  int  true  "Movie ID"
// @Success      200
// @Failure      400
// @Failure      500
// @Router       /api/movies/{id} [delete]
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

//Наложить ограничение уникальности на ключи в м-м

// AddActorToMovie godoc
// @Summary      Add an actor to movie
// @Description  Add an actor to movie by their ids
// @Tags         Movies
// @Accept       json
// @Param        id  path  int  true  "Movie ID"
// @Param        id  body  int  true  "Actor id"
// @Success      200
// @Failure      400
// @Failure      500
// @Router       /api/movies/{id}/actors [post]
func (mh *MoviesHandler) AddActorToMovie(w http.ResponseWriter, r *http.Request) {
	sliceOfURL := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(sliceOfURL[1])
	if err != nil {
		resp.JSONStatus(w, http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		resp.JSONStatus(w, http.StatusBadRequest)
		return
	}

	m := make(map[string]interface{})
	err = json.Unmarshal(body, &m)
	if err != nil {
		resp.JSONStatus(w, http.StatusBadRequest)
		return
	}

	actorId := m["id"].(int)
	err = mh.uc.AddActorToMovie(r.Context(), id, actorId)
	if err != nil {
		resp.JSONStatus(w, http.StatusInternalServerError)
		return
	}

	resp.JSONStatus(w, http.StatusOK)
}

// DeleteActorFromMovie godoc
// @Summary      Delete actor from movie
// @Description  Delete actor from movie by their ids
// @Tags         Movies
// @Accept       json
// @Param        movieId  path  int  true  "Movie ID"
// @Param        actorId  path  int  true  "Actor id"
// @Success      200
// @Failure      400
// @Failure      500
// @Router       /api/movies/{movieId}/actors/{actorId} [delete]
func (mh *MoviesHandler) DeleteActorFromMovie(w http.ResponseWriter, r *http.Request) {
	sliceOfURL := strings.Split(r.URL.Path, "/")
	movieId, err := strconv.Atoi(sliceOfURL[2])
	if err != nil {
		resp.JSONStatus(w, http.StatusBadRequest)
		return
	}

	actorId, err := strconv.Atoi(sliceOfURL[4])
	if err != nil {
		resp.JSONStatus(w, http.StatusBadRequest)
		return
	}

	err = mh.uc.DeleteActorFromMovie(r.Context(), movieId, actorId)
	if err != nil {
		resp.JSONStatus(w, http.StatusInternalServerError)
		return
	}

	resp.JSONStatus(w, http.StatusOK)
}
