package grpcsvc

import (
	"context"

	pb "github.com/kodinggo/rest-api-service-golang-private-1/pb/story"
	"google.golang.org/protobuf/types/known/emptypb"
)

type StoryService struct {
	// default method from story service protobuf
	// if we don't implement method yet, clients are able to call method stub but return error
	pb.UnimplementedStoryServiceServer
}

func NewStoryService() *StoryService {
	return &StoryService{}
}

func (s *StoryService) FindAll(ctx context.Context, in *pb.FindAllStoriesRequest) (*pb.Stories, error) {
	panic("need implementation")
}

func (s *StoryService) FindByID(ctx context.Context, in *pb.FindStoryByIDRequest) (*pb.Story, error) {
	panic("need implementation")
}

func (s *StoryService) Create(ctx context.Context, in *pb.CreateStoryRequest) (*pb.Story, error) {
	panic("need implementation")
}

func (s *StoryService) Update(ctx context.Context, in *pb.UpdateStoryRequest) (*pb.Story, error) {
	panic("need implementation")
}

func (s *StoryService) Delete(ctx context.Context, in *pb.DeleteStoryRequest) (*emptypb.Empty, error) {
	panic("need implementation")
}
