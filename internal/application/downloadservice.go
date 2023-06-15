package application

import (
	"context"

	downloader "github.com/jorzel/resource-downloader/internal/domain/downloader"
)

type DownloadService interface {
	DownloadResource(userId downloader.UserId, resourceId downloader.ResourceId) error
}

type DefaultDownloadService struct {
	DownloaderRepository downloader.DownloaderRepository
}

func NewDownloadService(downloaderRepository downloader.DownloaderRepository) DefaultDownloadService {
	return DefaultDownloadService{
		DownloaderRepository: downloaderRepository,
	}
}

func (s DefaultDownloadService) DownloadResource(ctx context.Context, userId downloader.UserId, resourceId downloader.ResourceId) error {
	resourceDownloader, err := s.DownloaderRepository.Get(ctx, userId)
	if err != nil {
		return err
	}
	err = resourceDownloader.RegisterDownload(resourceId)
	if err != nil {
		return err
	}

	// Take action to perform the download

	err = s.DownloaderRepository.Save(ctx, resourceDownloader)
	if err != nil {
		return err
	}
	return nil
}
