package limiter

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/yourusername/rate-limiter-service/pkg/storage"
)

// TokenBucketLimiter implements the Token Bucket algorithm
type TokenBucketLimiter struct {
	storage storage.Storage
	limit   int
	window  time.Duration
}

// NewTokenBucketLimiter creates a new Token Bucket limiter
func NewTokenBucketLimiter(storage storage.Storage, limit int, window time.Duration) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		storage: storage,
		limit:   limit,
		window:  window,
	}
}

// tokenBucketState represents the state of a token bucket
type tokenBucketState struct {
	Tokens     int   `json:"tokens"`
	LastRefill int64 `json:"last_refill"` // Unix timestamp
}

// Allow checks if a request should be allowed using Token Bucket algorithm
func (t *TokenBucketLimiter) Allow(key string) (bool, error) {
	now := time.Now().Unix()
	refillRate := float64(t.limit) / t.window.Seconds()

	// Get current state
	stateData, err := t.storage.Get(key)
	if err != nil {
		return false, fmt.Errorf("failed to get bucket state: %w", err)
	}

	var state tokenBucketState
	if stateData == nil {
		// Initialize bucket with full tokens
		state = tokenBucketState{
			Tokens:     t.limit,
			LastRefill: now,
		}
	} else {
		// Parse state
		var stateJSON string
		switch v := stateData.(type) {
		case string:
			stateJSON = v
		default:
			data, err := json.Marshal(v)
			if err != nil {
				return false, fmt.Errorf("failed to marshal state: %w", err)
			}
			stateJSON = string(data)
		}

		if err := json.Unmarshal([]byte(stateJSON), &state); err != nil {
			// If unmarshal fails, reset bucket
			state = tokenBucketState{
				Tokens:     t.limit,
				LastRefill: now,
			}
		} else {
			// Refill tokens based on time elapsed
			elapsed := now - state.LastRefill
			if elapsed > 0 {
				tokensToAdd := int(float64(elapsed) * refillRate)
				if tokensToAdd > 0 {
					state.Tokens = min(state.Tokens+tokensToAdd, t.limit)
					state.LastRefill = now
				}
			}
		}
	}

	// Check if we have tokens available
	if state.Tokens <= 0 {
		// Save state even if request is denied
		stateJSON, _ := json.Marshal(state)
		expiration := now + int64(t.window.Seconds())
		t.storage.Set(key, string(stateJSON), expiration)
		return false, nil
	}

	// Consume a token
	state.Tokens--
	stateJSON, _ := json.Marshal(state)
	expiration := now + int64(t.window.Seconds())
	if err := t.storage.Set(key, string(stateJSON), expiration); err != nil {
		return false, fmt.Errorf("failed to save bucket state: %w", err)
	}

	return true, nil
}

// Reset resets the token bucket for a key
func (t *TokenBucketLimiter) Reset(key string) error {
	return t.storage.Delete(key)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
