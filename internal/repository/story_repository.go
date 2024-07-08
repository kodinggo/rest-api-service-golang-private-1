package repository

import (
	"context"
	"database/sql"
	"log"

	"github.com/kodinggo/rest-api-service-golang-private-1/internal/model"
)

type storyRepository struct {
	db *sql.DB
}

func NewStoryRepository(db *sql.DB) model.StoryRepository {
	return &storyRepository{db: db}
}

func (r *storyRepository) FindAll(ctx context.Context, opt *model.StoryOptions) (results []model.Story, totalItems int64, err error) {
	// Run SQL query to select multiple rows
	rows, err := r.db.QueryContext(ctx, `SELECT id, title, content, author_id, created_at 
		FROM stories ORDER BY created_at DESC`)
	if err != nil {
		return
	}

	for rows.Next() {
		var story model.Story
		var authorID int64

		// Scan fields
		err = rows.Scan(&story.ID,
			&story.Title,
			&story.Content,
			&authorID,
			&story.CreatedAt)
		if err != nil {
			log.Printf("failed to scan field, error: %v", err)
			continue
		}

		// Collect all rows to "results"
		story.Author = model.User{
			ID: authorID,
		}
		results = append(results, story)
	}

	// TODO: Please find total items

	return
}
