package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/tsvetkovpa93tech/rate-limiter-service/internal/storage"
	"github.com/tsvetkovpa93tech/rate-limiter-service/pkg/interfaces"
)

// TokenBucketLimiter implements the Token Bucket algorithm
// Tokens are refilled at a constant rate, allowing for burst capacity
type TokenBucketLimiter struct {
	storage storage.Storage
	limit   int
	window  time.Duration
	logger  *slog.Logger
}

// NewTokenBucketLimiter creates a new Token Bucket limiter
func NewTokenBucketLimiter(
	storage storage.Storage,
	limit int,
	window time.Duration,
	logger *slog.Logger,
) *TokenBucketLimiter {
	if logger == nil {
		logger = slog.Default()
	}
	return &TokenBucketLimiter{
		storage: storage,
		limit:   limit,
		window:  window,
		logger:  logger,
	}
}

// tokenBucketState represents the state of a token bucket
type tokenBucketState struct {
	Tokens     int   `json:"tokens"`
	LastRefill int64 `json:"last_refill"` // Unix timestamp in nanoseconds
}

// Allow checks if a request should be allowed using Token Bucket algorithm
func (t *TokenBucketLimiter) Allow(ctx context.Context, key string) (bool, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		t.logger.Warn("Allow operation cancelled", "key", key, "error", ctx.Err())
		return false, ctx.Err()
	default:
	}

	now := time.Now()
	refillRate := float64(t.limit) / t.window.Seconds()

	// Get current state
	stateData, err := t.storage.Get(ctx, key)
	if err != nil {
		t.logger.Error("Failed to get bucket state", "key", key, "error", err)
		return false, fmt.Errorf("failed to get bucket state: %w", err)
	}

	var state tokenBucketState
	if stateData == nil {
		// Initialize bucket with full tokens
		state = tokenBucketState{
			Tokens:     t.limit,
			LastRefill: now.UnixNano(),
		}
		t.logger.Debug("Initialized new token bucket", "key", key, "tokens", t.limit)
	} else {
		// Parse state
		var stateJSON string
		switch v := stateData.(type) {
		case string:
			stateJSON = v
		default:
			data, err := json.Marshal(v)
			if err != nil {
				t.logger.Error("Failed to marshal state", "key", key, "error", err)
				return false, fmt.Errorf("failed to marshal state: %w", err)
			}
			stateJSON = string(data)
		}

		if err := json.Unmarshal([]byte(stateJSON), &state); err != nil {
			// If unmarshal fails, reset bucket
			t.logger.Warn("Failed to unmarshal state, resetting bucket", "key", key, "error", err)
			state = tokenBucketState{
				Tokens:     t.limit,
				LastRefill: now.UnixNano(),
			}
		} else {
			// Refill tokens based on time elapsed since last refill/emptying
			elapsedSeconds := float64(now.UnixNano()-state.LastRefill) / float64(time.Second)
			if elapsedSeconds > 0 {
				tokensToAdd := int(elapsedSeconds * refillRate)
				if tokensToAdd > 0 {
					oldTokens := state.Tokens
					state.Tokens = min(state.Tokens+tokensToAdd, t.limit)
					state.LastRefill = now.UnixNano()
					t.logger.Debug("Tokens refilled",
						"key", key,
						"old_tokens", oldTokens,
						"added", tokensToAdd,
						"new_tokens", state.Tokens,
					)
				}
			}
		}
	}

	// Check if we have tokens available
	if state.Tokens <= 0 {
		// Save state even if request is denied
		stateJSON, _ := json.Marshal(state)
		expiration := now.Add(t.window).Unix()
		if err := t.storage.Set(ctx, key, string(stateJSON), expiration); err != nil {
			t.logger.Error("Failed to save bucket state", "key", key, "error", err)
			return false, fmt.Errorf("failed to save bucket state: %w", err)
		}
		t.logger.Debug("Request denied: no tokens available", "key", key)
		return false, nil
	}

	// Consume a token
	state.Tokens--
	// When the bucket becomes empty, move the refill reference point
	// to the moment of emptying so that new tokens start accumulating
	// from this time, not from the initial creation time.
	if state.Tokens == 0 {
		state.LastRefill = now.UnixNano()
	}
	stateJSON, _ := json.Marshal(state)
	expiration := now.Add(t.window).Unix()
	if err := t.storage.Set(ctx, key, string(stateJSON), expiration); err != nil {
		t.logger.Error("Failed to save bucket state", "key", key, "error", err)
		return false, fmt.Errorf("failed to save bucket state: %w", err)
	}

	t.logger.Debug("Request allowed", "key", key, "remaining_tokens", state.Tokens)
	return true, nil
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Ensure TokenBucketLimiter implements interfaces.RateLimiter
var _ interfaces.RateLimiter = (*TokenBucketLimiter)(nil)
