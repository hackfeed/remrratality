package cacherepo

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/go-redis/redismock/v8"
	"github.com/hackfeed/remrratality/backend/internal/db/cache"
	"github.com/hackfeed/remrratality/backend/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewRedisRepo(t *testing.T) {
	db, _ := redismock.NewClientMock()

	redisTestClient := cache.RedisClient{Client: db}
	testTTL := 1 * time.Minute

	repo := NewRedisRepo(redisTestClient, testTTL)

	assert.NotNil(t, repo)
}

func TestGetMRR(t *testing.T) {
	db, mock := redismock.NewClientMock()

	redisTestClient := cache.RedisClient{Client: db}
	testTTL := 1 * time.Minute

	repo := NewRedisRepo(redisTestClient, testTTL)

	type testInput struct {
		key string
	}
	type testWant struct {
		mrr domain.TotalMRR
		err error
	}

	tests := []struct {
		input testInput
		want  testWant
	}{
		{
			input: testInput{
				key: "notExistingKey",
			},
			want: testWant{
				mrr: domain.TotalMRR{},
				err: nil,
			},
		},
		{
			input: testInput{
				key: "existingKeyWithRedisErr",
			},
			want: testWant{
				mrr: domain.TotalMRR{},
				err: errors.New("failed to get mrr from cache by key existingKeyWithRedisErr, error is: redis err"),
			},
		},
		{
			input: testInput{
				key: "existingKeyWithMarshalErr",
			},
			want: testWant{
				mrr: domain.TotalMRR{},
				err: errors.New("failed to unmarshal by key existingKeyWithMarshalErr, error is: invalid character 'b' looking for beginning of value"),
			},
		},
		{
			input: testInput{
				key: "existingKey",
			},
			want: testWant{
				mrr: domain.TotalMRR{
					New:          []float32{0},
					Old:          []float32{0},
					Reactivation: []float32{0},
					Expansion:    []float32{0},
					Contraction:  []float32{0},
					Churn:        []float32{0},
					Total:        []float32{0},
				},
				err: nil,
			},
		},
	}

	for _, test := range tests {
		if test.input.key == "notExistingKey" {
			mock.ExpectGet(test.input.key).RedisNil()
			mrr, err := repo.GetMRR(test.input.key)
			assert.Equal(t, test.want.mrr, mrr)
			assert.Equal(t, test.want.err, err)
			if err = mock.ExpectationsWereMet(); err != nil {
				assert.Error(t, err)
			}
			mock.ClearExpect()
		}
		if test.input.key == "existingKeyWithRedisErr" {
			redisErr := errors.New("redis err")
			mock.ExpectGet(test.input.key).SetErr(redisErr)
			mrr, err := repo.GetMRR(test.input.key)
			assert.Equal(t, test.want.mrr, mrr)
			assert.Equal(t, test.want.err, err)
			if err = mock.ExpectationsWereMet(); err != nil {
				assert.Error(t, err)
			}
			mock.ClearExpect()
		}
		if test.input.key == "existingKeyWithMarshalErr" {
			mock.ExpectGet(test.input.key).SetVal("brokenJSON")
			mrr, err := repo.GetMRR(test.input.key)
			assert.Equal(t, test.want.mrr, mrr)
			assert.Equal(t, test.want.err, err)
			if err = mock.ExpectationsWereMet(); err != nil {
				assert.Error(t, err)
			}
			mock.ClearExpect()
		}
		if test.input.key == "existingKey" {
			mock.ExpectGet(test.input.key).SetVal("{\"New\":[0],\"Old\":[0],\"Reactivation\":[0],\"Expansion\":[0],\"Contraction\":[0],\"Churn\":[0],\"Total\":[0]}")
			mrr, err := repo.GetMRR(test.input.key)
			assert.Equal(t, test.want.mrr, mrr)
			assert.Equal(t, test.want.err, err)
			if err = mock.ExpectationsWereMet(); err != nil {
				assert.Error(t, err)
			}
			mock.ClearExpect()
		}
	}
}

func TestSetMRR(t *testing.T) {
	db, mock := redismock.NewClientMock()

	redisTestClient := cache.RedisClient{Client: db}
	testTTL := 1 * time.Minute

	repo := NewRedisRepo(redisTestClient, testTTL)

	type testInput struct {
		key string
		mrr domain.TotalMRR
	}
	type testWant struct {
		mrr domain.TotalMRR
		err error
	}

	tests := []struct {
		input testInput
		want  testWant
	}{
		{
			input: testInput{
				key: "keyWithRedisErr",
				mrr: domain.TotalMRR{},
			},
			want: testWant{
				mrr: domain.TotalMRR{},
				err: errors.New("failed to set mrr to cache by key keyWithRedisErr, error is: redis err"),
			},
		},
		{
			input: testInput{
				key: "key",
				mrr: domain.TotalMRR{
					New:          []float32{0},
					Old:          []float32{0},
					Reactivation: []float32{0},
					Expansion:    []float32{0},
					Contraction:  []float32{0},
					Churn:        []float32{0},
					Total:        []float32{0},
				},
			},
			want: testWant{
				mrr: domain.TotalMRR{
					New:          []float32{0},
					Old:          []float32{0},
					Reactivation: []float32{0},
					Expansion:    []float32{0},
					Contraction:  []float32{0},
					Churn:        []float32{0},
					Total:        []float32{0},
				},
				err: nil,
			},
		},
	}

	for _, test := range tests {
		if test.input.key == "keyWithRedisErr" {
			redisErr := errors.New("redis err")
			bytes, _ := json.Marshal(test.input.mrr)
			mock.ExpectSet(test.input.key, bytes, testTTL).SetErr(redisErr)
			mrr, err := repo.SetMRR(test.input.key, test.input.mrr)
			assert.Equal(t, test.want.mrr, mrr)
			assert.Equal(t, test.want.err, err)
			if err = mock.ExpectationsWereMet(); err != nil {
				assert.Error(t, err)
			}
			mock.ClearExpect()
		}
		if test.input.key == "key" {
			bytes, _ := json.Marshal(test.input.mrr)
			mock.ExpectSet(test.input.key, bytes, testTTL).SetVal("")
			mrr, err := repo.SetMRR(test.input.key, test.input.mrr)
			assert.Equal(t, test.want.mrr, mrr)
			assert.Equal(t, test.want.err, err)
			if err = mock.ExpectationsWereMet(); err != nil {
				assert.Error(t, err)
			}
			mock.ClearExpect()
		}
	}
}
