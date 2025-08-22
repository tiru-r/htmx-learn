package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"htmx-learn/db"
	"htmx-learn/handlers"
	"htmx-learn/middleware"
)

func main() {
	// Get database path from environment or use default
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "./data/app.db"
	}

	// Initialize database
	database, err := db.New(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Initialize database schema
	if err := database.InitSchema(); err != nil {
		log.Fatalf("Failed to initialize database schema: %v", err)
	}

	// Initialize handlers with database
	h := handlers.New(database)

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
	mux.HandleFunc("POST /api/users", h.CreateUser)
	mux.HandleFunc("DELETE /api/users/{id}", h.DeleteUser)
	mux.HandleFunc("POST /api/search", h.SearchUsers)

	// Apply middleware
	handler := middleware.Recovery(
		middleware.Logger(
			middleware.CORS(mux),
		),
	)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on http://localhost%s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}