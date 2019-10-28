package image

import (
	"context"
	"strings"
	"time"

	imageentity "github.com/albertwidi/go_project_example/entity/image"
	"github.com/albertwidi/go_project_example/lib/redis"
)

// Repository of image
type Repository struct {
	redis *redis.Redis
}

// New repository for image
func New(redis *redis.Redis) *Repository {
	r := Repository{
		redis: redis,
	}
	return &r
}

func createImageKey(id string) string {
	return strings.Join([]string{"image_temp", id}, ":")
}

// SaveTempPath image path
func (r Repository) SaveTempPath(ctx context.Context, id, originalPath string, expiryTime time.Duration) error {
	key := createImageKey(id)
	_, err := r.redis.SetEX(key, originalPath, int(expiryTime.Seconds()))
	return err
}

// GetTempPath will return the original path from a temporary id
func (r Repository) GetTempPath(ctx context.Context, id string) (string, error) {
	key := createImageKey(id)
	out, err := r.redis.Get(key)
	if err != nil {
		if redis.IsErrNil(err) {
			return "", imageentity.ErrTempPathNotFound
		}
		return "", err
	}

	return out, err
}