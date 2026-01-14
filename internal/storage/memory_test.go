package storage

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"
)

func TestMemoryStorage_GetSetDelete(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	storage := NewMemoryStorage(logger)
	ctx := context.Background()

	key := "test-key"
	value := "test-value"

	// Test Set
	err := storage.Set(ctx, key, value, 0)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Test Get
	result, err := storage.Get(ctx, key)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if result != value {
		t.Errorf("Expected %v, got %v", value, result)
	}

	// Test Delete
	err = storage.Delete(ctx, key)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify deletion
	result, err = storage.Get(ctx, key)
	if err != nil {
		t.Fatalf("Get after delete failed: %v", err)
	}
	if result != nil {
		t.Errorf("Expected nil after delete, got %v", result)
	}
}

func TestMemoryStorage_Expiration(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	storage := NewMemoryStorage(logger)
	ctx := context.Background()

	key := "expiring-key"
	value := "expiring-value"
	// Expiration is stored as Unix timestamp in seconds, so we use
	// second-level precision in the test to avoid flakiness.
	expiration := time.Now().Add(1 * time.Second).Unix()

	// Set with expiration
	err := storage.Set(ctx, key, value, expiration)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Get before expiration (immediately after Set)
	result, err := storage.Get(ctx, key)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if result != value {
		t.Errorf("Expected %v, got %v", value, result)
	}

	// Wait for expiration (slightly more than 1 second)
	time.Sleep(1100 * time.Millisecond)

	// Get after expiration
	result, err = storage.Get(ctx, key)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if result != nil {
		t.Errorf("Expected nil after expiration, got %v", result)
	}
}

func TestMemoryStorage_ContextCancellation(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	storage := NewMemoryStorage(logger)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Test Get with cancelled context
	_, err := storage.Get(ctx, "key")
	if err != context.Canceled {
		t.Errorf("Expected context.Canceled, got %v", err)
	}

	// Test Set with cancelled context
	err = storage.Set(ctx, "key", "value", 0)
	if err != context.Canceled {
		t.Errorf("Expected context.Canceled, got %v", err)
	}

	// Test Delete with cancelled context
	err = storage.Delete(ctx, "key")
	if err != context.Canceled {
		t.Errorf("Expected context.Canceled, got %v", err)
	}
}

func TestMemoryStorage_Close(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	storage := NewMemoryStorage(logger)

	err := storage.Close()
	if err != nil {
		t.Errorf("Close should not return error, got %v", err)
	}
}

