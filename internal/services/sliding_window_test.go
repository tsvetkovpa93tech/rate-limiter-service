package services

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/yourusername/rate-limiter-service/internal/storage"
)

func TestSlidingWindowLimiter_Allow(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	memStorage := storage.NewMemoryStorage(logger)
	limiter := NewSlidingWindowLimiter(memStorage, 5, time.Second, logger)
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

func TestSlidingWindowLimiter_WindowSliding(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	memStorage := storage.NewMemoryStorage(logger)
	limiter := NewSlidingWindowLimiter(memStorage, 3, 500*time.Millisecond, logger)
	ctx := context.Background()

	key := "sliding-key"

	// Make 3 requests
	for i := 0; i < 3; i++ {
		allowed, err := limiter.Allow(ctx, key)
		if err != nil {
			t.Fatalf("Allow failed: %v", err)
		}
		if !allowed {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}

	// 4th request should be denied
	allowed, err := limiter.Allow(ctx, key)
	if err != nil {
		t.Fatalf("Allow failed: %v", err)
	}
	if allowed {
		t.Error("4th request should be denied")
	}

	// Wait for window to slide (oldest request should expire)
	time.Sleep(600 * time.Millisecond)

	// Should allow one more request (oldest expired)
	allowed, err = limiter.Allow(ctx, key)
	if err != nil {
		t.Fatalf("Allow failed: %v", err)
	}
	if !allowed {
		t.Error("Request should be allowed after window slides")
	}
}

func TestSlidingWindowLimiter_ContextCancellation(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	memStorage := storage.NewMemoryStorage(logger)
	limiter := NewSlidingWindowLimiter(memStorage, 5, time.Second, logger)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := limiter.Allow(ctx, "key")
	if err != context.Canceled {
		t.Errorf("Expected context.Canceled, got %v", err)
	}
}

func TestSlidingWindowLimiter_DifferentKeys(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	memStorage := storage.NewMemoryStorage(logger)
	limiter := NewSlidingWindowLimiter(memStorage, 2, time.Second, logger)
	ctx := context.Background()

	key1 := "key1"
	key2 := "key2"

	// Consume all requests for key1
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

	// key2 should still have capacity
	allowed, err = limiter.Allow(ctx, key2)
	if err != nil {
		t.Fatalf("Allow failed: %v", err)
	}
	if !allowed {
		t.Error("Request should be allowed for key2")
	}
}

