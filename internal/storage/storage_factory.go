package storage

import (
	"fmt"
	"log/slog"

	"github.com/tsvetkovpa93tech/rate-limiter-service/pkg/config"
)

// NewStorage creates a new storage instance based on the storage type
func NewStorage(storageType string, cfg config.StorageConfig, logger *slog.Logger) (Storage, error) {
	if logger == nil {
		logger = slog.Default()
	}

	switch storageType {
	case "memory":
		return NewMemoryStorage(logger), nil
	case "redis":
		return NewRedisStorage(cfg, logger)
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", storageType)
	}
}
