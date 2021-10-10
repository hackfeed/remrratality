package cacherepo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hackfeed/remrratality/backend/internal/db/cache"
	"github.com/hackfeed/remrratality/backend/internal/domain"
)

type RedisRepo struct {
	TTL         time.Duration
	CacheClient cache.RedisClient
}

func NewRedisRepo(cacheClient cache.RedisClient, ttl time.Duration) CacheRepository {
	return &RedisRepo{
		TTL:         ttl,
		CacheClient: cacheClient,
	}
}

func (rr *RedisRepo) GetMRR(key string) (domain.TotalMRR, error) {
	bytes, err := rr.CacheClient.Get(context.Background(), key)
	if err == redis.Nil {
		return domain.TotalMRR{}, nil
	}
	if err != redis.Nil && err != nil {
		return domain.TotalMRR{}, fmt.Errorf("failed to get mrr from cache by key %s, error is: %s", key, err)
	}

	var mrr domain.TotalMRR

	if err := json.Unmarshal(bytes, &mrr); err != nil {
		return domain.TotalMRR{}, fmt.Errorf("failed to unmarshal by key %s, error is: %s", key, err)
	}

	return mrr, nil
}

func (rr *RedisRepo) SetMRR(key string, mrr domain.TotalMRR) (domain.TotalMRR, error) {
	bytes, err := json.Marshal(mrr)
	if err != nil {
		return domain.TotalMRR{}, fmt.Errorf("failed to marshal by key %s, error is: %s", key, err)
	}

	if err := rr.CacheClient.Set(context.Background(), key, bytes, rr.TTL); err != nil {
		return domain.TotalMRR{}, fmt.Errorf("failed to set mrr to cache by key %s, error is: %s", key, err)
	}

	return mrr, nil
}
