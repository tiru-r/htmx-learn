package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"htmx-learn/config"
	"htmx-learn/db"
	"htmx-learn/handlers"
	"htmx-learn/middleware"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}
	
	// Initialize structured logging
	var logger *slog.Logger
	if cfg.LogFormat == "json" {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: parseLogLevel(cfg.LogLevel),
		}))
	} else {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: parseLogLevel(cfg.LogLevel),
		}))
	}
	slog.SetDefault(logger)
	
	slog.Info("Starting HTMX learning application",
		"version", "1.0.0",
		"environment", cfg.Environment,
		"port", cfg.Port)

	// Initialize database with pool configuration
	database, err := db.New(cfg.DatabaseURL, cfg.MaxConnections, cfg.MinConnections)
	if err != nil {
		slog.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer database.Close()

	// Initialize database schema
	ctx := context.Background()
	if err := database.InitSchema(ctx); err != nil {
		slog.Error("Failed to initialize database schema", "error", err)
		os.Exit(1)
	}

	// Initialize handlers with database and configuration
	h := handlers.New(database, cfg)

	mux := http.NewServeMux()

	// Static file serving
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

	// Page routes
	mux.HandleFunc("GET /", h.Home)
	mux.HandleFunc("GET /counter", h.CounterPage)
	mux.HandleFunc("GET /dynamic", h.DynamicPage)

	// API routes for counter
	mux.HandleFunc("POST /counter/increment", h.CounterIncrement)
	mux.HandleFunc("POST /counter/decrement", h.CounterDecrement)
	mux.HandleFunc("POST /counter/reset", h.CounterReset)

	// API routes for dynamic content
	mux.HandleFunc("GET /api/time", h.GetTime)
	mux.HandleFunc("GET /api/users", h.GetUsers)
	mux.HandleFunc("GET /api/users/paginated", h.GetUsersPaginated)
	mux.HandleFunc("POST /api/users", h.CreateUser)
	mux.HandleFunc("DELETE /api/users/{id}", h.DeleteUser)
	mux.HandleFunc("POST /api/search", h.SearchUsers)
	mux.HandleFunc("POST /api/search/paginated", h.SearchUsersPaginated)
	
	// Health check routes
	mux.HandleFunc("GET /health", h.HealthCheck)
	mux.HandleFunc("GET /health/ready", h.ReadinessCheck)
	mux.HandleFunc("GET /health/live", h.LivenessCheck)

	// Apply middleware with configuration
	handler := middleware.Recovery(
		middleware.Logger(
			middleware.SecurityHeaders(
				middleware.ConfigurableCORS(cfg.AllowedOrigins,
					middleware.RateLimit(cfg,
						mux),
				),
			),
		),
	)

	server := &http.Server{
		Addr:         cfg.GetServerAddress(),
		Handler:      handler,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		slog.Info("Server starting", "address", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down server...")

	// Create a deadline to wait for
	shutdownTimeout := 30 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("Server exited gracefully")
}

// parseLogLevel converts string log level to slog.Level
func parseLogLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}