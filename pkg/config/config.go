package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds application configuration
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Storage  StorageConfig  `mapstructure:"storage"`
	Limiter  LimiterConfig  `mapstructure:"limiter"`
	CORS     CORSConfig     `mapstructure:"cors"`
	Webhook  WebhookConfig  `mapstructure:"webhook"`
	Tenant   TenantConfig   `mapstructure:"tenant"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

// StorageConfig holds storage configuration
type StorageConfig struct {
	Type          string `mapstructure:"type"` // "memory" or "redis"
	RedisAddress  string `mapstructure:"redis_address"`
	RedisDB       int    `mapstructure:"redis_db"`
	RedisPassword string `mapstructure:"redis_password"`
	RedisPoolSize int    `mapstructure:"redis_pool_size"`
	RedisMinIdle  int    `mapstructure:"redis_min_idle"`
}

// LimiterConfig holds rate limiter configuration
type LimiterConfig struct {
	DefaultAlgorithm string        `mapstructure:"default_algorithm"` // "token_bucket" or "sliding_window"
	DefaultLimit     int           `mapstructure:"default_limit"`
	DefaultWindow    time.Duration `mapstructure:"default_window"`
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
}

// WebhookConfig holds webhook configuration
type WebhookConfig struct {
	URL     string        `mapstructure:"url"`
	Timeout time.Duration `mapstructure:"timeout"`
	Enabled bool          `mapstructure:"enabled"`
}

// TenantConfig holds tenant management configuration
type TenantConfig struct {
	Enabled bool `mapstructure:"enabled"`
}

// Load loads configuration from file and environment variables
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")

	// Set defaults
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.read_timeout", "15s")
	viper.SetDefault("server.write_timeout", "15s")
	viper.SetDefault("server.idle_timeout", "60s")
	viper.SetDefault("storage.type", "memory")
	viper.SetDefault("storage.redis_address", "localhost:6379")
	viper.SetDefault("storage.redis_db", 0)
	viper.SetDefault("storage.redis_pool_size", 10)
	viper.SetDefault("storage.redis_min_idle", 5)
	viper.SetDefault("limiter.default_algorithm", "token_bucket")
	viper.SetDefault("limiter.default_limit", 100)
	viper.SetDefault("limiter.default_window", "1m")
	viper.SetDefault("cors.allowed_origins", []string{"*"})
	viper.SetDefault("webhook.enabled", false)
	viper.SetDefault("webhook.timeout", "5s")
	viper.SetDefault("tenant.enabled", false)

	// Read from environment variables
	viper.SetEnvPrefix("RL")
	viper.AutomaticEnv()

	// Bind environment variables
	bindEnvVars()

	// Read config file if exists
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &config, nil
}

func bindEnvVars() {
	// Server
	viper.BindEnv("server.port", "RL_SERVER_PORT")
	viper.BindEnv("server.read_timeout", "RL_SERVER_READ_TIMEOUT")
	viper.BindEnv("server.write_timeout", "RL_SERVER_WRITE_TIMEOUT")
	viper.BindEnv("server.idle_timeout", "RL_SERVER_IDLE_TIMEOUT")

	// Storage
	viper.BindEnv("storage.type", "RL_STORAGE_TYPE")
	viper.BindEnv("storage.redis_address", "RL_REDIS_ADDRESS")
	viper.BindEnv("storage.redis_db", "RL_REDIS_DB")
	viper.BindEnv("storage.redis_password", "RL_REDIS_PASSWORD")
	viper.BindEnv("storage.redis_pool_size", "RL_REDIS_POOL_SIZE")
	viper.BindEnv("storage.redis_min_idle", "RL_REDIS_MIN_IDLE")

	// Limiter
	viper.BindEnv("limiter.default_algorithm", "RL_DEFAULT_ALGORITHM")
	viper.BindEnv("limiter.default_limit", "RL_DEFAULT_LIMIT")
	viper.BindEnv("limiter.default_window", "RL_DEFAULT_WINDOW")

	// CORS
	viper.BindEnv("cors.allowed_origins", "RL_CORS_ALLOWED_ORIGINS")

	// Webhook
	viper.BindEnv("webhook.url", "RL_WEBHOOK_URL")
	viper.BindEnv("webhook.enabled", "RL_WEBHOOK_ENABLED")
	viper.BindEnv("webhook.timeout", "RL_WEBHOOK_TIMEOUT")

	// Tenant
	viper.BindEnv("tenant.enabled", "RL_TENANT_ENABLED")

	// Override with direct env vars if set
	if port := os.Getenv("RL_SERVER_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			viper.Set("server.port", p)
		}
	}

	if limit := os.Getenv("RL_DEFAULT_LIMIT"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			viper.Set("limiter.default_limit", l)
		}
	}

	if window := os.Getenv("RL_DEFAULT_WINDOW"); window != "" {
		if d, err := time.ParseDuration(window); err == nil {
			viper.Set("limiter.default_window", d)
		}
	}

	if origins := os.Getenv("RL_CORS_ALLOWED_ORIGINS"); origins != "" {
		originsList := strings.Split(origins, ",")
		for i := range originsList {
			originsList[i] = strings.TrimSpace(originsList[i])
		}
		viper.Set("cors.allowed_origins", originsList)
	}
}

