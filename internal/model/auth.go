package model

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
)

type AuthKey string

const LoginKey AuthKey = "loginkey"
const JWTKey AuthKey = "jwtkey"

type CustomClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

type AuthUsecase interface {
	Login(ctx context.Context, username, password string) (token string, err error)
}
