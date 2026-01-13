package handlers

import (
	"net/http"

	"github.com/go-chi/render"
)

// HealthHandler handles health check requests
type HealthHandler struct{}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Check handles GET /health
func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, map[string]interface{}{
		"status":  "ok",
		"service": "rate-limiter-service",
		"version": "1.0.0",
	})
}

