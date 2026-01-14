package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/tsvetkovpa93tech/rate-limiter-service/internal"
	"github.com/tsvetkovpa93tech/rate-limiter-service/internal/metrics"
	"github.com/tsvetkovpa93tech/rate-limiter-service/internal/storage"
	"github.com/tsvetkovpa93tech/rate-limiter-service/pkg/config"
)

// RateLimiterService handles rate limiting logic
type RateLimiterService struct {
	storage          storage.Storage
	config           *config.Config
	metricsCollector *metrics.Collector
	logger           *slog.Logger
}

// NewRateLimiterService creates a new rate limiter service
func NewRateLimiterService(
	storage storage.Storage,
	cfg *config.Config,
	metricsCollector *metrics.Collector,
	logger *slog.Logger,
) *RateLimiterService {
	if logger == nil {
		logger = slog.Default()
	}
	return &RateLimiterService{
		storage:          storage,
		config:           cfg,
		metricsCollector: metricsCollector,
		logger:           logger,
	}
}

// CheckLimitRequest represents a request to check rate limit
type CheckLimitRequest struct {
	Key       string `json:"key"`
	Algorithm string `json:"algorithm,omitempty"` // Optional: overrides default
	Limit     int    `json:"limit,omitempty"`     // Optional: overrides default
	Window    string `json:"window,omitempty"`    // Optional: overrides default (e.g., "1m", "30s")
}

// CheckLimitResponse represents the response from rate limit check
type CheckLimitResponse struct {
	Allowed   bool   `json:"allowed"`
	Remaining int    `json:"remaining,omitempty"`
	ResetAt   int64  `json:"reset_at,omitempty"`
	Message   string `json:"message,omitempty"`
}

// CheckLimit checks if a request should be allowed based on rate limiting rules
func (s *RateLimiterService) CheckLimit(ctx context.Context, req *CheckLimitRequest) (*CheckLimitResponse, error) {
	// Determine algorithm
	algorithmStr := req.Algorithm
	if algorithmStr == "" {
		algorithmStr = s.config.Limiter.DefaultAlgorithm
	}

	// Convert string to AlgorithmType
	var algorithm internal.AlgorithmType
	switch algorithmStr {
	case "token_bucket":
		algorithm = internal.AlgorithmTokenBucket
	case "sliding_window":
		algorithm = internal.AlgorithmSlidingWindow
	default:
		return nil, fmt.Errorf("unsupported algorithm: %s", algorithmStr)
	}

	// Determine limit
	limit := req.Limit
	if limit == 0 {
		limit = s.config.Limiter.DefaultLimit
	}

	// Determine window
	window := s.config.Limiter.DefaultWindow
	if req.Window != "" {
		parsedWindow, err := time.ParseDuration(req.Window)
		if err != nil {
			return nil, fmt.Errorf("invalid window duration: %w", err)
		}
		window = parsedWindow
	}

	// Create limiter using factory
	limiterInstance, err := internal.NewRateLimiter(internal.LimiterConfig{
		Algorithm: algorithm,
		Limit:     limit,
		Window:    window,
		Storage:   s.storage,
		Logger:    s.logger,
	})
	if err != nil {
		s.logger.Error("Failed to create limiter", "error", err, "algorithm", algorithm)
		return nil, fmt.Errorf("failed to create limiter: %w", err)
	}

	// Check limit
	allowed, err := limiterInstance.Allow(ctx, req.Key)
	if err != nil {
		s.metricsCollector.IncLimitCheckErrors(algorithmStr)
		s.logger.Error("Failed to check limit", "error", err, "key", req.Key)
		return nil, fmt.Errorf("failed to check limit: %w", err)
	}

	// Update metrics
	if allowed {
		s.metricsCollector.IncAllowedRequests(algorithmStr)
	} else {
		s.metricsCollector.IncDeniedRequests(algorithmStr)
		s.metricsCollector.IncBlockedRequests(algorithmStr, req.Key)
	}

	// Calculate reset time
	resetAt := time.Now().Add(window).Unix()

	response := &CheckLimitResponse{
		Allowed: allowed,
		ResetAt: resetAt,
	}

	if !allowed {
		response.Message = "Rate limit exceeded"
	}

	return response, nil
}
