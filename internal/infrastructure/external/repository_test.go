package external

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"

	downloader "github.com/jorzel/resource-downloader/internal/domain/downloader"
	internalRedis "github.com/jorzel/resource-downloader/internal/infrastructure/redis"
	"github.com/stretchr/testify/assert"
)

func TestCachedExternalDownloaderRepositoryGet(t *testing.T) {

	userServiceClient := &mockUserServiceClient{}
	repository := setupExternalDownloaderRepository(t, userServiceClient)

	userServiceClient.requestMade = false
	resourceDownloader, err := repository.Get(context.Background(), downloader.UserId("1"))
	assert.NoError(t, err)
	assert.Equal(t, 10, resourceDownloader.Limit)
	assert.Equal(t, downloader.UserId("1"), resourceDownloader.UserId)
	assert.True(t, userServiceClient.requestMade)

	userServiceClient.requestMade = false
	resourceDownloader, err = repository.Get(context.Background(), downloader.UserId("1"))
	assert.NoError(t, err)
	assert.Equal(t, 10, resourceDownloader.Limit)
	assert.Equal(t, downloader.UserId("1"), resourceDownloader.UserId)
	assert.False(t, userServiceClient.requestMade)

	userServiceClient.requestMade = false
	resourceDownloader, err = repository.Get(context.Background(), downloader.UserId("2"))
	assert.NoError(t, err)
	assert.Equal(t, 10, resourceDownloader.Limit)
	assert.Equal(t, downloader.UserId("2"), resourceDownloader.UserId)
	assert.True(t, userServiceClient.requestMade)

}

// mockHTTPClient is a mock HTTP client for API requests
type mockUserServiceClient struct {
	requestMade bool
}

func (c *mockUserServiceClient) Get(ctx context.Context, userId downloader.UserId) (UserLimit, error) {
	// Simulate an API response with a user
	c.requestMade = true
	return UserLimit{UserId: userId, Limit: 10}, nil
}

func setupExternalDownloaderRepository(t *testing.T, userServiceClient *mockUserServiceClient) CachedExternalDownloaderRepository {
	mr := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     mr.Addr(),
		Password: "",
		DB:       0,
	})
	redisDownloaderRepository := internalRedis.NewRedisDownloaderRepository(redisClient)
	return NewCachedExternalDownloaderRepository(userServiceClient, redisDownloaderRepository)
}
