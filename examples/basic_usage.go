package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/yourusername/rate-limiter-service/internal"
	"github.com/yourusername/rate-limiter-service/internal/storage"
)

func main() {
	// Initialize logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// Create in-memory storage
	memStorage := storage.NewMemoryStorage(logger)

	// Create Token Bucket limiter using factory
	tokenBucketLimiter, err := internal.NewRateLimiter(internal.LimiterConfig{
		Algorithm: internal.AlgorithmTokenBucket,
		Limit:     5,
		Window:    time.Second,
		Storage:   memStorage,
		Logger:    logger,
	})
	if err != nil {
		logger.Error("Failed to create limiter", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()
	key := "user:123"

	// Test rate limiting
	fmt.Println("Testing Token Bucket Limiter:")
	for i := 0; i < 7; i++ {
		allowed, err := tokenBucketLimiter.Allow(ctx, key)
		if err != nil {
			logger.Error("Error checking limit", "error", err)
			continue
		}
		if allowed {
			fmt.Printf("Request %d: ALLOWED\n", i+1)
		} else {
			fmt.Printf("Request %d: DENIED (rate limit exceeded)\n", i+1)
		}
	}

	// Create Sliding Window limiter
	slidingWindowLimiter, err := internal.NewRateLimiter(internal.LimiterConfig{
		Algorithm: internal.AlgorithmSlidingWindow,
		Limit:     3,
		Window:    2 * time.Second,
		Storage:   memStorage,
		Logger:    logger,
	})
	if err != nil {
		logger.Error("Failed to create limiter", "error", err)
		os.Exit(1)
	}

	fmt.Println("\nTesting Sliding Window Limiter:")
	key2 := "user:456"
	for i := 0; i < 5; i++ {
		allowed, err := slidingWindowLimiter.Allow(ctx, key2)
		if err != nil {
			logger.Error("Error checking limit", "error", err)
			continue
		}
		if allowed {
			fmt.Printf("Request %d: ALLOWED\n", i+1)
		} else {
			fmt.Printf("Request %d: DENIED (rate limit exceeded)\n", i+1)
		}
		time.Sleep(300 * time.Millisecond)
	}

	// Cleanup
	memStorage.Close()
}
