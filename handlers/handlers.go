package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"htmx-learn/config"
	"htmx-learn/db"
	"htmx-learn/templates/components"
	"htmx-learn/templates/pages"
	"htmx-learn/validation"
	"github.com/jackc/pgx/v5"
)

type Handlers struct {
	counterStore db.CounterRepository
	userStore    db.UserRepository
	config       *config.Config
	database     *db.DB
}

func New(database *db.DB, cfg *config.Config) *Handlers {
	return &Handlers{
		counterStore: db.NewCounterStore(database),
		userStore:    db.NewUserStore(database),
		config:       cfg,
		database:     database,
	}
}

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, r, pages.Home())
}

func (h *Handlers) CounterPage(w http.ResponseWriter, r *http.Request) {
	count, err := h.counterStore.Get(r.Context())
	if err != nil {
		slog.Error("Error getting counter", "error", err)
		count = 0
	}
	renderTemplate(w, r, pages.CounterPage(count))
}

func (h *Handlers) DynamicPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, r, pages.DynamicPage())
}

func (h *Handlers) CounterIncrement(w http.ResponseWriter, r *http.Request) {
	count, err := h.counterStore.Increment(r.Context())
	if err != nil {
		handleError(w, "incrementing counter", err)
		return
	}
	renderTemplate(w, r, components.CountDisplay(count))
}

func (h *Handlers) CounterDecrement(w http.ResponseWriter, r *http.Request) {
	count, err := h.counterStore.Decrement(r.Context())
	if err != nil {
		handleError(w, "decrementing counter", err)
		return
	}
	renderTemplate(w, r, components.CountDisplay(count))
}

func (h *Handlers) CounterReset(w http.ResponseWriter, r *http.Request) {
	count, err := h.counterStore.Reset(r.Context())
	if err != nil {
		handleError(w, "resetting counter", err)
		return
	}
	renderTemplate(w, r, components.CountDisplay(count))
}

func (h *Handlers) GetTime(w http.ResponseWriter, r *http.Request) {
	currentTime := time.Now()
	renderTemplate(w, r, components.TimeDisplay(currentTime))
}

func (h *Handlers) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userStore.GetAll(r.Context())
	if err != nil {
		handleError(w, "getting users", err)
		return
	}
	
	templateUsers := convertToTemplateUsers(users)
	
	for _, user := range templateUsers {
		if err := components.UserCard(user).Render(r.Context(), w); err != nil {
			slog.Error("Template rendering error", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}

func (h *Handlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}
	
	// Sanitize and validate input
	input := validation.UserInput{
		Name:  validation.SanitizeInput(r.FormValue("user-name")),
		Email: validation.SanitizeInput(r.FormValue("user-email")),
	}
	
	if err := validation.ValidateUser(input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	user, err := h.userStore.Add(r.Context(), input.Name, input.Email)
	if err != nil {
		handleError(w, "creating user", err)
		return
	}
	
	templateUser := convertToTemplateUser(user)
	renderTemplate(w, r, components.UserCard(templateUser))
}

func (h *Handlers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	
	err = h.userStore.Delete(r.Context(), id)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		handleError(w, "deleting user", err)
		return
	}
	
	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) SearchUsers(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}
	
	// Sanitize search query
	query := validation.SanitizeInput(r.FormValue("search"))
	users, err := h.userStore.Search(r.Context(), query)
	if err != nil {
		handleError(w, "searching users", err)
		return
	}
	
	templateUsers := convertToTemplateUsers(users)
	renderTemplate(w, r, components.SearchResults(templateUsers))
}

