package cacherepo

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hackfeed/remrratality/backend/internal/db/cache"
	"github.com/hackfeed/remrratality/backend/internal/domain"
)

type redisRepo struct {
	ttl         time.Duration
	cacheClient cache.RedisClient
}

func NewRedisRepo(cacheClient cache.RedisClient, ttl time.Duration) CacheRepository {
	return &redisRepo{
		ttl:         ttl,
		cacheClient: cacheClient,
	}
}

func (rr *redisRepo) GetMRR(key string) (domain.TotalMRR, error) {
	bytes, err := rr.cacheClient.Get(context.Background(), key)
	if err == redis.Nil {
		return domain.TotalMRR{}, nil
	}
	if err != redis.Nil && err != nil {
		return domain.TotalMRR{}, err
	}

	var mrr domain.TotalMRR

	if err := json.Unmarshal(bytes, &mrr); err != nil {
		return domain.TotalMRR{}, err
	}

	return mrr, nil
}

func (rr *redisRepo) SetMRR(key string, mrr domain.TotalMRR) (domain.TotalMRR, error) {
	bytes, err := json.Marshal(mrr)
	if err != nil {
		return domain.TotalMRR{}, err
	}

	if err := rr.cacheClient.Set(context.Background(), key, bytes, rr.ttl); err != nil {
		return domain.TotalMRR{}, err
	}

	return mrr, nil
}