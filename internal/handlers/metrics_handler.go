package handlers

import (
	"net/http"

	"github.com/yourusername/rate-limiter-service/internal/metrics"
)

// MetricsHandler handles metrics requests
type MetricsHandler struct {
	collector *metrics.Collector
}

// NewMetricsHandler creates a new metrics handler
func NewMetricsHandler(collector *metrics.Collector) *MetricsHandler {
	return &MetricsHandler{
		collector: collector,
	}
}

// Serve handles GET /metrics
func (h *MetricsHandler) Serve(w http.ResponseWriter, r *http.Request) {
	h.collector.Handler().ServeHTTP(w, r)
}