// GetUsersPaginated handles paginated user listing - reverted to original approach
func (h *Handlers) GetUsersPaginated(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters
	params, err := parsePaginationParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get paginated users
	result, err := h.userStore.GetAllPaginated(r.Context(), params)
	if err != nil {
		handleError(w, "getting paginated users", err)
		return
	}

	templateUsers := convertToTemplateUsers(result.Data)

	// For HTMX requests, return just the user cards and pagination
	if r.Header.Get("HX-Request") == "true" {
		// Render user cards
		for _, user := range templateUsers {
			if err := components.UserCard(user).Render(r.Context(), w); err != nil {
				slog.Error("Template rendering error", "error", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		}
		
		// Render pagination component
		paginationData := components.PaginationData{
			CurrentPage: result.Page,
			TotalPages:  result.TotalPages,
			HasPrev:     result.HasPrev,
			HasNext:     result.HasNext,
			BaseURL:     "/api/users/paginated",
			SearchQuery: "",
		}
		renderTemplate(w, r, components.Pagination(paginationData))
		return
	}

	// For non-HTMX requests, render the full page
	renderTemplate(w, r, pages.DynamicPage())
}

// SearchUsersPaginated handles paginated user search
func (h *Handlers) SearchUsersPaginated(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Parse pagination parameters
	params, err := parsePaginationParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Sanitize search query
	query := validation.SanitizeInput(r.FormValue("search"))
	
	result, err := h.userStore.SearchPaginated(r.Context(), query, params)
	if err != nil {
		handleError(w, "searching users with pagination", err)
		return
	}

	templateUsers := convertToTemplateUsers(result.Data)
	renderTemplate(w, r, components.SearchResults(templateUsers))
	
	// Also render pagination component for search results
	paginationData := components.PaginationData{
		CurrentPage: result.Page,
		TotalPages:  result.TotalPages,
		HasPrev:     result.HasPrev,
		HasNext:     result.HasNext,
		BaseURL:     "/api/search/paginated",
		SearchQuery: query,
	}
	renderTemplate(w, r, components.Pagination(paginationData))
}

// HealthStatus represents the health status of the application
type HealthStatus struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Version   string            `json:"version"`
	Checks    map[string]Health `json:"checks"`
}

// Health represents individual health check status
type Health struct {
	Status  string        `json:"status"`
	Message string        `json:"message,omitempty"`
	Latency time.Duration `json:"latency"`
}

// HealthCheck provides a health check endpoint
func (h *Handlers) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	checks := make(map[string]Health)
	overallStatus := "healthy"
	
	// Database health check
	dbStart := time.Now()
	if err := h.checkDatabaseHealth(r.Context()); err != nil {
		checks["database"] = Health{
			Status:  "unhealthy",
			Message: err.Error(),
			Latency: time.Since(dbStart),
		}
		overallStatus = "unhealthy"
	} else {
		checks["database"] = Health{
			Status:  "healthy",
			Latency: time.Since(dbStart),
		}
	}
	
	status := HealthStatus{
		Status:    overallStatus,
		Timestamp: time.Now(),
		Version:   "1.0.0",
		Checks:    checks,
	}
	
	statusCode := http.StatusOK
	if overallStatus != "healthy" {
		statusCode = http.StatusServiceUnavailable
	}
	
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(status)
}

// ReadinessCheck provides a readiness check endpoint
func (h *Handlers) ReadinessCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Check if all dependencies are ready
	if err := h.checkDatabaseHealth(r.Context()); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":    "not ready",
			"timestamp": time.Now(),
			"error":     err.Error(),
		})
		return
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "ready",
		"timestamp": time.Now(),
	})
}

// LivenessCheck provides a liveness check endpoint
func (h *Handlers) LivenessCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "alive",
		"timestamp": time.Now(),
	})
}

// checkDatabaseHealth performs a simple database health check
func (h *Handlers) checkDatabaseHealth(ctx context.Context) error {
	// Create a timeout context for the health check
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	
	// Simple ping to check database connectivity
	conn, err := h.database.Pool.Acquire(timeoutCtx)
	if err != nil {
		return err
	}
	defer conn.Release()
	
	// Execute a simple query
	var result int
	err = conn.QueryRow(timeoutCtx, "SELECT 1").Scan(&result)
	if err != nil {
		return err
	}
	
	return nil
}