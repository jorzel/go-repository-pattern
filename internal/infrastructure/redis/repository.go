package redis

import (
	"context"
	"encoding/json"
	"fmt"

	downloader "github.com/jorzel/resource-downloader/internal/domain/downloader"
	"github.com/redis/go-redis/v9"
)

type RedisDownloaderRepository struct {
	client *redis.Client
}

func NewRedisDownloaderRepository(client *redis.Client) RedisDownloaderRepository {
	return RedisDownloaderRepository{client: client}
}

var _ downloader.DownloaderRepository = (*RedisDownloaderRepository)(nil)

func (r RedisDownloaderRepository) Get(ctx context.Context, userId downloader.UserId) (downloader.ResourceDownloader, error) {
	result, err := r.client.Get(ctx, generateDownloaderKey(userId)).Result()
	if err == redis.Nil {
		return downloader.ResourceDownloader{}, fmt.Errorf("downloader not found")
	} else if err != nil {
		return downloader.ResourceDownloader{}, err
	}

	var resourceDownloader downloader.ResourceDownloader
	err = json.Unmarshal([]byte(result), &resourceDownloader)
	if err != nil {
		return downloader.ResourceDownloader{}, err
	}

	return resourceDownloader, nil
}

func (r RedisDownloaderRepository) Save(ctx context.Context, resourceDownloader downloader.ResourceDownloader) error {
	resourceDownloaderJSON, err := json.Marshal(resourceDownloader)
	if err != nil {
		return err
	}

	err = r.client.Set(ctx, generateDownloaderKey(resourceDownloader.UserId), resourceDownloaderJSON, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func generateDownloaderKey(userId downloader.UserId) string {
	return fmt.Sprintf("downloader:%s", userId)
}
