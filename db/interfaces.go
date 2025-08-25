// Package db provides database connection management and data access operations
// for the HTMX learning application using PostgreSQL with pgx driver.
package db

import "context"

// UserRepository defines the interface for user data operations
type UserRepository interface {
	GetAll(ctx context.Context) ([]*User, error)
	GetAllPaginated(ctx context.Context, params PaginationParams) (*PaginatedResult[*User], error)
	Add(ctx context.Context, name, email string) (*User, error)
	Delete(ctx context.Context, id int) error
	Search(ctx context.Context, query string) ([]*User, error)
	SearchPaginated(ctx context.Context, query string, params PaginationParams) (*PaginatedResult[*User], error)
}

// CounterRepository defines the interface for counter state operations
type CounterRepository interface {
	Get(ctx context.Context) (int, error)
	Increment(ctx context.Context) (int, error)
	Decrement(ctx context.Context) (int, error)
	Reset(ctx context.Context) (int, error)
}

// Ensure our concrete types implement the interfaces at compile time
var (
	_ UserRepository    = (*UserStore)(nil)
	_ CounterRepository = (*CounterStore)(nil)
)