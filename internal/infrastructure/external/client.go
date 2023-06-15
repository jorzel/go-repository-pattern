package external

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	downloader "github.com/jorzel/resource-downloader/internal/domain/downloader"
)

type UserServiceClient interface {
	Get(ctx context.Context, userId downloader.UserId) (UserLimit, error)
}

type DefaultUserServiceClient struct {
	baseURL string
}

func NewDefaultUserServiceClient(baseURL string) DefaultUserServiceClient {
	return DefaultUserServiceClient{baseURL: baseURL}
}

func (c *DefaultUserServiceClient) Get(ctx context.Context, userId downloader.UserId) (UserLimit, error) {
	url := fmt.Sprintf("%s/users/%s", c.baseURL, userId)

	response, err := http.Get(url)
	if err != nil {
		return UserLimit{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return UserLimit{}, fmt.Errorf("failed to get user with ID %s", userId)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return UserLimit{}, err
	}

	var userLimit UserLimit
	err = json.Unmarshal(body, &userLimit)
	if err != nil {
		return UserLimit{}, err
	}
	return userLimit, nil
}
