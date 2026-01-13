package services

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/yourusername/rate-limiter-service/internal/storage"
)

func BenchmarkTokenBucket_Allow(b *testing.B) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	memStorage := storage.NewMemoryStorage(logger)
	limiter := NewTokenBucketLimiter(memStorage, 1000, time.Second, logger)
	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		key := "bench:token-bucket"
		for pb.Next() {
			_, _ = limiter.Allow(ctx, key)
		}
	})
}

func BenchmarkSlidingWindow_Allow(b *testing.B) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	memStorage := storage.NewMemoryStorage(logger)
	limiter := NewSlidingWindowLimiter(memStorage, 1000, time.Second, logger)
	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		key := "bench:sliding-window"
		for pb.Next() {
			_, _ = limiter.Allow(ctx, key)
		}
	})
}

func BenchmarkTokenBucket_Allow_DifferentKeys(b *testing.B) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	memStorage := storage.NewMemoryStorage(logger)
	limiter := NewTokenBucketLimiter(memStorage, 100, time.Second, logger)
	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := "bench:token-bucket:" + string(rune(i%1000))
			_, _ = limiter.Allow(ctx, key)
			i++
		}
	})
}

func BenchmarkSlidingWindow_Allow_DifferentKeys(b *testing.B) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	memStorage := storage.NewMemoryStorage(logger)
	limiter := NewSlidingWindowLimiter(memStorage, 100, time.Second, logger)
	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := "bench:sliding-window:" + string(rune(i%1000))
			_, _ = limiter.Allow(ctx, key)
			i++
		}
	})
}

