package usecase

import (
	"context"
	"errors"
	"sync"

	"github.com/kodinggo/rest-api-service-golang-private-1/internal/model"
	log "github.com/sirupsen/logrus"
)

type storyUsecase struct {
	storyRepo model.StoryRepository
	userRepo  model.UserRepository
}

func NewStoryUsecase(
	storyRepo model.StoryRepository,
	userRepo model.UserRepository,
) model.StoryUsecase {
	return &storyUsecase{
		storyRepo: storyRepo,
		userRepo:  userRepo,
	}
}

func (u *storyUsecase) FindAll(ctx context.Context, opt *model.StoryOptions) (results []model.Story, totalItems int64, err error) {
	results, _, err = u.storyRepo.FindAll(ctx, opt)
	if err != nil {
		log.Errorf("faled when find all stories from repo, error: %v", err)
		return
	}

	// resolve author
	var wg sync.WaitGroup
	wg.Add(len(results))

	for idx, result := range results {
		go func(idx int, result model.Story) {
			defer wg.Done()
			user, err := u.userRepo.FindByID(ctx, result.Author.ID)
			if err != nil {
				log.Errorf("failed when resolve author, authorID: %d, error: %v", result.Author.ID, err)
				return
			}
			if user == nil {
				return
			}
			results[idx].Author = *user
		}(idx, result)
	}

	wg.Wait()

	return
}

func (u *storyUsecase) Create(ctx context.Context, data model.Story) (*model.Story, error) {
	_, err := u.userRepo.FindByID(ctx, data.Author.ID)
	if err != nil {
		return nil, errors.New("invalid author")
	}

	insertedData, err := u.storyRepo.Create(ctx, data)
	if err != nil {
		log.Errorf("failed create new story, error: %v", err)
	}
	return insertedData, err
}
