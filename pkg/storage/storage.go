package storage

import (
	"fmt"

	"github.com/yourusername/rate-limiter-service/internal/config"
)

// Storage defines the interface for rate limiter storage
type Storage interface {
	// Get retrieves a value from storage
	Get(key string) (interface{}, error)
	// Set stores a value in storage with optional expiration
	Set(key string, value interface{}, expiration int64) error
	// Delete removes a value from storage
	Delete(key string) error
	// Close closes the storage connection
	Close() error
}

// NewStorage creates a new storage instance based on the storage type
func NewStorage(storageType string, cfg config.StorageConfig) (Storage, error) {
	switch storageType {
	case "memory":
		return NewMemoryStorage(), nil
	case "redis":
		return NewRedisStorage(cfg)
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", storageType)
	}
}

