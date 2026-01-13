package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/yourusername/rate-limiter-service/internal/config"
)

// RedisStorage implements Redis storage for distributed rate limiting
type RedisStorage struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisStorage creates a new Redis storage instance
func NewRedisStorage(cfg config.StorageConfig) (*RedisStorage, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddress,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	ctx := context.Background()

	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisStorage{
		client: client,
		ctx:    ctx,
	}, nil
}

// Get retrieves a value from Redis storage
func (r *RedisStorage) Get(key string) (interface{}, error) {
	val, err := r.client.Get(r.ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get value from Redis: %w", err)
	}

	// Try to unmarshal as JSON, if fails return as string
	var result interface{}
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		return val, nil
	}
	return result, nil
}

// Set stores a value in Redis storage with optional expiration
func (r *RedisStorage) Set(key string, value interface{}, expiration int64) error {
	// Marshal value to JSON if it's not a string
	var val string
	switch v := value.(type) {
	case string:
		val = v
	default:
		data, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value: %w", err)
		}
		val = string(data)
	}

	var expirationDuration time.Duration
	if expiration > 0 {
		expirationDuration = time.Until(time.Unix(expiration, 0))
		if expirationDuration < 0 {
			expirationDuration = 0
		}
	}

	if expirationDuration > 0 {
		return r.client.Set(r.ctx, key, val, expirationDuration).Err()
	}
	return r.client.Set(r.ctx, key, val, 0).Err()
}

// Delete removes a value from Redis storage
func (r *RedisStorage) Delete(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

// Close closes the Redis connection
func (r *RedisStorage) Close() error {
	return r.client.Close()
}

