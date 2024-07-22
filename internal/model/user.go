package model

import (
	"context"
	"time"

	pb "github.com/kodinggo/rest-api-service-golang-private-1/pb/story"
)

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

// TODO: Add other interface methods
type UserUsecase interface {
	FindByID(ctx context.Context, id int64) (*User, error)
}

// TODO: Add other interface methods
type UserRepository interface {
	FindByID(ctx context.Context, id int64) (*User, error)
	FindByUsername(ctx context.Context, username string) (result *User, err error)
}

func (u User) ToProto() *pb.User {
	return &pb.User{
		Id:       u.ID,
		Username: u.Username,
	}
}
