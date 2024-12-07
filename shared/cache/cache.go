package cache

import (
	"context"
	"sync"

	"github.com/go-redis/redis/v8"
)

var (
	client *redis.Client
	once   sync.Once
)

func InitRedis(addr, password string, db int) {
	once.Do(func() {
		client = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		})
	})
}

func GetRedisClient() *redis.Client {
	return client
}

func Ping(ctx context.Context) error {
	return client.Ping(ctx).Err()
}
