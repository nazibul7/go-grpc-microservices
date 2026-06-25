package model

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UserID    int64
	Email     string
	Role      Role
	TokenType TokenType
	jwt.RegisteredClaims
}