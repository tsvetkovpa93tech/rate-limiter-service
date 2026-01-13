package storage

import (
	"sync"
	"time"
)

// MemoryStorage implements in-memory storage using sync.Map
type MemoryStorage struct {
	data sync.Map
}

// NewMemoryStorage creates a new in-memory storage instance
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data: sync.Map{},
	}
}

// Get retrieves a value from memory storage
func (m *MemoryStorage) Get(key string) (interface{}, error) {
	value, ok := m.data.Load(key)
	if !ok {
		return nil, nil
	}

	// Check if the value has expired
	if item, ok := value.(*memoryItem); ok {
		if item.expiresAt > 0 && time.Now().Unix() > item.expiresAt {
			m.data.Delete(key)
			return nil, nil
		}
		return item.value, nil
	}

	return value, nil
}

// Set stores a value in memory storage with optional expiration
func (m *MemoryStorage) Set(key string, value interface{}, expiration int64) error {
	item := &memoryItem{
		value:     value,
		expiresAt: expiration,
	}
	m.data.Store(key, item)
	return nil
}

// Delete removes a value from memory storage
func (m *MemoryStorage) Delete(key string) error {
	m.data.Delete(key)
	return nil
}

// Close closes the memory storage (no-op for in-memory storage)
func (m *MemoryStorage) Close() error {
	return nil
}

// memoryItem represents a stored item with expiration
type memoryItem struct {
	value     interface{}
	expiresAt int64 // Unix timestamp, 0 means no expiration
}

