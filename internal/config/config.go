package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server  ServerConfig  `mapstructure:"server"`
	Storage StorageConfig `mapstructure:"storage"`
	Limiter LimiterConfig `mapstructure:"limiter"`
}

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

type StorageConfig struct {
	Type         string `mapstructure:"type"` // "memory" or "redis"
	RedisAddress string `mapstructure:"redis_address"`
	RedisDB      int    `mapstructure:"redis_db"`
	RedisPassword string `mapstructure:"redis_password"`
}

type LimiterConfig struct {
	DefaultAlgorithm string        `mapstructure:"default_algorithm"` // "token_bucket" or "sliding_window"
	DefaultLimit     int           `mapstructure:"default_limit"`
	DefaultWindow    time.Duration `mapstructure:"default_window"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")

	// Set defaults
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("storage.type", "memory")
	viper.SetDefault("storage.redis_address", "localhost:6379")
	viper.SetDefault("storage.redis_db", 0)
	viper.SetDefault("limiter.default_algorithm", "token_bucket")
	viper.SetDefault("limiter.default_limit", 100)
	viper.SetDefault("limiter.default_window", "1m")

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

	// Storage
	viper.BindEnv("storage.type", "RL_STORAGE_TYPE")
	viper.BindEnv("storage.redis_address", "RL_REDIS_ADDRESS")
	viper.BindEnv("storage.redis_db", "RL_REDIS_DB")
	viper.BindEnv("storage.redis_password", "RL_REDIS_PASSWORD")

	// Limiter
	viper.BindEnv("limiter.default_algorithm", "RL_DEFAULT_ALGORITHM")
	viper.BindEnv("limiter.default_limit", "RL_DEFAULT_LIMIT")
	viper.BindEnv("limiter.default_window", "RL_DEFAULT_WINDOW")

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
}

