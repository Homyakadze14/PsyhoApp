package redisrepo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Homyakadze14/PsyhoApp/AuthMicroservice/internal/usecase"
	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	redis *redis.Client
}

func NewRedisRepository(redis *redis.Client) *RedisRepository {
	return &RedisRepository{redis}
}

func (r *RedisRepository) Set(ctx context.Context, key string, value any, expTime time.Duration) error {
	const op = "repositories.RedisRepository.Set"

	p, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = r.redis.Set(ctx, key, p, expTime).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *RedisRepository) Get(ctx context.Context, key string, dest any) error {
	const op = "repositories.RedisRepository.Get"

	var value []byte
	err := r.redis.Get(ctx, key).Scan(&value)
	if err != nil {
		if err == redis.Nil {
			return usecase.ErrCacheNotFound
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return json.Unmarshal(value, dest)
}

func (r *RedisRepository) Del(ctx context.Context, key string) (res int64, err error) {
	const op = "repositories.RedisRepository.Del"
	res, err = r.redis.Del(ctx, key).Result()
	if err != nil {
		return res, fmt.Errorf("%s: %w", op, err)
	}
	return res, nil
}

func (r *RedisRepository) Expire(ctx context.Context, key string, expiration time.Duration) error {
	const op = "repositories.RedisRepository.Expire"

	_, err := r.redis.Expire(ctx, key, expiration).Result()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
