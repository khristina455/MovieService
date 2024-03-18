package jwt

import (
	"MovieService/internal/models"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"math/rand"
	"time"
)

type TokenManager interface {
	NewJWT(userId int, isAdmin bool) (string, error)
	Parse(accessToken string) (int, bool, error)
	NewRefreshToken() (string, error)
}

type Manager struct {
	signingKey string
}

var TokenManagerSingletone *Manager

func LoadSecret(signingKey string) error {
	if signingKey == "" {
		return errors.New("empty signing key")
	}

	TokenManagerSingletone = &Manager{signingKey: signingKey}
	return nil
}

func (m *Manager) NewJWT(userId int, isAdmin bool) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &models.JwtClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserId:  userId,
		IsAdmin: isAdmin,
	})

	return token.SignedString([]byte(m.signingKey))
}

func (m *Manager) Parse(accessToken string) (int, bool, error) {
	token, err := jwt.ParseWithClaims(accessToken, &models.JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.signingKey), nil
	})

	if err != nil {
		return 0, false, err
	}

	myClaims := token.Claims.(*models.JwtClaims)
	return myClaims.UserId, myClaims.IsAdmin, nil
}

func (m *Manager) NewRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
