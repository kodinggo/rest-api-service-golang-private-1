package grpcsvc

import (
	"context"

	pb "github.com/kodinggo/rest-api-service-golang-private-1/pb/story"
)

type UserService struct {
	// default method from story service protobuf
	// if we don't implement method yet, clients are able to call method stub but return error
	pb.UnimplementedUserServiceServer
}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) FindByID(ctx context.Context, in *pb.FindUserByIDRequest) (*pb.User, error) {
	panic("need implementation")
}
