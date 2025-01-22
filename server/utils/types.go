package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JsonUser struct {
	ID             string    `json:"id"`
	Email          string    `json:"email"`
    Username       string    `json:"username"`
    CreatedAt      time.Time `json:"createdAt"`
    UpdatedAt      time.Time `json:"updatedAt"`
}

type MyJWTClaims struct {
	ID             string    `json:"id"`
	Email          string    `json:"email"`
    Username       string    `json:"username"`
    CreatedAt      time.Time `json:"createdAt"`
    UpdatedAt      time.Time `json:"updatedAt"`
	jwt.RegisteredClaims
}