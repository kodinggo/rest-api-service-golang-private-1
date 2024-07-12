package usecase

import (
	"context"

	"github.com/kodinggo/rest-api-service-golang-private-1/internal/model"
)

type authUsecase struct {
	userRepo model.UserRepository
}

func NewAuthUsecase(userRepo model.UserRepository) model.AuthUsecase {
	return &authUsecase{userRepo: userRepo}
}

func (u *authUsecase) Login(ctx context.Context, username, password string) (token string, err error) {
	// TODO: Handle login
	// 1. Buat method di user repo findByUsername
	// 2. cek password
	// 3. convert ke JWT (user_id)
	return "token", nil
}
