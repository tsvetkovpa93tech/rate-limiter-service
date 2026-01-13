package storage

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

// MemoryStorage implements in-memory storage using sync.Map
// Suitable for single-instance deployments
type MemoryStorage struct {
	data   sync.Map
	logger *slog.Logger
}

// NewMemoryStorage creates a new in-memory storage instance
func NewMemoryStorage(logger *slog.Logger) *MemoryStorage {
	if logger == nil {
		logger = slog.Default()
	}
	return &MemoryStorage{
		data:   sync.Map{},
		logger: logger,
	}
}

// Get retrieves a value from memory storage
func (m *MemoryStorage) Get(ctx context.Context, key string) (interface{}, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		m.logger.Warn("Get operation cancelled", "key", key, "error", ctx.Err())
		return nil, ctx.Err()
	default:
	}

	value, ok := m.data.Load(key)
	if !ok {
		return nil, nil
	}

	// Check if the value has expired
	if item, ok := value.(*memoryItem); ok {
		if item.expiresAt > 0 && time.Now().Unix() >= item.expiresAt {
			m.data.Delete(key)
			m.logger.Debug("Item expired and removed", "key", key)
			return nil, nil
		}
		return item.value, nil
	}

	return value, nil
}

// Set stores a value in memory storage with optional expiration
func (m *MemoryStorage) Set(ctx context.Context, key string, value interface{}, expiration int64) error {
	// Check context cancellation
	select {
	case <-ctx.Done():
		m.logger.Warn("Set operation cancelled", "key", key, "error", ctx.Err())
		return ctx.Err()
	default:
	}

	item := &memoryItem{
		value:     value,
		expiresAt: expiration,
	}
	m.data.Store(key, item)
	m.logger.Debug("Item stored", "key", key, "expiration", expiration)
	return nil
}

// Delete removes a value from memory storage
func (m *MemoryStorage) Delete(ctx context.Context, key string) error {
	// Check context cancellation
	select {
	case <-ctx.Done():
		m.logger.Warn("Delete operation cancelled", "key", key, "error", ctx.Err())
		return ctx.Err()
	default:
	}

	m.data.Delete(key)
	m.logger.Debug("Item deleted", "key", key)
	return nil
}

// Close closes the memory storage (no-op for in-memory storage)
func (m *MemoryStorage) Close() error {
	m.logger.Info("Memory storage closed")
	return nil
}

// memoryItem represents a stored item with expiration
type memoryItem struct {
	value     interface{}
	expiresAt int64 // Unix timestamp, 0 means no expiration
}

