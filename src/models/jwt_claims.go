package models

import "github.com/dgrijalva/jwt-go"

type JWTClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}
