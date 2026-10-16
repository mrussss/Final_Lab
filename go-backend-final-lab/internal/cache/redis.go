package cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisStore(addr string, ttl time.Duration) *RedisStore {
	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		DialTimeout:  250 * time.Millisecond,
		ReadTimeout:  250 * time.Millisecond,
		WriteTimeout: 250 * time.Millisecond,
		MaxRetries:   -1,
	})

	return &RedisStore{
		client: client,
		ttl:    ttl,
	}
}

func (s *RedisStore) Get(
	ctx context.Context,
	key string,
) (string, bool, error) {
	opCtx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	value, err := s.client.Get(opCtx, key).Result()

	if errors.Is(err, redis.Nil) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}

	return value, true, nil
}

func (s *RedisStore) Set(
	ctx context.Context,
	key string,
	value string,
) error {
	opCtx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	return s.client.Set(opCtx, key, value, s.ttl).Err()
}

func (s *RedisStore) Delete(
	ctx context.Context,
	key string,
) error {
	opCtx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	return s.client.Del(opCtx, key).Err()
}

func (s *RedisStore) Ping(ctx context.Context) error {
	opCtx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	return s.client.Ping(opCtx).Err()
}

func (s *RedisStore) Close() error {
	return s.client.Close()
}
