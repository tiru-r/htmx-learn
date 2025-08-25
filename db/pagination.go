// Package db provides database connection management and data access operations
// for the HTMX learning application using PostgreSQL with pgx driver.
package db

const (
	// Default pagination settings
	DefaultPageSize = 10
	MaxPageSize     = 100
	MinPageSize     = 5
)

// PaginationParams holds pagination parameters
type PaginationParams struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Offset   int `json:"offset"`
}

// PaginatedResult holds paginated query results with metadata
type PaginatedResult[T any] struct {
	Data       []T `json:"data"`
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}

// NewPaginationParams creates validated pagination parameters
func NewPaginationParams(page, pageSize int) PaginationParams {
	// Validate and set defaults
	if page < 1 {
		page = 1
	}
	
	if pageSize < MinPageSize {
		pageSize = DefaultPageSize
	}
	
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}
	
	offset := (page - 1) * pageSize
	
	return PaginationParams{
		Page:     page,
		PageSize: pageSize,
		Offset:   offset,
	}
}

// NewPaginatedResult creates a paginated result with metadata
func NewPaginatedResult[T any](data []T, params PaginationParams, total int) *PaginatedResult[T] {
	totalPages := (total + params.PageSize - 1) / params.PageSize // Ceiling division
	if totalPages < 1 {
		totalPages = 1
	}
	
	return &PaginatedResult[T]{
		Data:       data,
		Page:       params.Page,
		PageSize:   params.PageSize,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    params.Page < totalPages,
		HasPrev:    params.Page > 1,
	}
}