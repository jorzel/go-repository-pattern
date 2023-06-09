package downloading

import "fmt"

type UserId string
type ResourceId string

type ResourceDownloader struct {
	UserId    UserId       `json:"user_id"`
	Resources []ResourceId `json:"resources"`
	Limit     int          `json:"limit"`
}

func EmptyResourceDownloader() ResourceDownloader {
	return ResourceDownloader{}
}

func NewResourceDownloader(userId UserId, resources []ResourceId, limit int) *ResourceDownloader {
	return &ResourceDownloader{
		UserId:    userId,
		Resources: resources,
		Limit:     limit,
	}
}

func (d *ResourceDownloader) isLimitReached() bool {
	return len(d.Resources) >= d.Limit
}

func (d *ResourceDownloader) RegisterDownload(resourceId ResourceId) error {
	if d.isLimitReached() {
		return fmt.Errorf("limit reached")
	}
	d.Resources = append(d.Resources, resourceId)
	return nil
}
