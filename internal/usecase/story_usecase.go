package usecase

import (
	"context"
	"errors"
	"sync"

	pbCategory "github.com/kodinggo/category-service-gp1/pb/category"
	pbComment "github.com/kodinggo/comment-service-gp1/pb/comment"
	"github.com/kodinggo/rest-api-service-golang-private-1/internal/model"
	"github.com/kodinggo/rest-api-service-golang-private-1/internal/worker"
	log "github.com/sirupsen/logrus"
)

type storyUsecase struct {
	storyRepo          model.StoryRepository
	userRepo           model.UserRepository
	grpcCommentClient  pbComment.CommentServiceClient
	grpcCategoryClient pbCategory.CategoryServiceClient
	workerClient       *worker.WorkerClient
}

func NewStoryUsecase(
	storyRepo model.StoryRepository,
	userRepo model.UserRepository,
	grpcCommentClient pbComment.CommentServiceClient,
	grpcCategoryClient pbCategory.CategoryServiceClient,
	workerClient *worker.WorkerClient,
) model.StoryUsecase {
	return &storyUsecase{
		storyRepo:          storyRepo,
		userRepo:           userRepo,
		grpcCommentClient:  grpcCommentClient,
		grpcCategoryClient: grpcCategoryClient,
		workerClient:       workerClient,
	}
}

func (u *storyUsecase) FindAll(ctx context.Context, opt *model.StoryOptions) (results []model.Story, totalItems int64, err error) {
	results, _, err = u.storyRepo.FindAllES(ctx, opt)
	if err != nil {
		log.Errorf("faled when find all stories from repo, error: %v", err)
		return
	}

	// resolve author
	var wg sync.WaitGroup
	wg.Add(len(results) * 3)

	for idx, result := range results {
		go func(idx int, result model.Story) {
			defer wg.Done()
			// Manggil gRPC comment
			commentsPB, err := u.grpcCommentClient.FindAllCommentsByStoryID(ctx,
				&pbComment.FindAllCommentsByStoryIDRequest{
					StoryId: result.ID,
				})
			if err != nil || commentsPB == nil {
				log.Errorf("failed when resolve comments, storyID: %d, error: %v", result.ID, err)
				return
			}

			// Convert dari protobuf comment ke main comment entity
			var comments []model.StoryComment
			for _, pbComment := range commentsPB.Comments {
				comments = append(comments, model.StoryComment{
					ID:      pbComment.Id,
					Comment: pbComment.Comment,
				})
			}

			results[idx].Comments = comments
		}(idx, result)

		go func(idx int, result model.Story) {
			defer wg.Done()
			// Manggil gRPC category
			categoryPB, err := u.grpcCategoryClient.FindCategoryByID(ctx,
				&pbCategory.FindCategoryByIDRequest{
					Id: result.ID,
				})
			if err != nil || categoryPB == nil {
				log.Errorf("failed when resolve category, storyID: %d, error: %v", result.ID, err)
				return
			}

			results[idx].Category = model.StoryCategory{
				ID:   categoryPB.Id,
				Name: categoryPB.Name,
			}
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

func (u *storyUsecase) FindByID(ctx context.Context, id int64) (*model.Story, error) {
	story, err := u.storyRepo.FindByID(ctx, id)
	if err != nil {
		log.Errorf("failed find story by id, error: %v", err)
	}
	return story, err
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

	// Send Task to Queue
	_, err = u.workerClient.SendEmail(worker.SendEmail{
		From:    "john@gmail.com",
		To:      "mark@gmail.com",
		Subject: "Test",
	})
	if err != nil {
		log.Errorf("failed send task for sending email to worker, error: %v", err)
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
