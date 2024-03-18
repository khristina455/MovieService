package middleware

import (
	"MovieService/internal/models"
	"MovieService/internal/pkg/utils/jwt"
	resp "MovieService/internal/pkg/utils/responser"
	"log"
	"net/http"
	"strings"
)

const (
	jwtPrefix = "Bearer "
)

func RoleCheck(w http.ResponseWriter, r *http.Request, next func(w http.ResponseWriter, r *http.Request), roles []models.Role) {
	cookie, err := r.Cookie("AccessToken")
	if err != nil && len(roles) != 0 {
		resp.JSONStatus(w, http.StatusForbidden)
		return
	}

	jwtStr := ""
	if cookie != nil {
		jwtStr = cookie.Value
	}

	if !strings.HasPrefix(jwtStr, jwtPrefix) && len(roles) != 0 {
		resp.JSONStatus(w, http.StatusForbidden)
		return
	}

	if len(jwtStr) >= len(jwtPrefix) {
		jwtStr = jwtStr[len(jwtPrefix):]
	}

	userId, isAdmin, err := jwt.TokenManagerSingletone.Parse(jwtStr)

	if len(roles) == 1 {
		if !isAdmin && roles[0] == 1 || isAdmin && roles[0] == 0 {
			resp.JSONStatus(w, http.StatusForbidden)
			log.Printf("user %v is not admin", userId)
			return
		}
	}

	next(w, r)
}
