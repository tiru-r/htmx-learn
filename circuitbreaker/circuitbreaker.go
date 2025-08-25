// Package circuitbreaker provides a simple circuit breaker implementation
// for protecting external dependencies like databases from cascading failures.
package circuitbreaker

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"
)

var (
	ErrCircuitBreakerOpen     = errors.New("circuit breaker is open")
	ErrCircuitBreakerTimeout  = errors.New("circuit breaker timeout")
)

// State represents the current state of the circuit breaker
type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

// Config holds circuit breaker configuration
type Config struct {
	MaxFailures     int           // Maximum failures before opening
	ResetTimeout    time.Duration // Time to wait before transitioning to half-open
	FailureTimeout  time.Duration // Timeout for individual calls
	MaxRequests     int           // Maximum requests allowed in half-open state
}

// DefaultConfig returns a default circuit breaker configuration
func DefaultConfig() Config {
	return Config{
		MaxFailures:     5,
		ResetTimeout:    30 * time.Second,
		FailureTimeout:  10 * time.Second,
		MaxRequests:     3,
	}
}

// CircuitBreaker implements a circuit breaker pattern
type CircuitBreaker struct {
	config       Config
	state        State
	failures     int
	requests     int
	lastFailTime time.Time
	mu           sync.RWMutex
}

// New creates a new circuit breaker with the given configuration
func New(config Config) *CircuitBreaker {
	return &CircuitBreaker{
		config: config,
		state:  StateClosed,
	}
}

// Execute runs the given function with circuit breaker protection
func (cb *CircuitBreaker) Execute(ctx context.Context, fn func(context.Context) error) error {
	if !cb.allowRequest() {
		return ErrCircuitBreakerOpen
	}

	// Create timeout context
	timeoutCtx, cancel := context.WithTimeout(ctx, cb.config.FailureTimeout)
	defer cancel()

	// Execute with timeout
	done := make(chan error, 1)
	go func() {
		done <- fn(timeoutCtx)
	}()

	select {
	case err := <-done:
		if err != nil {
			cb.recordFailure()
			return err
		}
		cb.recordSuccess()
		return nil
	case <-timeoutCtx.Done():
		cb.recordFailure()
		return ErrCircuitBreakerTimeout
	}
}

// allowRequest checks if a request should be allowed through the circuit breaker
func (cb *CircuitBreaker) allowRequest() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	now := time.Now()

	switch cb.state {
	case StateClosed:
		return true
	case StateOpen:
		if now.Sub(cb.lastFailTime) > cb.config.ResetTimeout {
			cb.state = StateHalfOpen
			cb.requests = 0
			slog.Info("Circuit breaker transitioning to half-open state")
			return true
		}
		return false
	case StateHalfOpen:
		return cb.requests < cb.config.MaxRequests
	default:
		return false
	}
}

// recordSuccess records a successful request
func (cb *CircuitBreaker) recordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateHalfOpen:
		cb.requests++
		if cb.requests >= cb.config.MaxRequests {
			cb.state = StateClosed
			cb.failures = 0
			cb.requests = 0
			slog.Info("Circuit breaker transitioning to closed state")
		}
	case StateClosed:
		cb.failures = 0
	}
}

// recordFailure records a failed request
func (cb *CircuitBreaker) recordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failures++
	cb.lastFailTime = time.Now()

	switch cb.state {
	case StateClosed:
		if cb.failures >= cb.config.MaxFailures {
			cb.state = StateOpen
			slog.Warn("Circuit breaker opening due to failures",
				"failures", cb.failures,
				"max_failures", cb.config.MaxFailures)
		}
	case StateHalfOpen:
		cb.state = StateOpen
		slog.Warn("Circuit breaker opening from half-open state due to failure")
	}
}

// GetState returns the current state of the circuit breaker
func (cb *CircuitBreaker) GetState() State {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// GetStats returns statistics about the circuit breaker
func (cb *CircuitBreaker) GetStats() map[string]interface{} {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	return map[string]interface{}{
		"state":         cb.state.String(),
		"failures":      cb.failures,
		"requests":      cb.requests,
		"last_failure":  cb.lastFailTime,
	}
}