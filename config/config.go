// Package config provides centralized configuration management for the HTMX learning application.
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all application configuration
type Config struct {
	// Server configuration
	Port         string        `env:"PORT"`
	Host         string        `env:"HOST"`
	ReadTimeout  time.Duration `env:"READ_TIMEOUT"`
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT"`
	IdleTimeout  time.Duration `env:"IDLE_TIMEOUT"`
	
	// Database configuration
	DatabaseURL     string `env:"DATABASE_URL"`
	MaxConnections  int32  `env:"DB_MAX_CONNECTIONS"`
	MinConnections  int32  `env:"DB_MIN_CONNECTIONS"`
	ConnMaxLifetime time.Duration `env:"DB_CONN_MAX_LIFETIME"`
	
	// Security configuration
	AllowedOrigins []string `env:"ALLOWED_ORIGINS"`
	TrustedProxies []string `env:"TRUSTED_PROXIES"`
	SecretKey      string   `env:"SECRET_KEY"`
	
	// Logging configuration
	LogLevel  string `env:"LOG_LEVEL"`
	LogFormat string `env:"LOG_FORMAT"`
	
	// Rate limiting configuration
	RateLimit         int           `env:"RATE_LIMIT"`
	RateLimitWindow   time.Duration `env:"RATE_LIMIT_WINDOW"`
	RateLimitBurst    int           `env:"RATE_LIMIT_BURST"`
	
	// Application configuration
	Environment string `env:"ENVIRONMENT"`
	Debug       bool   `env:"DEBUG"`
}

// Load loads configuration from environment variables with sensible defaults
func Load() (*Config, error) {
	config := &Config{
		// Server defaults
		Port:         getEnv("PORT", "8080"),
		Host:         getEnv("HOST", "localhost"),
		ReadTimeout:  parseDuration("READ_timeout", getEnv("READ_TIMEOUT", "15s")),
		WriteTimeout: parseDuration("write_timeout", getEnv("WRITE_TIMEOUT", "15s")),
		IdleTimeout:  parseDuration("idle_timeout", getEnv("IDLE_TIMEOUT", "60s")),
		
		// Database defaults
		DatabaseURL:     getRequiredEnv("DATABASE_URL"),
		MaxConnections:  int32(parseInt("DB_MAX_CONNECTIONS", getEnv("DB_MAX_CONNECTIONS", "10"))),
		MinConnections:  int32(parseInt("DB_MIN_CONNECTIONS", getEnv("DB_MIN_CONNECTIONS", "2"))),
		ConnMaxLifetime: parseDuration("db_conn_max_lifetime", getEnv("DB_CONN_MAX_LIFETIME", "1h")),
		
		// Security defaults
		AllowedOrigins: parseStringSlice(getEnv("ALLOWED_ORIGINS", "http://localhost:8080,https://localhost:8080")),
		TrustedProxies: parseStringSlice(getEnv("TRUSTED_PROXIES", "127.0.0.1,::1")),
		SecretKey:      getRequiredEnv("SECRET_KEY"),
		
		// Logging defaults
		LogLevel:  getEnv("LOG_LEVEL", "info"),
		LogFormat: getEnv("LOG_FORMAT", "json"),
		
		// Rate limiting defaults
		RateLimit:       parseInt("RATE_LIMIT", getEnv("RATE_LIMIT", "100")),
		RateLimitWindow: parseDuration("rate_limit_window", getEnv("RATE_LIMIT_WINDOW", "1m")),
		RateLimitBurst:  parseInt("RATE_LIMIT_BURST", getEnv("RATE_LIMIT_BURST", "20")),
		
		// Application defaults
		Environment: getEnv("ENVIRONMENT", "development"),
		Debug:       parseBool("DEBUG", getEnv("DEBUG", "false")),
	}
	
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}
	
	return config, nil
}

// Validate ensures the configuration is valid
func (c *Config) Validate() error {
	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	
	if c.SecretKey == "" {
		return fmt.Errorf("SECRET_KEY is required")
	}
	
	if len(c.SecretKey) < 32 {
		return fmt.Errorf("SECRET_KEY must be at least 32 characters long")
	}
	
	if c.MaxConnections < c.MinConnections {
		return fmt.Errorf("DB_MAX_CONNECTIONS must be greater than DB_MIN_CONNECTIONS")
	}
	
	if len(c.AllowedOrigins) == 0 {
		return fmt.Errorf("ALLOWED_ORIGINS must be specified")
	}
	
	validEnvs := map[string]bool{"development": true, "staging": true, "production": true}
	if !validEnvs[c.Environment] {
		return fmt.Errorf("ENVIRONMENT must be one of: development, staging, production")
	}
	
	return nil
}


// GetServerAddress returns the full server address
func (c *Config) GetServerAddress() string {
	if strings.HasPrefix(c.Port, ":") {
		return c.Port
	}
	return ":" + c.Port
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getRequiredEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("required environment variable %s is not set", key))
	}
	return value
}

func parseInt(key, value string) int {
	if i, err := strconv.Atoi(value); err == nil {
		return i
	}
	panic(fmt.Sprintf("invalid integer value for %s: %s", key, value))
}

func parseBool(key, value string) bool {
	if b, err := strconv.ParseBool(value); err == nil {
		return b
	}
	panic(fmt.Sprintf("invalid boolean value for %s: %s", key, value))
}

func parseDuration(key, value string) time.Duration {
	if d, err := time.ParseDuration(value); err == nil {
		return d
	}
	panic(fmt.Sprintf("invalid duration value for %s: %s", key, value))
}

func parseStringSlice(value string) []string {
	if value == "" {
		return []string{}
	}
	
	parts := strings.Split(value, ",")
	result := make([]string, len(parts))
	for i, part := range parts {
		result[i] = strings.TrimSpace(part)
	}
	return result
}