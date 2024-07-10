package repository

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/kodinggo/rest-api-service-golang-private-1/internal/model"
	log "github.com/sirupsen/logrus"
)

type storyRepository struct {
	db *sql.DB
}

func NewStoryRepository(db *sql.DB) model.StoryRepository {
	return &storyRepository{db: db}
}

func (r *storyRepository) FindAll(ctx context.Context, opt *model.StoryOptions) (results []model.Story, totalItems int64, err error) {
	selectQ := squirrel.Select("id, title, content, author_id, created_at").
		From("stories").
		OrderBy("created_at DESC")

	if opt.Search != "" {
		selectQ = selectQ.Where(squirrel.Like{
			"title": fmt.Sprintf("%%%s%%", opt.Search),
		})
	}

	if strings.ToLower(opt.SortBy) == "asc" {
		selectQ = selectQ.OrderBy("created_at ASC")
	}

	if opt.Cursor != "" {
		// decode base64 to time string
		decodedCursor, err := base64.StdEncoding.DecodeString(opt.Cursor)
		if err != nil {
			return results, 0, err
		}

		// parse ke time.Time
		cursorTime, err := time.Parse(time.RFC3339, string(decodedCursor))
		if err != nil {
			return results, 0, err
		}

		if strings.ToLower(opt.SortBy) == "asc" {
			selectQ = selectQ.Where(squirrel.Gt{"created_at": cursorTime})
		} else {
			selectQ = selectQ.Where(squirrel.Lt{"created_at": cursorTime})
		}
	}

	// Get raw sql query
	// queryRaw, _, _ := selectQ.ToSql()
	// fmt.Println(queryRaw)

	// Run SQL query to select multiple rows
	rows, err := selectQ.
		RunWith(r.db).
		QueryContext(ctx)
	if err != nil {
		log.Errorf("faled when run query sql select, error: %v", err)
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
			log.Errorf("failed to scan field, error: %v", err)
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
