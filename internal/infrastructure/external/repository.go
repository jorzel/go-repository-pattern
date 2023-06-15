package external

import (
	"context"

	downloader "github.com/jorzel/resource-downloader/internal/domain/downloader"
)

type UserLimit struct {
	UserId downloader.UserId `json:"user_id"`
	Limit  int               `json:"limit"`
}

type CachedExternalDownloaderRepository struct {
	cache  downloader.DownloaderRepository
	client UserServiceClient
}

var _ downloader.DownloaderRepository = (*CachedExternalDownloaderRepository)(nil)

func NewCachedExternalDownloaderRepository(client UserServiceClient, cache downloader.DownloaderRepository) CachedExternalDownloaderRepository {
	return CachedExternalDownloaderRepository{
		cache:  cache,
		client: client,
	}
}

func (r CachedExternalDownloaderRepository) Get(ctx context.Context, userId downloader.UserId) (downloader.ResourceDownloader, error) {
	resourceDownloaderFromCache, err := r.cache.Get(ctx, userId)
	if err == nil {
		return resourceDownloaderFromCache, nil
	}

	userLimit, err := r.fetchUserLimitsFromExternalAPI(ctx, userId)
	if err != nil {
		return downloader.ResourceDownloader{}, err
	}

	resourceDownloader := downloader.NewResourceDownloader(userId, []downloader.ResourceId{}, userLimit.Limit)
	_ = r.cache.Save(ctx, *resourceDownloader)
	return *resourceDownloader, nil

}

func (r CachedExternalDownloaderRepository) Save(ctx context.Context, resourceDownloader downloader.ResourceDownloader) error {
	r.cache.Save(ctx, resourceDownloader)
	return nil
}

func (r CachedExternalDownloaderRepository) fetchUserLimitsFromExternalAPI(ctx context.Context, userId downloader.UserId) (UserLimit, error) {
	userLimit, err := r.client.Get(ctx, userId)
	if err != nil {
		return UserLimit{}, err
	}
	return userLimit, nil
}
