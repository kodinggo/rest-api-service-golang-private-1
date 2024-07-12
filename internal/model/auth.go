package model

import "context"

type AuthKey string

const LoginKey AuthKey = "loginkey"

type AuthUsecase interface {
	Login(ctx context.Context, username, password string) (token string, err error)
}
