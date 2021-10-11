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

	tests := []struct {
		input string
		want  domain.TotalMRR
	}{
		{input: "notExistingKey", want: domain.TotalMRR{}},
		{input: "existingKeyWithRedisErr", want: domain.TotalMRR{}},
		{input: "existingKeyWithMarshalErr", want: domain.TotalMRR{}},
		{input: "existingKey", want: domain.TotalMRR{
			New:          []float32{0},
			Old:          []float32{0},
			Reactivation: []float32{0},
			Expansion:    []float32{0},
			Contraction:  []float32{0},
			Churn:        []float32{0},
			Total:        []float32{0},
		}},
	}

	for i := range tests {
		if tests[i].input == "notExistingKey" {
			mock.ExpectGet(tests[i].input).RedisNil()
			mrr, err := repo.GetMRR(tests[i].input)
			assert.Equal(t, tests[i].want, mrr)
			assert.NoError(t, err)
			if err = mock.ExpectationsWereMet(); err != nil {
				assert.Error(t, err)
			}
			mock.ClearExpect()
		}
		if tests[i].input == "existingKeyWithRedisErr" {
			redisErr := errors.New("redis err")
			mock.ExpectGet(tests[i].input).SetErr(redisErr)
			mrr, err := repo.GetMRR(tests[i].input)
			assert.Equal(t, tests[i].want, mrr)
			assert.Errorf(t, err, "failed to get mrr from cache by key %s, error is: %s", tests[i].input, err)
			if err = mock.ExpectationsWereMet(); err != nil {
				assert.Error(t, err)
			}
			mock.ClearExpect()
		}
		if tests[i].input == "existingKeyWithMarshalErr" {
			mock.ExpectGet(tests[i].input).SetVal("brokenJSON")
			mrr, err := repo.GetMRR(tests[i].input)
			assert.Equal(t, tests[i].want, mrr)
			assert.Errorf(t, err, "failed to unmarshal by key %s, error is: %s", tests[i].input, err)
			if err = mock.ExpectationsWereMet(); err != nil {
				assert.Error(t, err)
			}
			mock.ClearExpect()
		}
		if tests[i].input == "existingKey" {
			mock.ExpectGet(tests[i].input).SetVal("{\"New\":[0],\"Old\":[0],\"Reactivation\":[0],\"Expansion\":[0],\"Contraction\":[0],\"Churn\":[0],\"Total\":[0]}")
			mrr, err := repo.GetMRR(tests[i].input)
			assert.Equal(t, tests[i].want, mrr)
			assert.NoError(t, err)
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

	tests := []struct {
		inputKey   string
		inputValue domain.TotalMRR
		want       domain.TotalMRR
	}{
		{inputKey: "keyWithRedisErr", inputValue: domain.TotalMRR{}, want: domain.TotalMRR{}},
		{
			inputKey: "key",
			inputValue: domain.TotalMRR{
				New:          []float32{0},
				Old:          []float32{0},
				Reactivation: []float32{0},
				Expansion:    []float32{0},
				Contraction:  []float32{0},
				Churn:        []float32{0},
				Total:        []float32{0},
			},
			want: domain.TotalMRR{
				New:          []float32{0},
				Old:          []float32{0},
				Reactivation: []float32{0},
				Expansion:    []float32{0},
				Contraction:  []float32{0},
				Churn:        []float32{0},
				Total:        []float32{0},
			},
		},
	}

	for i := range tests {
		if tests[i].inputKey == "keyWithRedisErr" {
			redisErr := errors.New("redis err")
			mock.ExpectSet(tests[i].inputKey, tests[i].inputValue, testTTL).SetErr(redisErr)
			mrr, err := repo.SetMRR(tests[i].inputKey, tests[i].inputValue)
			assert.Equal(t, tests[i].want, mrr)
			assert.Errorf(t, err, "failed to set mrr to cache by key %s, error is: %s", tests[i].inputKey, err)
			if err = mock.ExpectationsWereMet(); err != nil {
				assert.Error(t, err)
			}
			mock.ClearExpect()
		}
		if tests[i].inputKey == "key" {
			bytes, _ := json.Marshal(tests[i].inputValue)
			mock.ExpectSet(tests[i].inputKey, bytes, testTTL).SetVal("")
			mrr, err := repo.SetMRR(tests[i].inputKey, tests[i].inputValue)
			assert.Equal(t, tests[i].want, mrr)
			assert.NoError(t, err)
			if err = mock.ExpectationsWereMet(); err != nil {
				assert.Error(t, err)
			}
			mock.ClearExpect()
		}
	}
}
