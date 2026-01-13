package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/yourusername/rate-limiter-service/internal/storage"
	"github.com/yourusername/rate-limiter-service/pkg/interfaces"
)

// SlidingWindowLimiter implements the Sliding Window Log algorithm
// Tracks individual request timestamps for precise rate limiting
type SlidingWindowLimiter struct {
	storage storage.Storage
	limit   int
	window  time.Duration
	logger  *slog.Logger
}

// NewSlidingWindowLimiter creates a new Sliding Window Log limiter
func NewSlidingWindowLimiter(
	storage storage.Storage,
	limit int,
	window time.Duration,
	logger *slog.Logger,
) *SlidingWindowLimiter {
	if logger == nil {
		logger = slog.Default()
	}
	return &SlidingWindowLimiter{
		storage: storage,
		limit:   limit,
		window:  window,
		logger:  logger,
	}
}

// slidingWindowState represents the state of a sliding window
type slidingWindowState struct {
	Timestamps []int64 `json:"timestamps"` // Unix timestamps
}

// Allow checks if a request should be allowed using Sliding Window Log algorithm
func (s *SlidingWindowLimiter) Allow(ctx context.Context, key string) (bool, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		s.logger.Warn("Allow operation cancelled", "key", key, "error", ctx.Err())
		return false, ctx.Err()
	default:
	}

	now := time.Now()
	nowUnix := now.Unix()
	windowStart := now.Add(-s.window).Unix()

	// Get current state
	stateData, err := s.storage.Get(ctx, key)
	if err != nil {
		s.logger.Error("Failed to get window state", "key", key, "error", err)
		return false, fmt.Errorf("failed to get window state: %w", err)
	}

	var state slidingWindowState
	if stateData == nil {
		state = slidingWindowState{
			Timestamps: []int64{},
		}
		s.logger.Debug("Initialized new sliding window", "key", key)
	} else {
		// Parse state
		var stateJSON string
		switch v := stateData.(type) {
		case string:
			stateJSON = v
		default:
			data, err := json.Marshal(v)
			if err != nil {
				s.logger.Error("Failed to marshal state", "key", key, "error", err)
				return false, fmt.Errorf("failed to marshal state: %w", err)
			}
			stateJSON = string(data)
		}

		if err := json.Unmarshal([]byte(stateJSON), &state); err != nil {
			// If unmarshal fails, reset window
			s.logger.Warn("Failed to unmarshal state, resetting window", "key", key, "error", err)
			state = slidingWindowState{
				Timestamps: []int64{},
			}
		}
	}

	// Remove timestamps outside the window
	validTimestamps := make([]int64, 0, len(state.Timestamps))
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
		expiration := nowUnix + int64(s.window.Seconds())
		if err := s.storage.Set(ctx, key, string(stateJSON), expiration); err != nil {
			s.logger.Error("Failed to save window state", "key", key, "error", err)
			return false, fmt.Errorf("failed to save window state: %w", err)
		}
		s.logger.Debug("Request denied: limit exceeded",
			"key", key,
			"count", len(state.Timestamps),
			"limit", s.limit,
		)
		return false, nil
	}

	// Add current timestamp
	state.Timestamps = append(state.Timestamps, nowUnix)
	stateJSON, _ := json.Marshal(state)
	expiration := nowUnix + int64(s.window.Seconds())
	if err := s.storage.Set(ctx, key, string(stateJSON), expiration); err != nil {
		s.logger.Error("Failed to save window state", "key", key, "error", err)
		return false, fmt.Errorf("failed to save window state: %w", err)
	}

	s.logger.Debug("Request allowed",
		"key", key,
		"count", len(state.Timestamps),
		"limit", s.limit,
	)
	return true, nil
}

// Ensure SlidingWindowLimiter implements interfaces.RateLimiter
var _ interfaces.RateLimiter = (*SlidingWindowLimiter)(nil)

