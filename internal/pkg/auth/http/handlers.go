package http

import (
	"MovieService/internal/pkg/auth"
	"log/slog"
	"net/http"
	"regexp"
)

var (
	signinRe = regexp.MustCompile(`^\/auth\/signIn[\/]*$`)
	signupRe = regexp.MustCompile(`^\/auth\/signUp[\/]*$`)
)

type AuthHandler struct {
	log *slog.Logger
	uc  auth.AuthUsecase
}

//Разобраться с jwt и заполнить ручкм для auth

func NewAuthHandler(log *slog.Logger, uc auth.AuthUsecase) AuthHandler {
	return AuthHandler{
		log: log,
		uc:  uc,
	}
}

func (ah *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	switch {
	case r.Method == http.MethodPost && signinRe.MatchString(r.URL.Path):
		ah.SignIn(w, r)
		return
	case r.Method == http.MethodPost && signupRe.MatchString(r.URL.Path):
		ah.SignUp(w, r)
		return
	default:
		return
	}
}

func (ah *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	//ah.uc.SignIn()
}

func (ah *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	// ah.uc.SignUp()
}
