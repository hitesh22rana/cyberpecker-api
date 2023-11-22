package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	cybernews "github.com/hitesh22rana/cyberpecker-api/src/cybernews"

	"github.com/redis/go-redis/v9"
)

const dataExpirationTime = time.Hour * 6

var (
	ErrorConnecting   = errors.New("unable to connect to the database")
	ErrorJsonEncoding = errors.New("json encoding error")
	ErrorJsonParsing  = errors.New("json parsing error")
	ErrorNewsNotFound = errors.New("news not found")
)

type RedisClient struct {
	*redis.Client
}

func NewRedisClient(config *redis.Options) *RedisClient {
	return &RedisClient{
		redis.NewClient(config),
	}
}

func (client *RedisClient) Health() error {
	if err := client.Ping(context.Background()).Err(); err != nil {
		return ErrorConnecting
	}

	return nil
}

func (c *RedisClient) SetNews(ctx context.Context, key string, value []cybernews.News) error {
	data, err := json.Marshal(value)
	if err != nil {
		return ErrorJsonEncoding
	}

	err = c.Set(ctx, key, data, dataExpirationTime).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *RedisClient) GetNews(ctx context.Context, key string) ([]cybernews.News, error) {
	data, err := c.Get(ctx, key).Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrorNewsNotFound
		}

		return nil, ErrorConnecting
	}

	var news []cybernews.News
	err = json.Unmarshal([]byte(data), &news)
	if err != nil {
		return nil, ErrorJsonParsing
	}

	return news, nil
}
