package database

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

func NewRedisClient(config *redis.Options) *redis.Client {
	return redis.NewClient(config)
}

func Health(client *redis.Client) error {
	if err := client.Ping(context.Background()).Err(); err != nil {
		return ErrorConnecting
	}

	return nil
}

func SaveNews(client *redis.Client, ctx context.Context, key string, value []cybernews.News) error {
	data, err := json.Marshal(value)
	if err != nil {
		return ErrorJsonEncoding
	}

	err = client.Set(ctx, key, data, dataExpirationTime).Err()
	if err != nil {
		return err
	}
	return nil
}

func RetrieveNews(client *redis.Client, ctx context.Context, key string) ([]cybernews.News, error) {
	data, err := client.Get(ctx, key).Result()
	fmt.Println()

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

// dial tcp [::1]:6379: connectex: No connection could be made because the target machine actively refused it.
// redis: nil
