package utils

import (
	"backend/models"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"strings"
	"time"
)

var signKey = []byte(os.Getenv("JWT_SECRET_KEY"))

var AuthJWTDuration = time.Hour
var RefreshJWTDuration = time.Hour * 24 * 30

// CreateToken creates a JWT from an email
func CreateToken(email string, lifetime time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(lifetime).Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(signKey)

	if err != nil {
		return "", fmt.Errorf("error signing token: %v", err)
	}

	return tokenString, nil
}

// ExtractTokenFromHeader extracts the auth token from the request header
func ExtractTokenFromHeader(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")

	if header == "" {
		return "", fmt.Errorf("missing Authorization header")
	}

	args := strings.Split(header, " ")

	if len(args) != 2 || args[0] != "Bearer" {
		return "", fmt.Errorf("invalid Authorization header")
	}

	return args[1], nil
}

// ValidateToken validates that a JWT is valid
func ValidateToken(tokenStr string) (*models.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return signKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	if claims, ok := token.Claims.(*models.JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
