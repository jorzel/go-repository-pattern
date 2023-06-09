package inmemory

import (
	"context"
	"fmt"

	downloader "github.com/jorzel/resource-downloader/internal/domain/downloader"
)

type InMemoryDownloaderRepository struct {
	Storage map[downloader.UserId]downloader.ResourceDownloader
}

var _ downloader.DownloaderRepository = (*InMemoryDownloaderRepository)(nil)

func NewInMemoryDownloaderRepository() InMemoryDownloaderRepository {
	return InMemoryDownloaderRepository{
		Storage: make(map[downloader.UserId]downloader.ResourceDownloader),
	}
}

func (r InMemoryDownloaderRepository) Get(ctx context.Context, userId downloader.UserId) (downloader.ResourceDownloader, error) {
	resourceDownloader, ok := r.Storage[userId]
	if !ok {
		return downloader.ResourceDownloader{}, fmt.Errorf("downloader not found")
	}
	return resourceDownloader, nil
}

func (r InMemoryDownloaderRepository) Save(ctx context.Context, resourceDownloader downloader.ResourceDownloader) error {
	r.Storage[resourceDownloader.UserId] = resourceDownloader
	return nil
}
