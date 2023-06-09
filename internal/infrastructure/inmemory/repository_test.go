package inmemory

import (
	"context"
	"testing"

	downloader "github.com/jorzel/resource-downloader/internal/domain/downloader"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryDownloaderRepositoryGetWhenDownloaderExists(t *testing.T) {
	repository := NewInMemoryDownloaderRepository()
	userId := downloader.UserId("user1")
	resourceDownloader := downloader.NewResourceDownloader(userId, []downloader.ResourceId{}, 5)
	repository.Storage[userId] = *resourceDownloader

	result, err := repository.Get(context.Background(), userId)

	assert.NoError(t, err)
	assert.Equal(t, *resourceDownloader, result)
}

func TestInMemoryDownloaderRepositorySave(t *testing.T) {
	repository := NewInMemoryDownloaderRepository()
	userId := downloader.UserId("user1")
	resourceDownloader := downloader.NewResourceDownloader(userId, []downloader.ResourceId{}, 5)

	err := repository.Save(context.Background(), *resourceDownloader)

	assert.NoError(t, err)
	assert.Equal(t, *resourceDownloader, repository.Storage[userId])
}
