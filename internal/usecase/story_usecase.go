package usecase

import (
	"context"
	"errors"
	"sync"

	pbComment "github.com/kodinggo/comment-service-gp1/pb/comment"
	"github.com/kodinggo/rest-api-service-golang-private-1/internal/model"
	log "github.com/sirupsen/logrus"
)

type storyUsecase struct {
	storyRepo         model.StoryRepository
	userRepo          model.UserRepository
	grpcCommentClient pbComment.CommentServiceClient
}

func NewStoryUsecase(
	storyRepo model.StoryRepository,
	userRepo model.UserRepository,
	grpcCommentClient pbComment.CommentServiceClient,
) model.StoryUsecase {
	return &storyUsecase{
		storyRepo:         storyRepo,
		userRepo:          userRepo,
		grpcCommentClient: grpcCommentClient,
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
	wg.Add(len(results) * 2)

	for idx, result := range results {
		go func(idx int, result model.Story) {
			defer wg.Done()
			// Manggil gRPC comment
			pbComments, err := u.grpcCommentClient.FindAllCommentsByStoryID(ctx,
				&pbComment.FindAllCommentsByStoryIDRequest{
					StoryId: result.ID,
				})
			if err != nil || pbComments == nil {
				log.Errorf("failed when resolve comments, storyID: %d, error: %v", result.ID, err)
				return
			}

			// Convert dari protobuf comment ke main comment entity
			var comments []model.StoryComment
			for _, pbComment := range pbComments.Comments {
				comments = append(comments, model.StoryComment{
					ID:      pbComment.Id,
					Comment: pbComment.Comment,
				})
			}

			results[idx].Comments = comments
		}(idx, result)

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

func (u *storyUsecase) Update(ctx context.Context, data model.Story) (*model.Story, error) {
	_, err := u.userRepo.FindByID(ctx, data.Author.ID)
	if err != nil {
		return nil, model.NewErrorUnAuthorized("invalid author")
	}

	story, err := u.storyRepo.FindByID(ctx, data.ID)
	if err != nil || story == nil {
		return nil, model.NewErrorNotFound("story not found") // TODO: Create error type
	}

	insertedData, err := u.storyRepo.Update(ctx, data)
	if err != nil {
		log.Errorf("failed create new story, error: %v", err)
	}
	return insertedData, err
}
