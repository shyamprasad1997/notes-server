package models

import "github.com/golang-jwt/jwt/v4"

type User struct {
	Id       int32
	Name     string
	Email    string
	Password string
}

type Claims struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	jwt.RegisteredClaims
}
