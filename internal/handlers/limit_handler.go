package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"log/slog"

	"github.com/yourusername/rate-limiter-service/internal/service"
)

// LimitHandler handles rate limiting requests
type LimitHandler struct {
	service *service.RateLimiterService
	logger  *slog.Logger
}

// NewLimitHandler creates a new limit handler
func NewLimitHandler(svc *service.RateLimiterService, logger *slog.Logger) *LimitHandler {
	if logger == nil {
		logger = slog.Default()
	}
	return &LimitHandler{
		service: svc,
		logger:  logger,
	}
}

// CheckLimit handles POST /api/v1/limit-check
func (h *LimitHandler) CheckLimit(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	ctx := r.Context()

	var req service.CheckLimitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Invalid request body", "error", err, "remote_addr", r.RemoteAddr)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request body"})
		return
	}

	if req.Key == "" {
		h.logger.Warn("Missing key in request", "remote_addr", r.RemoteAddr)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Key is required"})
		return
	}

	response, err := h.service.CheckLimit(ctx, &req)
	duration := time.Since(start)

	if err != nil {
		h.logger.Error("Failed to check limit",
			"error", err,
			"key", req.Key,
			"duration_ms", duration.Milliseconds(),
		)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	statusCode := http.StatusOK
	if !response.Allowed {
		statusCode = http.StatusTooManyRequests
		h.logger.Info("Rate limit exceeded",
			"key", req.Key,
			"algorithm", req.Algorithm,
			"duration_ms", duration.Milliseconds(),
		)
	} else {
		h.logger.Debug("Request allowed",
			"key", req.Key,
			"algorithm", req.Algorithm,
			"duration_ms", duration.Milliseconds(),
		)
	}

	render.Status(r, statusCode)
	render.JSON(w, r, response)
}

