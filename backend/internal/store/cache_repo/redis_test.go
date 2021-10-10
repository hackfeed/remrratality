package cacherepo

import (
	"testing"
	"time"

	"github.com/go-redis/redismock/v8"
	"github.com/hackfeed/remrratality/backend/internal/db/cache"
	"github.com/stretchr/testify/assert"
)

var (
	redisTestRepo *RedisRepo
)

func TestNewRedisRepo(t *testing.T) {
	db, _ := redismock.NewClientMock()

	redisTestClient := cache.RedisClient{Client: db}
	testTTL := 1 * time.Minute

	repo := NewRedisRepo(redisTestClient, testTTL)

	assert.NotNil(t, repo)
}
