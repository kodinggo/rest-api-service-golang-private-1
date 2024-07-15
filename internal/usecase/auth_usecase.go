package usecase

import (
	"context"
	"errors"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/kodinggo/rest-api-service-golang-private-1/internal/config"
	"github.com/kodinggo/rest-api-service-golang-private-1/internal/model"
	"github.com/labstack/gommon/log"
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
	user, err := u.userRepo.FindByUsername(ctx, username)
	if err != nil {
		log.Errorf("failed find user by username, error: %v", err)
		err = errors.New("invalid usernamae")
		return
	}

	// 2. cek password
	// TODO: Implement password encryption
	if user.Password != password {
		err = errors.New("invalid password")
		return
	}

	// 3. convert ke JWT (user_id)
	timeNowUTC := time.Now().UTC()
	expiredTime := timeNowUTC.Add(config.JWTExp()) // expired 24 jam
	jwtToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     expiredTime.Unix(),
	}).SignedString([]byte(config.JWTSigningKey()))
	if err != nil {
		log.Errorf("failed when generate jwt token, error: %v", err)
		err = errors.New("failed when generate jwt token")
		return
	}

	return jwtToken, nil
}
