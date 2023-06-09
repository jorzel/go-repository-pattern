## Overview

The Repository Pattern is a software design pattern that provides an abstraction layer between the business logic of an application and the persistence layer (typically a database). It helps to separate the concerns and provides a consistent interface to access and manage data.

In the Repository Pattern, a repository acts as a mediator between the application and the data source. It encapsulates the logic for retrieving, storing, and querying data, allowing the application to interact with the repository instead of directly accessing the database or other data sources.

While the repository pattern is commonly associated with Domain-Driven Design (DDD) and Clean Architecture, it is a flexible and versatile pattern that can be used in various software development contexts.


## Application core

The main benefit is we can write application code that is totally agnostic of data provider (e.g. filesystem, database, external API).

There is a domain object `ResourceDownloader` that keeps already downloaded resources, current limit and is responsible for enforcing business invariant (the number of downloads cannot exceed the limit). 
```go
type UserId string
type ResourceId string

type ResourceDownloader struct {
	UserId    UserId       `json:"user_id"`
	Resources []ResourceId `json:"resources"`
	Limit     int          `json:"limit"`
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

```

Here is example of the application service service that is responsible for registering downloads if limit value is not reached. As you can see `DefaultDownloadService` is dependent on `DownloaderRepository`. But it is an interface that could be implemented by any infrastracture provider.

```go
type DownloaderRepository interface {
	Get(ctx context.Context, userId UserId) (ResourceDownloader, error)
	Save(ctx context.Context, downloader ResourceDownloader) error
}


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
```


## Repository implementations
### In-Memory

In memory repository implementation should not be used in production-ready system. However, it is a great building block for use case unit tests to avoid providing a mock.
```go

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
```

### Redis
Redis is an open-source, in-memory data structure store that can be used as a database, cache, and message broker. Redis is designed for high-performance data storage and retrieval, offering low-latency access to data. 
xw
```go
import	"github.com/redis/go-redis/v9"

type RedisDownloaderRepository struct {
	client *redis.Client
}

func NewRedisDownloaderRepository(client *redis.Client) RedisDownloaderRepository {
	return RedisDownloaderRepository{client: client}
}

var _ downloader.DownloaderRepository = (*RedisDownloaderRepository)(nil)

func (r RedisDownloaderRepository) Get(ctx context.Context, userId downloader.UserId) (downloader.ResourceDownloader, error) {
	result, err := r.client.Get(ctx, generateDownloaderKey(userId)).Result()
	if err == redis.Nil {
		return downloader.ResourceDownloader{}, fmt.Errorf("downloader not found")
	} else if err != nil {
		return downloader.ResourceDownloader{}, err
	}

	var resourceDownloader downloader.ResourceDownloader
	err = json.Unmarshal([]byte(result), &resourceDownloader)
	if err != nil {
		return downloader.ResourceDownloader{}, err
	}

	return resourceDownloader, nil
}

func (r RedisDownloaderRepository) Save(ctx context.Context, resourceDownloader downloader.ResourceDownloader) error {
	resourceDownloaderJSON, err := json.Marshal(resourceDownloader)
	if err != nil {
		return err
	}

	err = r.client.Set(ctx, generateDownloaderKey(resourceDownloader.UserId), resourceDownloaderJSON, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func generateDownloaderKey(userId downloader.UserId) string {
	return fmt.Sprintf("downloader:%s", userId)
}
``` 


### External API with cache (e.g. lru cache, Redis)
Repository pattern is not restricted only to databases and filesystems. In distributed system repository pattern can be expoited as an abstraction for an external REST API data source. In this example, we set up `CachedExternalDownloaderRepository`. We then use the repository's Get method to fetch a `ResourceDownloader` by `UserId`. If the object is not found in the cache (that would be accessed also by a repository pattern), the repository will make a call to the External API and cache the retrieved `UserLimit` for subsequent requests.

This approach allows you to leverage a local cache to minimize the number of external API calls and improve the performance of your application. The cache acts as a layer between the repository and the external API, providing a faster data retrieval option when the data is already available locally.

```go
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
```
