package repository

import "fmt"

func newStoryByIDCacheKey(id int) string {
	return fmt.Sprintf("story:%d", id)
}
