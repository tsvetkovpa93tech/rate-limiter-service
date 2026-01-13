package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/render"

	"github.com/yourusername/rate-limiter-service/internal/service"
)

// RateLimiterHandler handles HTTP requests for rate limiting
type RateLimiterHandler struct {
	service *service.RateLimiterService
}

// NewRateLimiterHandler creates a new rate limiter handler
func NewRateLimiterHandler(svc *service.RateLimiterService) *RateLimiterHandler {
	return &RateLimiterHandler{
		service: svc,
	}
}

// CheckLimit handles POST /api/v1/limit-check
func (h *RateLimiterHandler) CheckLimit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req service.CheckLimitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request body"})
		return
	}

	if req.Key == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Key is required"})
		return
	}

	response, err := h.service.CheckLimit(ctx, &req)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	statusCode := http.StatusOK
	if !response.Allowed {
		statusCode = http.StatusTooManyRequests
	}

	render.Status(r, statusCode)
	render.JSON(w, r, response)
}

// HealthCheck handles GET /health
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, map[string]string{
		"status": "ok",
		"service": "rate-limiter-service",
	})
}

