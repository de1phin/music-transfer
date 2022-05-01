package redis

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type Redis[Key comparable, T any] struct {
	redis     *redis.Client
	expiresIn time.Duration
	ctx       context.Context
}

func NewRedis[Key comparable, T any](config Config, DB int, expiresIn time.Duration) *Redis[Key, T] {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       DB,
	})
	return &Redis[Key, T]{
		redis:     rdb,
		expiresIn: expiresIn,
		ctx:       context.Background(),
	}
}

func (r *Redis[Key, T]) Get(key int64) (result T, err error) {
	jsonValue, err := r.redis.Get(r.ctx, strconv.FormatInt(key, 10)).Result()
	if err != nil {
		return result, err
	}
	err = json.Unmarshal([]byte(jsonValue), &result)
	return result, nil
}

func (r *Redis[Key, T]) Exist(key int64) (bool, error) {
	result, err := r.redis.Exists(r.ctx, strconv.FormatInt(key, 10)).Result()
	return result == 1, err
}

func (r *Redis[Key, T]) Set(key int64, value T) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.redis.Set(r.ctx, strconv.FormatInt(key, 10), jsonValue, r.expiresIn).Err()
}
