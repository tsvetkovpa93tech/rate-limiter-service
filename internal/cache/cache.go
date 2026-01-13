package cache

import (
	"sync"
	"time"
)

// Cache provides in-memory caching for rate limit results
type Cache struct {
	data            sync.Map
	ttl             time.Duration
	cleanupInterval time.Duration
	stop            chan struct{}
}

// cacheItem represents a cached item with expiration
type cacheItem struct {
	value     interface{}
	expiresAt time.Time
}

// NewCache creates a new cache instance
func NewCache(ttl time.Duration, cleanupInterval time.Duration) *Cache {
	c := &Cache{
		ttl:             ttl,
		cleanupInterval: cleanupInterval,
		stop:            make(chan struct{}),
	}

	// Start cleanup goroutine
	go c.cleanup()

	return c
}

// Get retrieves a value from cache
func (c *Cache) Get(key string) (interface{}, bool) {
	item, ok := c.data.Load(key)
	if !ok {
		return nil, false
	}

	cachedItem := item.(*cacheItem)
	if time.Now().After(cachedItem.expiresAt) {
		c.data.Delete(key)
		return nil, false
	}

	return cachedItem.value, true
}

// Set stores a value in cache
func (c *Cache) Set(key string, value interface{}) {
	c.data.Store(key, &cacheItem{
		value:     value,
		expiresAt: time.Now().Add(c.ttl),
	})
}

// Delete removes a value from cache
func (c *Cache) Delete(key string) {
	c.data.Delete(key)
}

// cleanup periodically removes expired items
func (c *Cache) cleanup() {
	ticker := time.NewTicker(c.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			now := time.Now()
			c.data.Range(func(key, value interface{}) bool {
				item := value.(*cacheItem)
				if now.After(item.expiresAt) {
					c.data.Delete(key)
				}
				return true
			})
		case <-c.stop:
			return
		}
	}
}

// Close stops the cache cleanup goroutine
func (c *Cache) Close() {
	close(c.stop)
}

// CacheKey generates a cache key from parameters
func CacheKey(key string, algorithm string, limit int, window time.Duration) string {
	return key + ":" + algorithm + ":" + string(rune(limit)) + ":" + window.String()
}
