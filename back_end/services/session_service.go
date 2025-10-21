package services

import (
	"context"
	"strconv"
	"time"

	"github.com/Chabuduo04/mate/back_end/config"
	"github.com/redis/go-redis/v9"
)

type SessionStore interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
}

type RedisSessionStore struct {
	rdb *redis.Client
}

func NewRedisSessionStore() (*RedisSessionStore, error) {
	cfg := config.GetConfig()
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host + ":" + strconv.Itoa(cfg.Redis.Port),
		Password: cfg.Redis.Password,
	})
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return &RedisSessionStore{rdb: rdb}, nil
}

func (s *RedisSessionStore) Get(ctx context.Context, key string) (string, error) {
	val, err := s.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

func (s *RedisSessionStore) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return s.rdb.Set(ctx, key, value, ttl).Err()
}

// --------------- Memory fallback ----------------
type MemorySessionStore struct {
	data map[string]string
}

func NewMemorySessionStore() *MemorySessionStore {
	return &MemorySessionStore{data: make(map[string]string)}
}

func (m *MemorySessionStore) Get(ctx context.Context, key string) (string, error) {
	if v, ok := m.data[key]; ok {
		return v, nil
	}
	return "", nil
}

func (m *MemorySessionStore) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	m.data[key] = value
	// no TTL enforcement for simplicity in memory mode
	return nil
}
