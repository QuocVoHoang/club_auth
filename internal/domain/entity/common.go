package entity

import "github.com/golang-jwt/jwt"

type AuthClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   int    `json:"role"`
	jwt.StandardClaims
}
