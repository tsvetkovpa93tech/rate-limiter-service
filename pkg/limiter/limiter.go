package limiter

import (
	"fmt"
	"time"

	"github.com/tsvetkovpa93tech/rate-limiter-service/pkg/storage"
)

// Limiter defines the interface for rate limiting algorithms
type Limiter interface {
	// Allow checks if a request should be allowed
	Allow(key string) (bool, error)
	// Reset resets the limiter state for a key
	Reset(key string) error
}

// NewLimiter creates a new limiter instance based on the algorithm type
func NewLimiter(algorithm string, storage storage.Storage, limit int, window time.Duration) (Limiter, error) {
	switch algorithm {
	case "token_bucket":
		return NewTokenBucketLimiter(storage, limit, window), nil
	case "sliding_window":
		return NewSlidingWindowLimiter(storage, limit, window), nil
	default:
		return nil, fmt.Errorf("unsupported algorithm: %s", algorithm)
	}
}
