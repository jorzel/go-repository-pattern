package application

import (
	"context"
	"testing"

	downloader "github.com/jorzel/resource-downloader/internal/domain/downloader"
	"github.com/jorzel/resource-downloader/internal/infrastructure/inmemory"
	"github.com/stretchr/testify/assert"
)

func TestDownloadService(t *testing.T) {
	repository := inmemory.NewInMemoryDownloaderRepository()
	userId := downloader.UserId("user1")
	resourceDownloader := downloader.NewResourceDownloader(userId, []downloader.ResourceId{}, 5)
	repository.Save(context.Background(), *resourceDownloader)
	service := NewDownloadService(repository)

	err := service.DownloadResource(context.Background(), userId, "resource1")

	assert.NoError(t, err)
	resourceDownloaderFromRepo, _ := repository.Get(context.Background(), userId)
	assert.NoError(t, err)
	assert.Equal(t, []downloader.ResourceId{"resource1"}, resourceDownloaderFromRepo.Resources)
}

func TestDownloadServiceReachLimit(t *testing.T) {
	repository := inmemory.NewInMemoryDownloaderRepository()
	userId := downloader.UserId("user1")
	resourceDownloader := downloader.NewResourceDownloader(userId, []downloader.ResourceId{}, 2)
	repository.Save(context.Background(), *resourceDownloader)
	service := NewDownloadService(repository)

	err := service.DownloadResource(context.Background(), userId, "resource1")
	assert.NoError(t, err)
	err = service.DownloadResource(context.Background(), userId, "resource2")
	assert.NoError(t, err)
	err = service.DownloadResource(context.Background(), userId, "resource3")
	assert.Error(t, err)
	resourceDownloaderFromRepo, err := repository.Get(context.Background(), userId)
	assert.NoError(t, err)
	assert.Equal(t, []downloader.ResourceId{"resource1", "resource2"}, resourceDownloaderFromRepo.Resources)
}
