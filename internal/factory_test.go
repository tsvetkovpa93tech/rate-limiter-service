package internal

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/yourusername/rate-limiter-service/internal/storage"
)

func TestNewRateLimiter_TokenBucket(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	memStorage := storage.NewMemoryStorage(logger)

	limiter, err := NewRateLimiter(LimiterConfig{
		Algorithm: AlgorithmTokenBucket,
		Limit:     5,
		Window:    time.Second,
		Storage:   memStorage,
		Logger:    logger,
	})

	if err != nil {
		t.Fatalf("NewRateLimiter failed: %v", err)
	}

	if limiter == nil {
		t.Fatal("Limiter should not be nil")
	}

	// Test that it works
	ctx := context.Background()
	allowed, err := limiter.Allow(ctx, "test-key")
	if err != nil {
		t.Fatalf("Allow failed: %v", err)
	}
	if !allowed {
		t.Error("First request should be allowed")
	}
}

func TestNewRateLimiter_SlidingWindow(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	memStorage := storage.NewMemoryStorage(logger)

	limiter, err := NewRateLimiter(LimiterConfig{
		Algorithm: AlgorithmSlidingWindow,
		Limit:     5,
		Window:    time.Second,
		Storage:   memStorage,
		Logger:    logger,
	})

	if err != nil {
		t.Fatalf("NewRateLimiter failed: %v", err)
	}

	if limiter == nil {
		t.Fatal("Limiter should not be nil")
	}

	// Test that it works
	ctx := context.Background()
	allowed, err := limiter.Allow(ctx, "test-key")
	if err != nil {
		t.Fatalf("Allow failed: %v", err)
	}
	if !allowed {
		t.Error("First request should be allowed")
	}
}

func TestNewRateLimiter_InvalidAlgorithm(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	memStorage := storage.NewMemoryStorage(logger)

	_, err := NewRateLimiter(LimiterConfig{
		Algorithm: AlgorithmType("invalid"),
		Limit:     5,
		Window:    time.Second,
		Storage:   memStorage,
		Logger:    logger,
	})

	if err == nil {
		t.Error("Expected error for invalid algorithm")
	}
}

func TestNewRateLimiter_InvalidLimit(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	memStorage := storage.NewMemoryStorage(logger)

	_, err := NewRateLimiter(LimiterConfig{
		Algorithm: AlgorithmTokenBucket,
		Limit:     0,
		Window:    time.Second,
		Storage:   memStorage,
		Logger:    logger,
	})

	if err == nil {
		t.Error("Expected error for invalid limit")
	}
}

func TestNewRateLimiter_InvalidWindow(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	memStorage := storage.NewMemoryStorage(logger)

	_, err := NewRateLimiter(LimiterConfig{
		Algorithm: AlgorithmTokenBucket,
		Limit:     5,
		Window:    0,
		Storage:   memStorage,
		Logger:    logger,
	})

	if err == nil {
		t.Error("Expected error for invalid window")
	}
}

func TestNewRateLimiter_NilStorage(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	_, err := NewRateLimiter(LimiterConfig{
		Algorithm: AlgorithmTokenBucket,
		Limit:     5,
		Window:    time.Second,
		Storage:   nil,
		Logger:    logger,
	})

	if err == nil {
		t.Error("Expected error for nil storage")
	}
}

