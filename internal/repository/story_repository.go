package repository

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/kodinggo/rest-api-service-golang-private-1/internal/db"
	"github.com/kodinggo/rest-api-service-golang-private-1/internal/model"
	"github.com/olivere/elastic/v7"
	log "github.com/sirupsen/logrus"
)

type storyRepository struct {
	db          *sql.DB
	redisClient model.RedisClient
	esClient    *elastic.Client
}

func NewStoryRepository(db *sql.DB, redisClient model.RedisClient, esClient *elastic.Client) model.StoryRepository {
	return &storyRepository{db: db, redisClient: redisClient, esClient: esClient}
}

func (r *storyRepository) FindAll(ctx context.Context, opt *model.StoryOptions) (results []model.Story, totalItems int64, err error) {
	// Cek apakah data ada pada redis, jika ada maka ambil dari redis
	// Jika tidak maka lanjut ke mysql
	cacheKey := newStoriesCacheKey(opt)
	err = r.redisClient.HGet(ctx, storiesBucketKey, cacheKey, &results)
	if err != nil {
		log.Errorf("failed get data from redis, error: %v", err)
	}
	if len(results) > 0 {
		return
	}

	selectQ := sq.Select("id, title, content, author_id, created_at").
		From("stories").
		OrderBy("created_at DESC")

	if opt != nil && opt.Search != "" {
		selectQ = selectQ.Where(sq.Like{
			"title": fmt.Sprintf("%%%s%%", opt.Search),
		})
	}

	if opt != nil && strings.ToLower(opt.SortBy) == "asc" {
		selectQ = selectQ.OrderBy("created_at ASC")
	}

	if opt != nil && opt.Cursor != "" {
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
			selectQ = selectQ.Where(sq.Gt{"created_at": cursorTime})
		} else {
			selectQ = selectQ.Where(sq.Lt{"created_at": cursorTime})
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

	// Set data to redis
	// Store ke redis
	err = r.redisClient.HSet(ctx, storiesBucketKey, cacheKey, results, 5*time.Minute)
	if err != nil {
		log.Errorf("failed set data to redis, error: %v", err)
	}

	return
}

func (r *storyRepository) FindAllES(ctx context.Context, opt *model.StoryOptions) (results []model.Story, totalItems int64, err error) {
	boolQuery := elastic.NewBoolQuery()

	if opt != nil && opt.Search != "" {
		boolQuery = boolQuery.Should(
			elastic.NewMatchQuery("title", opt.Search),
			elastic.NewMatchQuery("content", opt.Search),
		).
			MinimumShouldMatch("1")
	}

	searchResult, err := r.esClient.Search().
		Index(db.IndexName).
		Query(boolQuery).
		Do(context.Background())
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	for _, res := range searchResult.Hits.Hits {
		var result model.Story
		json.Unmarshal(res.Source, &result)
		results = append(results, result)
	}

	return
}

func (r *storyRepository) Create(ctx context.Context, data model.Story) (result *model.Story, err error) {
	createdAt := time.Now().UTC()
	res, err := sq.Insert("stories").
		Columns("title", "content", "author_id", "created_at").
		Values(data.Title, data.Content, data.Author.ID, createdAt).
		RunWith(r.db).
		ExecContext(ctx)
	if err != nil {
		log.WithField("data", data).
			Errorf("failed when insert data to story, error: %v", err)
		return
	}
	data.ID, _ = res.LastInsertId()
	data.CreatedAt = createdAt
	result = &data

	go func() {
		// Invalidate redis
		err = r.redisClient.HDelByBucketKey(context.Background(), storiesBucketKey)
		if err != nil {
			log.Errorf("failed when delete data from redis, error: %v", err)
		}

		// Insert data to elasticsearch
		_, err = r.esClient.Index().
			Index(db.IndexName).
			Id(fmt.Sprintf("%d", data.ID)).
			BodyJson(result).
			Do(context.Background())
		if err != nil {
			log.Errorf("faild indexing document ID=%d: %s", data.ID, err)
		}
	}()

	return
}

func (r *storyRepository) Update(ctx context.Context, data model.Story) (result *model.Story, err error) {
	updatedAt := time.Now().UTC()
	_, err = sq.Update("stories").
		Set("title", data.Title).
		Set("content", data.Content).
		Set("created_at", updatedAt).
		Where(sq.Eq{"id": data.ID}).
		RunWith(r.db).
		ExecContext(ctx)
	if err != nil {
		log.WithField("data", data).
			Errorf("failed when insert data to story, error: %v", err)
		return
	}
	data.UpdatedAt = updatedAt
	result = &data

	go func() {
		// Invalidate redis
		err = r.redisClient.Del(context.Background(), newStoryByIDCacheKey(int(data.ID)))
		if err != nil {
			log.Errorf("failed when delete data from redis, error: %v", err)
		}
		err = r.redisClient.HDelByBucketKey(context.Background(), storiesBucketKey)
		if err != nil {
			log.Errorf("failed when delete data from redis, error: %v", err)
		}

		// Update data from elasticsearch
		_, err = r.esClient.Update().
			Index(db.IndexName).
			Id(fmt.Sprintf("%d", data.ID)).
			Doc(result).
			Do(context.Background())
		if err != nil {
			log.Errorf("faild updating document ID=%d: %s", data.ID, err)
		}
	}()

	return
}

func (r *storyRepository) FindByID(ctx context.Context, id int64) (*model.Story, error) {
	row := sq.Select("id, title, content, author_id, created_at").
		From("stories").
		OrderBy("created_at DESC").
		Where(sq.Eq{"id": id}).
		RunWith(r.db).
		QueryRowContext(ctx)

	var story model.Story
	var authorID int64

	// Cek apakah data ada pada redis, jika ada maka ambil dari redis
	// Jika tidak maka lanjut ke mysql
	cacheKey := newStoryByIDCacheKey(int(id))
	err := r.redisClient.Get(ctx, cacheKey, &story)
	if err != nil {
		log.Errorf("failed get data from redis, error: %v", err)
	}
	if story.ID != 0 {
		return &story, nil
	}

	// Scan fields
	err = row.Scan(&story.ID,
		&story.Title,
		&story.Content,
		&authorID,
		&story.CreatedAt)
	if err != nil {
		log.Errorf("failed to scan field, error: %v", err)
		return nil, err
	}

	// Collect all rows to "results"
	story.Author = model.User{
		ID: authorID,
	}

	// Store ke redis
	err = r.redisClient.Set(ctx, cacheKey, story, 5*time.Minute)
	if err != nil {
		log.Errorf("failed set data to redis, error: %v", err)
	}

	return &story, nil
}
