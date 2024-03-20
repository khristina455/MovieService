package http

import (
	"MovieService/internal/models"
	"MovieService/internal/pkg/auth"
	"MovieService/internal/pkg/utils/jwt"
	resp "MovieService/internal/pkg/utils/responser"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"regexp"
)

var (
	signinRe = regexp.MustCompile(`^\/api\/auth\/signIn[\/]*$`)
	signupRe = regexp.MustCompile(`^\/api\/auth\/signUp[\/]*$`)
)

type AuthHandler struct {
	log *slog.Logger
	uc  auth.AuthUsecase
}

func NewAuthHandler(log *slog.Logger, uc auth.AuthUsecase) AuthHandler {
	return AuthHandler{
		log: log,
		uc:  uc,
	}
}

func (ah *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	fmt.Println(r.URL.Path)
	switch {
	case r.Method == http.MethodPost && signinRe.MatchString(r.URL.Path):
		ah.SignIn(w, r)
		return
	case r.Method == http.MethodPost && signupRe.MatchString(r.URL.Path):
		ah.SignUp(w, r)
		return
	default:
		resp.JSONStatus(w, http.StatusNotFound)
	}
}

// SignIn godoc
// @Summary      User sign-in
// @Description  Authenticates a user and generates an access token
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        user  body  models.User  true  "User information"
// @Success      200
// @Failure      400
// @Failure      500
// @Router       /api/signIn [post]
func (ah *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		resp.JSONStatus(w, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	u := &models.User{}
	err = json.Unmarshal(body, u)
	if err != nil {
		resp.JSONStatus(w, http.StatusInternalServerError)
		return
	}

	err = ah.uc.SignIn(r.Context(), u)
	if err != nil {
		fmt.Println(err)
		resp.JSONStatus(w, http.StatusInternalServerError)
		return
	}

	token, err := jwt.TokenManagerSingletone.NewJWT(u.Id, u.IsAdmin)
	if err != nil {
		resp.JSONStatus(w, http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{Name: "AccessToken", Value: "Bearer " + token})
	resp.JSONStatus(w, http.StatusOK)
}

// SignUp godoc
// @Summary      Sign up a new user
// @Description  Creates a new user account
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        user  body  models.User  true  "User information"
// @Success      200
// @Failure      400
// @Failure      500
// @Router       /api/signUp [post]
func (ah *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	fmt.Println("signup")
	body, err := io.ReadAll(r.Body)

	if err != nil {
		resp.JSONStatus(w, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	u := &models.User{}
	err = json.Unmarshal(body, u)
	if err != nil {
		resp.JSONStatus(w, http.StatusInternalServerError)
		return
	}

	u.Id, err = ah.uc.SignUp(r.Context(), u)
	if err != nil {
		fmt.Println(err)
		resp.JSONStatus(w, http.StatusInternalServerError)
		return
	}

	token, err := jwt.TokenManagerSingletone.NewJWT(u.Id, u.IsAdmin)
	if err != nil {
		resp.JSONStatus(w, http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{Name: "AccessToken", Value: "Bearer " + token})
	resp.JSONStatus(w, http.StatusOK)
}
