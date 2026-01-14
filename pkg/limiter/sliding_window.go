package limiter

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/tsvetkovpa93tech/rate-limiter-service/pkg/storage"
)

// SlidingWindowLimiter implements the Sliding Window Log algorithm
type SlidingWindowLimiter struct {
	storage storage.Storage
	limit   int
	window  time.Duration
}

// NewSlidingWindowLimiter creates a new Sliding Window Log limiter
func NewSlidingWindowLimiter(storage storage.Storage, limit int, window time.Duration) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		storage: storage,
		limit:   limit,
		window:  window,
	}
}

// slidingWindowState represents the state of a sliding window
type slidingWindowState struct {
	Timestamps []int64 `json:"timestamps"` // Unix timestamps
}

// Allow checks if a request should be allowed using Sliding Window Log algorithm
func (s *SlidingWindowLimiter) Allow(key string) (bool, error) {
	now := time.Now().Unix()
	windowStart := now - int64(s.window.Seconds())

	// Get current state
	stateData, err := s.storage.Get(key)
	if err != nil {
		return false, fmt.Errorf("failed to get window state: %w", err)
	}

	var state slidingWindowState
	if stateData == nil {
		state = slidingWindowState{
			Timestamps: []int64{},
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
			// If unmarshal fails, reset window
			state = slidingWindowState{
				Timestamps: []int64{},
			}
		}
	}

	// Remove timestamps outside the window
	validTimestamps := []int64{}
	for _, ts := range state.Timestamps {
		if ts > windowStart {
			validTimestamps = append(validTimestamps, ts)
		}
	}
	state.Timestamps = validTimestamps

	// Check if we're within the limit
	if len(state.Timestamps) >= s.limit {
		// Save state even if request is denied
		stateJSON, _ := json.Marshal(state)
		expiration := now + int64(s.window.Seconds())
		s.storage.Set(key, string(stateJSON), expiration)
		return false, nil
	}

	// Add current timestamp
	state.Timestamps = append(state.Timestamps, now)
	stateJSON, _ := json.Marshal(state)
	expiration := now + int64(s.window.Seconds())
	if err := s.storage.Set(key, string(stateJSON), expiration); err != nil {
		return false, fmt.Errorf("failed to save window state: %w", err)
	}

	return true, nil
}

// Reset resets the sliding window for a key
func (s *SlidingWindowLimiter) Reset(key string) error {
	return s.storage.Delete(key)
}
