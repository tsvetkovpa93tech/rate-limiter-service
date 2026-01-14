package services

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/tsvetkovpa93tech/rate-limiter-service/internal/storage"
)

func TestTokenBucketLimiter_Allow(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	memStorage := storage.NewMemoryStorage(logger)
	limiter := NewTokenBucketLimiter(memStorage, 5, time.Second, logger)
	ctx := context.Background()

	key := "test-key"

	// Should allow first 5 requests
	for i := 0; i < 5; i++ {
		allowed, err := limiter.Allow(ctx, key)
		if err != nil {
			t.Fatalf("Allow failed: %v", err)
		}
		if !allowed {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}

	// 6th request should be denied
	allowed, err := limiter.Allow(ctx, key)
	if err != nil {
		t.Fatalf("Allow failed: %v", err)
	}
	if allowed {
		t.Error("6th request should be denied")
	}
}

func TestTokenBucketLimiter_TokenRefill(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	memStorage := storage.NewMemoryStorage(logger)
	// 2 tokens per second
	limiter := NewTokenBucketLimiter(memStorage, 2, time.Second, logger)
	ctx := context.Background()

	key := "refill-key"

	// Consume all tokens
	for i := 0; i < 2; i++ {
		allowed, err := limiter.Allow(ctx, key)
		if err != nil {
			t.Fatalf("Allow failed: %v", err)
		}
		if !allowed {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}

	// Next request should be denied
	allowed, err := limiter.Allow(ctx, key)
	if err != nil {
		t.Fatalf("Allow failed: %v", err)
	}
	if allowed {
		t.Error("Request should be denied when no tokens available")
	}

	// Wait for token refill (slightly more than 0.5 seconds for 1 token)
	time.Sleep(600 * time.Millisecond)

	// Should allow one more request
	allowed, err = limiter.Allow(ctx, key)
	if err != nil {
		t.Fatalf("Allow failed: %v", err)
	}
	if !allowed {
		t.Error("Request should be allowed after token refill")
	}
}

func TestTokenBucketLimiter_ContextCancellation(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	memStorage := storage.NewMemoryStorage(logger)
	limiter := NewTokenBucketLimiter(memStorage, 5, time.Second, logger)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := limiter.Allow(ctx, "key")
	if err != context.Canceled {
		t.Errorf("Expected context.Canceled, got %v", err)
	}
}

func TestTokenBucketLimiter_DifferentKeys(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	memStorage := storage.NewMemoryStorage(logger)
	limiter := NewTokenBucketLimiter(memStorage, 2, time.Second, logger)
	ctx := context.Background()

	key1 := "key1"
	key2 := "key2"

	// Consume all tokens for key1
	allowed, err := limiter.Allow(ctx, key1)
	if err != nil {
		t.Fatalf("Allow failed: %v", err)
	}
	if !allowed {
		t.Error("Request should be allowed")
	}

	allowed, err = limiter.Allow(ctx, key1)
	if err != nil {
		t.Fatalf("Allow failed: %v", err)
	}
	if !allowed {
		t.Error("Request should be allowed")
	}

	// key1 should be exhausted
	allowed, err = limiter.Allow(ctx, key1)
	if err != nil {
		t.Fatalf("Allow failed: %v", err)
	}
	if allowed {
		t.Error("Request should be denied for key1")
	}

	// key2 should still have tokens
	allowed, err = limiter.Allow(ctx, key2)
	if err != nil {
		t.Fatalf("Allow failed: %v", err)
	}
	if !allowed {
		t.Error("Request should be allowed for key2")
	}
}
