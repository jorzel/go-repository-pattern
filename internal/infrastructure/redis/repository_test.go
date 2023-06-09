package redis

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/alicebob/miniredis/v2"
	downloader "github.com/jorzel/resource-downloader/internal/domain/downloader"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRedisDownloaderRepositoryGetWhenDownloaderExists(t *testing.T) {
	redisDownloaderRepository := setupRedisDownloaderRepository(t)
	redisClient := redisDownloaderRepository.client
	userId := downloader.UserId("user1")
	resourceDownloader := downloader.NewResourceDownloader(userId, []downloader.ResourceId{}, 5)
	serializedResourceDownloader, err := json.Marshal(resourceDownloader)
	assert.NoError(t, err)
	_, err = redisClient.Set(context.Background(), generateDownloaderKey(userId), string(serializedResourceDownloader), 0).Result()
	assert.NoError(t, err)

	result, err := redisDownloaderRepository.Get(context.Background(), userId)

	assert.NoError(t, err)
	assert.Equal(t, *resourceDownloader, result)
}

func TestRedisDownloaderRepositoryGetWhenDownloaderDoesNotExist(t *testing.T) {
	redisDownloaderRepository := setupRedisDownloaderRepository(t)
	userId := downloader.UserId("user1")

	result, err := redisDownloaderRepository.Get(context.Background(), userId)

	assert.Error(t, err)
	assert.Equal(t, downloader.ResourceDownloader{}, result)
}

func TestRedisDownloaderRepositorySave(t *testing.T) {
	redisDownloaderRepository := setupRedisDownloaderRepository(t)
	redisClient := redisDownloaderRepository.client
	userId := downloader.UserId("user1")
	resourceDownloader := downloader.NewResourceDownloader(userId, []downloader.ResourceId{}, 5)
	serializedResourceDownloader, err := json.Marshal(resourceDownloader)
	assert.NoError(t, err)

	err = redisDownloaderRepository.Save(context.Background(), *resourceDownloader)

	assert.NoError(t, err)
	result, err := redisClient.Get(context.Background(), generateDownloaderKey(userId)).Result()
	assert.NoError(t, err)
	assert.Equal(t, string(result), string(serializedResourceDownloader))
}

func setupRedisDownloaderRepository(t *testing.T) RedisDownloaderRepository {
	mr := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     mr.Addr(),
		Password: "",
		DB:       0,
	})
	redisDownloaderRepository := NewRedisDownloaderRepository(redisClient)
	return redisDownloaderRepository
}
