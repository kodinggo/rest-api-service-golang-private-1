package model

import (
	"context"
	"time"

	pb "github.com/kodinggo/rest-api-service-golang-private-1/pb/story"
)

type Story struct {
	ID         int64          `json:"id"`
	Title      string         `json:"title" validate:"required"`
	Content    string         `json:"content" validate:"required"`
	Author     User           `json:"author"`
	CreatedAt  time.Time      `json:"created_at"`
	Comments   []StoryComment `json:"comments"`
	Category   StoryCategory  `json:"category"`
	CategoryID int64          `json:"-"`
	UpdatedAt  time.Time      `json:"-"`
	DeletedAt  time.Time      `json:"-"`
}

type StoryOptions struct {
	Search string `query:"search"`
	SortBy string `query:"sort_by"`
	Cursor string `query:"cursor"`
}

type StoryComment struct {
	ID      int64  `json:"id"`
	Comment string `json:"comment"`
}

type StoryCategory struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// TODO: Add other interface methods
type StoryUsecase interface {
	FindAll(ctx context.Context, opt *StoryOptions) ([]Story, int64, error)
	Create(ctx context.Context, data Story) (*Story, error)
	Update(ctx context.Context, data Story) (*Story, error)
}

// TODO: Add other interface methods
type StoryRepository interface {
	FindAll(ctx context.Context, opt *StoryOptions) ([]Story, int64, error)
	FindByID(ctx context.Context, id int64) (*Story, error)
	Create(ctx context.Context, data Story) (*Story, error)
	Update(ctx context.Context, data Story) (*Story, error)
}

func (s Story) ToProto() *pb.Story {
	return &pb.Story{
		Id:      s.ID,
		Title:   s.Title,
		Content: s.Content,
		Author:  s.Author.ToProto(),
	}
}

func NewStoryFromProto(p *pb.Story) Story {
	story := Story{
		ID:      p.Id,
		Title:   p.Title,
		Content: p.Content,
	}
	story.Author = NewUserFromProto(p.Author)

	return story
}
