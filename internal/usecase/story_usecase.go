package usecase

import (
	"context"

	"github.com/kodinggo/rest-api-service-golang-private-1/internal/model"
	log "github.com/sirupsen/logrus"
)

type storyUsecase struct {
	storyRepo model.StoryRepository
}

func NewStoryUsecase(storyRepo model.StoryRepository) model.StoryUsecase {
	return &storyUsecase{storyRepo: storyRepo}
}

func (u *storyUsecase) FindAll(ctx context.Context, opt *model.StoryOptions) (results []model.Story, totalItems int64, err error) {
	results, _, err = u.storyRepo.FindAll(ctx, opt)
	if err != nil {
		log.Errorf("faled when find all stories from repo, error: %v", err)
		return
	}

	// TODO: resolve author

	return
}
