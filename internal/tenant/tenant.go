package tenant

import (
	"sync"
	"time"
)

// TenantConfig holds configuration for a tenant
type TenantConfig struct {
	APIKey    string
	Algorithm string
	Limit     int
	Window    time.Duration
	Enabled   bool
}

// Manager manages tenant configurations
type Manager struct {
	configs sync.Map // map[string]*TenantConfig
	mu      sync.RWMutex
}

// NewManager creates a new tenant manager
func NewManager() *Manager {
	return &Manager{}
}

// GetConfig retrieves tenant configuration
func (m *Manager) GetConfig(apiKey string) (*TenantConfig, bool) {
	config, ok := m.configs.Load(apiKey)
	if !ok {
		return nil, false
	}
	return config.(*TenantConfig), true
}

// SetConfig sets tenant configuration
func (m *Manager) SetConfig(config *TenantConfig) {
	m.configs.Store(config.APIKey, config)
}

// UpdateConfig updates tenant configuration dynamically
func (m *Manager) UpdateConfig(apiKey string, limit int, window time.Duration) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	config, ok := m.configs.Load(apiKey)
	if !ok {
		return false
	}

	tenantConfig := config.(*TenantConfig)
	tenantConfig.Limit = limit
	tenantConfig.Window = window
	m.configs.Store(apiKey, tenantConfig)
	return true
}

// DeleteConfig removes tenant configuration
func (m *Manager) DeleteConfig(apiKey string) {
	m.configs.Delete(apiKey)
}

// ListConfigs returns all tenant configurations
func (m *Manager) ListConfigs() []*TenantConfig {
	var configs []*TenantConfig
	m.configs.Range(func(key, value interface{}) bool {
		configs = append(configs, value.(*TenantConfig))
		return true
	})
	return configs
}

