package downloading

import "context"

type DownloaderRepository interface {
	Get(ctx context.Context, userId UserId) (ResourceDownloader, error)
	Save(ctx context.Context, downloader ResourceDownloader) error
}
