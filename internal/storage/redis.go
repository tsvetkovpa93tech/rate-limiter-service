package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/yourusername/rate-limiter-service/pkg/config"
)

// RedisStorage implements Redis storage for distributed rate limiting
// Suitable for multi-instance deployments
type RedisStorage struct {
	client *redis.Client
	logger *slog.Logger
}

// NewRedisStorage creates a new Redis storage instance
func NewRedisStorage(cfg config.StorageConfig, logger *slog.Logger) (*RedisStorage, error) {
	if logger == nil {
		logger = slog.Default()
	}

	poolSize := cfg.RedisPoolSize
	if poolSize == 0 {
		poolSize = 10
	}
	minIdle := cfg.RedisMinIdle
	if minIdle == 0 {
		minIdle = 5
	}

	client := redis.NewClient(&redis.Options{
		Addr:         cfg.RedisAddress,
		Password:     cfg.RedisPassword,
		DB:           cfg.RedisDB,
		PoolSize:     poolSize,              // Connection pool size
		MinIdleConns: minIdle,               // Minimum idle connections
		MaxRetries:   3,                     // Retry failed commands
		DialTimeout:  5 * time.Second,      // Connection timeout
		ReadTimeout:  3 * time.Second,      // Read timeout
		WriteTimeout: 3 * time.Second,      // Write timeout
		PoolTimeout:  4 * time.Second,      // Pool timeout
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		logger.Error("Failed to connect to Redis", "error", err, "address", cfg.RedisAddress)
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	logger.Info("Connected to Redis", "address", cfg.RedisAddress, "db", cfg.RedisDB)

	return &RedisStorage{
		client: client,
		logger: logger,
	}, nil
}

// Get retrieves a value from Redis storage
func (r *RedisStorage) Get(ctx context.Context, key string) (interface{}, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		r.logger.Debug("Key not found in Redis", "key", key)
		return nil, nil
	}
	if err != nil {
		r.logger.Error("Failed to get value from Redis", "key", key, "error", err)
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
func (r *RedisStorage) Set(ctx context.Context, key string, value interface{}, expiration int64) error {
	// Marshal value to JSON if it's not a string
	var val string
	switch v := value.(type) {
	case string:
		val = v
	default:
		data, err := json.Marshal(value)
		if err != nil {
			r.logger.Error("Failed to marshal value", "key", key, "error", err)
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

	err := r.client.Set(ctx, key, val, expirationDuration).Err()
	if err != nil {
		r.logger.Error("Failed to set value in Redis", "key", key, "error", err)
		return fmt.Errorf("failed to set value in Redis: %w", err)
	}

	r.logger.Debug("Value set in Redis", "key", key, "expiration", expirationDuration)
	return nil
}

// Delete removes a value from Redis storage
func (r *RedisStorage) Delete(ctx context.Context, key string) error {
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		r.logger.Error("Failed to delete value from Redis", "key", key, "error", err)
		return fmt.Errorf("failed to delete value from Redis: %w", err)
	}

	r.logger.Debug("Value deleted from Redis", "key", key)
	return nil
}

// Close closes the Redis connection
func (r *RedisStorage) Close() error {
	if err := r.client.Close(); err != nil {
		r.logger.Error("Failed to close Redis connection", "error", err)
		return err
	}
	r.logger.Info("Redis connection closed")
	return nil
}

