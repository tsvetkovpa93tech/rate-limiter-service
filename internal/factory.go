package internal

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/yourusername/rate-limiter-service/internal/services"
	"github.com/yourusername/rate-limiter-service/internal/storage"
	"github.com/yourusername/rate-limiter-service/pkg/interfaces"
)

// AlgorithmType represents the type of rate limiting algorithm
type AlgorithmType string

const (
	// AlgorithmTokenBucket represents the Token Bucket algorithm
	AlgorithmTokenBucket AlgorithmType = "token_bucket"
	// AlgorithmSlidingWindow represents the Sliding Window Log algorithm
	AlgorithmSlidingWindow AlgorithmType = "sliding_window"
)

// LimiterConfig holds configuration for creating a rate limiter
type LimiterConfig struct {
	Algorithm AlgorithmType
	Limit     int
	Window    time.Duration
	Storage   storage.Storage
	Logger    *slog.Logger
}

// NewRateLimiter creates a new rate limiter based on the algorithm type
func NewRateLimiter(config LimiterConfig) (interfaces.RateLimiter, error) {
	if config.Storage == nil {
		return nil, fmt.Errorf("storage is required")
	}

	if config.Limit <= 0 {
		return nil, fmt.Errorf("limit must be greater than 0, got %d", config.Limit)
	}

	if config.Window <= 0 {
		return nil, fmt.Errorf("window must be greater than 0, got %v", config.Window)
	}

	switch config.Algorithm {
	case AlgorithmTokenBucket:
		return services.NewTokenBucketLimiter(
			config.Storage,
			config.Limit,
			config.Window,
			config.Logger,
		), nil
	case AlgorithmSlidingWindow:
		return services.NewSlidingWindowLimiter(
			config.Storage,
			config.Limit,
			config.Window,
			config.Logger,
		), nil
	default:
		return nil, fmt.Errorf("unsupported algorithm: %s", config.Algorithm)
	}
}

