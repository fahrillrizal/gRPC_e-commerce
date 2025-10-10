package entity

import "github.com/golang-jwt/jwt/v5"

type JwtClaims struct {
	UserID   uint   `json:"user_id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	RoleCode string `json:"role_code"`
	jwt.RegisteredClaims
}