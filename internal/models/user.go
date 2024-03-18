package models

import "github.com/golang-jwt/jwt"

type User struct {
	Id       int    `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"isAdmin"`
}

type JwtClaims struct {
	jwt.StandardClaims
	UserId  int  `json:"userId"`
	IsAdmin bool `json:"isAdmin"`
}

type Role int

const (
	Client Role = iota // 0
	Admin              // 1
)
