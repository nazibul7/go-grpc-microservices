package model

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UserID    string
	Email     string
	Role      Role
	jwt.RegisteredClaims
}
