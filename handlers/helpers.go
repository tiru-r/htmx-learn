// Package handlers provides HTTP request handlers for the HTMX learning application.
package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"htmx-learn/db"
	"htmx-learn/templates/components"
	"github.com/a-h/templ"
)

// renderTemplate renders a templ component and handles errors consistently
func renderTemplate(w http.ResponseWriter, r *http.Request, component templ.Component) {
	if err := component.Render(r.Context(), w); err != nil {
		slog.Error("Template rendering error", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// handleError logs an error with context and sends an appropriate HTTP error response
func handleError(w http.ResponseWriter, context string, err error) {
	slog.Error("Handler error", "context", context, "error", err)
	http.Error(w, "Internal server error", http.StatusInternalServerError)
}

// convertToTemplateUsers converts database users to template users
func convertToTemplateUsers(users []*db.User) []components.User {
	if users == nil {
		return nil
	}

	result := make([]components.User, len(users))
	for i, user := range users {
		result[i] = components.User{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		}
	}
	return result
}

// convertToTemplateUser converts a single database user to template user
func convertToTemplateUser(user *db.User) components.User {
	return components.User{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
}

// parsePaginationParams extracts and validates pagination parameters from request
func parsePaginationParams(r *http.Request) (db.PaginationParams, error) {
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")
	
	page := 1
	pageSize := db.DefaultPageSize
	
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil {
			page = p
		}
	}
	
	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil {
			pageSize = ps
		}
	}
	
	// Create validated pagination params
	params := db.NewPaginationParams(page, pageSize)
	return params, nil
}