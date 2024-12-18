package models

import (
	"github.com/golang-jwt/jwt/v4"
)

type TokenClaims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}
