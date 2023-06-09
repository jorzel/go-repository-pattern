package inmemory

import (
	"context"
	"fmt"

	downloading "github.com/jorzel/resource-downloader/internal/domain/downloader"
)

type InMemoryDownloaderRepository struct {
	storage map[downloading.UserId]downloading.ResourceDownloader
}

var _ downloading.DownloaderRepository = (*InMemoryDownloaderRepository)(nil)

func NewInMemoryDownloaderRepository() InMemoryDownloaderRepository {
	return InMemoryDownloaderRepository{
		storage: make(map[downloading.UserId]downloading.ResourceDownloader),
	}
}

func (r InMemoryDownloaderRepository) Get(ctx context.Context, userId downloading.UserId) (downloading.ResourceDownloader, error) {
	resourceDownloader, ok := r.storage[userId]
	if !ok {
		return downloading.ResourceDownloader{}, fmt.Errorf("downloader not found")
	}
	return resourceDownloader, nil
}

func (r InMemoryDownloaderRepository) Save(ctx context.Context, resourceDownloader downloading.ResourceDownloader) error {
	r.storage[resourceDownloader.UserId] = resourceDownloader
	return nil
}
