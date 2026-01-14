package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/tsvetkovpa93tech/rate-limiter-service/internal/handlers"
	"github.com/tsvetkovpa93tech/rate-limiter-service/internal/metrics"
	appmw "github.com/tsvetkovpa93tech/rate-limiter-service/internal/middleware"
	"github.com/tsvetkovpa93tech/rate-limiter-service/internal/service"
	"github.com/tsvetkovpa93tech/rate-limiter-service/internal/storage"
	"github.com/tsvetkovpa93tech/rate-limiter-service/pkg/config"
)

func main() {
	// Initialize logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Initialize storage
	var storageInstance storage.Storage
	switch cfg.Storage.Type {
	case "memory":
		storageInstance = storage.NewMemoryStorage(logger)
	case "redis":
		storageInstance, err = storage.NewRedisStorage(cfg.Storage, logger)
		if err != nil {
			logger.Error("Failed to initialize Redis storage", "error", err)
			os.Exit(1)
		}
	default:
		logger.Error("Unsupported storage type", "type", cfg.Storage.Type)
		os.Exit(1)
	}
	defer storageInstance.Close()

	// Initialize metrics
	metricsCollector := metrics.NewCollector()
	metricsCollector.Register()

	// Initialize service
	rateLimiterService := service.NewRateLimiterService(storageInstance, cfg, metricsCollector, logger)

	// Initialize handlers
	limitHandler := handlers.NewLimitHandler(rateLimiterService, logger)
	healthHandler := handlers.NewHealthHandler()
	metricsHandler := handlers.NewMetricsHandler(metricsCollector)

	// Setup router
	router := chi.NewRouter()

	// Middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(appmw.RequestLogger(logger))
	router.Use(appmw.RecoveryMiddleware(logger))
	router.Use(appmw.CORS(cfg.CORS.AllowedOrigins))

	// Metrics middleware
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			duration := time.Since(start)
			status := strconv.Itoa(ww.Status())
			metricsCollector.IncTotalRequests(r.Method, r.URL.Path, status)
			metricsCollector.ObserveRequestDuration(duration, r.Method, r.URL.Path, status)
		})
	})

	// Routes
	router.Get("/health", healthHandler.Check)
	router.Get("/metrics", metricsHandler.Serve)
	router.Route("/api/v1", func(r chi.Router) {
		r.Post("/limit-check", limitHandler.CheckLimit)
	})

	// Start server
	readTimeout := cfg.Server.ReadTimeout
	if readTimeout == 0 {
		readTimeout = 15 * time.Second
	}
	writeTimeout := cfg.Server.WriteTimeout
	if writeTimeout == 0 {
		writeTimeout = 15 * time.Second
	}
	idleTimeout := cfg.Server.IdleTimeout
	if idleTimeout == 0 {
		idleTimeout = 60 * time.Second
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	// Graceful shutdown
	go func() {
		logger.Info("Server starting", "port", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	logger.Info("Server exited")
}
