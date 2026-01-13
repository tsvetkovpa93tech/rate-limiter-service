package storage

import (
	"context"
)

// Storage defines the interface for rate limiter storage backends
type Storage interface {
	// Get retrieves a value from storage
	Get(ctx context.Context, key string) (interface{}, error)
	// Set stores a value in storage with optional expiration (Unix timestamp)
	Set(ctx context.Context, key string, value interface{}, expiration int64) error
	// Delete removes a value from storage
	Delete(ctx context.Context, key string) error
	// Close closes the storage connection
	Close() error
}

