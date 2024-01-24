package redis

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/hosseintrz/gaterun/config/models"
)

var client *redis.Client
var once sync.Once

func InitRedis(ctx context.Context, cfg models.RedisConfig) (err error) {
	once.Do(func() {
		client = redis.NewClient(&redis.Options{
			Addr:     cfg.Address,
			Password: "",
			DB:       cfg.DB,
		})

		_, err = client.Ping(ctx).Result()
		if err != nil {
			return
		}

		fmt.Println("connected to redis")
	})

	return
}

func GetClient(ctx context.Context) *redis.Client {
	return client.WithContext(ctx)
}
