package downloading

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDownloaderSuccesfulWhenLimitIsNotReached(t *testing.T) {
	downloader := NewResourceDownloader("user1", []ResourceId{}, 2)

	err := downloader.RegisterDownload("resource1")

	assert.NoError(t, err)
	assert.Equal(t, []ResourceId{"resource1"}, downloader.Resources)
}

func TestDownloaderFailedWhenLimitReached(t *testing.T) {
	downloader := NewResourceDownloader("user1", []ResourceId{"resource1", "resource2"}, 2)

	err := downloader.RegisterDownload("resource3")

	assert.Error(t, err)
	assert.Equal(t, []ResourceId{"resource1", "resource2"}, downloader.Resources)
}
