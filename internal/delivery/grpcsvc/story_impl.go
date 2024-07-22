package grpcsvc

import (
	"context"

	"github.com/kodinggo/rest-api-service-golang-private-1/internal/model"
	pb "github.com/kodinggo/rest-api-service-golang-private-1/pb/story"
	"github.com/labstack/gommon/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

type StoryService struct {
	// default method from story service protobuf
	// if we don't implement method yet, clients are able to call method stub but return error
	pb.UnimplementedStoryServiceServer

	storyUsecase model.StoryUsecase
}

func NewStoryService(storyUsecase model.StoryUsecase) *StoryService {
	return &StoryService{storyUsecase: storyUsecase}
}

func (s *StoryService) FindAll(ctx context.Context, in *pb.FindAllStoriesRequest) (*pb.Stories, error) {
	stories, _, err := s.storyUsecase.FindAll(ctx, &model.StoryOptions{
		Search: in.Search,
		SortBy: in.SortBy,
		// Cursor: in.Cursor,
	})
	if err != nil {
		log.Errorf("failed find all stories, error: %v", err)
		return nil, err
	}

	// Konversi dari main entity ke protobuf
	var pbStories []*pb.Story
	for _, story := range stories {
		pbStories = append(pbStories, story.ToProto())
	}

	return &pb.Stories{
		Stories: pbStories,
	}, nil
}

func (s *StoryService) FindByID(ctx context.Context, in *pb.FindStoryByIDRequest) (*pb.Story, error) {
	// Assume this is from real db
	story := model.Story{
		ID:    2,
		Title: "Contoh title",
	}

	return story.ToProto(), nil
}

func (s *StoryService) Create(ctx context.Context, in *pb.CreateStoryRequest) (*pb.Story, error) {
	newStory, err := s.storyUsecase.Create(ctx, model.Story{
		Title:   in.Title,
		Content: in.Content,
		Author:  model.User{ID: in.AuthorId},
	})
	if err != nil {
		log.Errorf("failed create new story, error: %v", err)
		return nil, err
	}

	return newStory.ToProto(), nil
}

func (s *StoryService) Update(ctx context.Context, in *pb.UpdateStoryRequest) (*pb.Story, error) {
	panic("need implementation")
}

func (s *StoryService) Delete(ctx context.Context, in *pb.DeleteStoryRequest) (*emptypb.Empty, error) {
	panic("need implementation")
}
