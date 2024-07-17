package model

import (
	"context"
	"time"
)

type Story struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title" validate:"required"`
	Content   string    `json:"content" validate:"required"`
	Author    User      `json:"author"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"-"`
	DeletedAt time.Time `json:"-"`
}

type StoryOptions struct {
	Search string `query:"search"`
	SortBy string `query:"sort_by"`
	Cursor string `query:"cursor"`
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
