package main

import (
	"fmt"

	"github.com/jorzel/resource-downloader/internal/infrastructure/external"
	internalRedis "github.com/jorzel/resource-downloader/internal/infrastructure/redis"
	"github.com/redis/go-redis/v9"
)

func main() {
	fmt.Println("Repository pattern example")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "secret",
		DB:       0,
	})
	cache := internalRedis.NewRedisDownloaderRepository(redisClient)
	userServiceClient := external.NewDefaultUserServiceClient("http://localhost")
	external.NewCachedExternalDownloaderRepository(
		&userServiceClient, cache,
	)
}
