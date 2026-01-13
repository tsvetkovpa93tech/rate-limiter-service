package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Collector manages Prometheus metrics
type Collector struct {
	totalRequests    *prometheus.CounterVec
	blockedRequests  *prometheus.CounterVec
	allowedRequests  *prometheus.CounterVec
	deniedRequests   *prometheus.CounterVec
	limitCheckErrors *prometheus.CounterVec
	requestDuration  *prometheus.HistogramVec
}

// NewCollector creates a new metrics collector
func NewCollector() *Collector {
	return &Collector{
		totalRequests: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "rate_limiter_total_requests",
				Help: "Total number of requests",
			},
			[]string{"method", "endpoint", "status"},
		),
		blockedRequests: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "rate_limiter_blocked_requests",
				Help: "Total number of blocked requests (rate limit exceeded)",
			},
			[]string{"algorithm", "key"},
		),
		allowedRequests: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "rate_limiter_allowed_requests_total",
				Help: "Total number of allowed requests",
			},
			[]string{"algorithm"},
		),
		deniedRequests: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "rate_limiter_denied_requests_total",
				Help: "Total number of denied requests",
			},
			[]string{"algorithm"},
		),
		limitCheckErrors: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "rate_limiter_check_errors_total",
				Help: "Total number of rate limit check errors",
			},
			[]string{"algorithm"},
		),
		requestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "rate_limiter_request_duration_seconds",
				Help:    "Request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "endpoint", "status"},
		),
	}
}

// Register registers all metrics with Prometheus
func (c *Collector) Register() {
	prometheus.MustRegister(c.totalRequests)
	prometheus.MustRegister(c.blockedRequests)
	prometheus.MustRegister(c.allowedRequests)
	prometheus.MustRegister(c.deniedRequests)
	prometheus.MustRegister(c.limitCheckErrors)
	prometheus.MustRegister(c.requestDuration)
}

// IncTotalRequests increments the total requests counter
func (c *Collector) IncTotalRequests(method, endpoint, status string) {
	c.totalRequests.WithLabelValues(method, endpoint, status).Inc()
}

// IncBlockedRequests increments the blocked requests counter
func (c *Collector) IncBlockedRequests(algorithm, key string) {
	c.blockedRequests.WithLabelValues(algorithm, key).Inc()
}

// IncAllowedRequests increments the allowed requests counter
func (c *Collector) IncAllowedRequests(algorithm string) {
	c.allowedRequests.WithLabelValues(algorithm).Inc()
}

// IncDeniedRequests increments the denied requests counter
func (c *Collector) IncDeniedRequests(algorithm string) {
	c.deniedRequests.WithLabelValues(algorithm).Inc()
}

// IncLimitCheckErrors increments the error counter
func (c *Collector) IncLimitCheckErrors(algorithm string) {
	c.limitCheckErrors.WithLabelValues(algorithm).Inc()
}

// ObserveRequestDuration records the request duration
func (c *Collector) ObserveRequestDuration(duration time.Duration, method, endpoint, status string) {
	c.requestDuration.WithLabelValues(method, endpoint, status).Observe(duration.Seconds())
}

// Handler returns the HTTP handler for metrics endpoint
func (c *Collector) Handler() http.Handler {
	return promhttp.Handler()
}

