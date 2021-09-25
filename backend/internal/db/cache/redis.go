package cache

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	client *redis.Client
}

type Options struct {
	Host     string
	Port     string
	Password string
	DB       string
}

var (
	redisClient *RedisClient
	lock        = &sync.Mutex{}
)

func NewRedisClient(ctx context.Context, options *Options) (*RedisClient, error) {
	lock.Lock()
	defer lock.Unlock()

	if redisClient == nil {
		client, err := getRedisClient(ctx, options)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize redis client, error is: %s", err)
		}
		redisClient = &RedisClient{
			client: client,
		}
		return redisClient, nil
	}

	return redisClient, nil
}

func getRedisClient(ctx context.Context, options *Options) (*redis.Client, error) {
	dbOpt, _ := strconv.Atoi(options.DB)
	opts := redis.Options{
		Addr:     fmt.Sprintf("%s:%s", options.Host, options.Port),
		Password: options.Password,
		DB:       dbOpt,
	}
	client := redis.NewClient(&opts)

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis client with ttl %s failed to ping address %s, error is: %s",
			opts.IdleTimeout, opts.Addr, err)
	}

	return client, nil
}

func (rc *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return rc.client.Set(ctx, key, value, expiration).Err()
}

func (rc *RedisClient) Get(ctx context.Context, key string) ([]byte, error) {
	return rc.client.Get(ctx, key).Bytes()
}
