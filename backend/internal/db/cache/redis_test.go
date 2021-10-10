package cache

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
)

var (
	redisTestClient *RedisClient
	ctx             = context.TODO()
)

func TestSet(t *testing.T) {
	db, mock := redismock.NewClientMock()
	redisTestClient = &RedisClient{
		Client: db,
	}

	mock.ExpectSet("key", "val", 1*time.Minute).SetVal("val")

	err := redisTestClient.Set(ctx, "key", "val", 1*time.Minute)
	if err != nil {
		t.Error(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestGet(t *testing.T) {
	db, mock := redismock.NewClientMock()
	redisTestClient = &RedisClient{
		Client: db,
	}

	mock.ExpectGet("key").SetVal("key")

	res, err := redisTestClient.Get(ctx, "key")
	if err != nil {
		t.Error(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}

	assert.Equal(t, []byte("key"), res)
}
