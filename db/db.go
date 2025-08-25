package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"htmx-learn/circuitbreaker"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DB holds the database connection pool and circuit breaker
type DB struct {
	*pgxpool.Pool
	CircuitBreaker *circuitbreaker.CircuitBreaker
}

// New creates a new database connection pool with configurable pool settings
func New(databaseURL string, maxConns, minConns int32) (*DB, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	// Set connection pool settings
	config.MaxConns = maxConns
	config.MinConns = minConns

	// Use context with timeout for initialization
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test the connection
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Initialize circuit breaker
	cb := circuitbreaker.New(circuitbreaker.DefaultConfig())

	return &DB{
		Pool:           pool,
		CircuitBreaker: cb,
	}, nil
}

const (
	// Maximum schema file size (1MB) to prevent memory exhaustion
	maxSchemaFileSize = 1024 * 1024
)

// InitSchema initializes the database schema with size limits for security
func (db *DB) InitSchema(ctx context.Context) error {
	// Check file size before reading to prevent memory exhaustion attacks
	fileInfo, err := os.Stat("db/schema.sql")
	if err != nil {
		return fmt.Errorf("failed to stat schema file: %w", err)
	}

	if fileInfo.Size() > maxSchemaFileSize {
		return fmt.Errorf("schema file too large: %d bytes (max %d bytes)", 
			fileInfo.Size(), maxSchemaFileSize)
	}

	schemaSQL, err := os.ReadFile("db/schema.sql")
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	if _, err := db.Exec(ctx, string(schemaSQL)); err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	return nil
}

// ExecuteWithCircuitBreaker executes a database operation with circuit breaker protection
func (db *DB) ExecuteWithCircuitBreaker(ctx context.Context, operation func(context.Context) error) error {
	return db.CircuitBreaker.Execute(ctx, operation)
}

// Close closes the database connection pool
func (db *DB) Close() {
	db.Pool.Close()
}