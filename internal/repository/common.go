package repository

import (
	"fmt"

	"github.com/kodinggo/rest-api-service-golang-private-1/internal/model"
)

const storiesBucketKey = "stories"

func newStoryByIDCacheKey(id int) string {
	return fmt.Sprintf("story:%d", id)
}

func newStoriesCacheKey(opt *model.StoryOptions) string {
	var search, sortBy, cursor string
	if opt != nil && opt.Search != "" {
		search = opt.Search
	}
	if opt != nil && opt.SortBy != "" {
		sortBy = opt.SortBy
	}
	if opt != nil && opt.Cursor != "" {
		cursor = opt.Cursor
	}
	return fmt.Sprintf("stories:search:%s:sort_by:%s:cursor:%s", search, sortBy, cursor)
}
