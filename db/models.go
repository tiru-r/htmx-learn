// Package db provides database connection management and data access operations
// for the HTMX learning application using PostgreSQL with pgx driver.
package db

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

const (
	// CounterID represents the single counter state record ID
	counterID = 1
)

// User represents a user in the database
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CounterState represents the counter state in the database
type CounterState struct {
	ID        int       `json:"id"`
	Count     int       `json:"count"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserStore provides database operations for users
type UserStore struct {
	db *DB
}

// NewUserStore creates a new UserStore
func NewUserStore(db *DB) *UserStore {
	return &UserStore{db: db}
}

// GetAll retrieves all users from the database
func (us *UserStore) GetAll(ctx context.Context) ([]*User, error) {
	query := "SELECT id, name, email, created_at, updated_at FROM users ORDER BY created_at DESC"
	rows, err := us.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user rows: %w", err)
	}

	return users, nil
}


// Add creates a new user in the database
func (us *UserStore) Add(ctx context.Context, name, email string) (*User, error) {
	query := "INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id, name, email, created_at, updated_at"
	row := us.db.Pool.QueryRow(ctx, query, name, email)

	user := &User{}
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create user %s <%s>: %w", name, email, err)
	}

	return user, nil
}


// Delete removes a user from the database
func (us *UserStore) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM users WHERE id = $1"
	result, err := us.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user ID %d: %w", id, err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

// Search finds users by name or email
func (us *UserStore) Search(ctx context.Context, query string) ([]*User, error) {
	sqlQuery := `
		SELECT id, name, email, created_at, updated_at 
		FROM users 
		WHERE name ILIKE $1 OR email ILIKE $1 
		ORDER BY created_at DESC
	`
	searchTerm := "%" + strings.ToLower(query) + "%"
	rows, err := us.db.Query(ctx, sqlQuery, searchTerm)
	if err != nil {
		return nil, fmt.Errorf("failed to search users with query '%s': %w", query, err)
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan search result: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating search results: %w", err)
	}

	return users, nil
}

// SearchPaginated finds users by name or email with pagination
func (us *UserStore) SearchPaginated(ctx context.Context, query string, params PaginationParams) (*PaginatedResult[*User], error) {
	// First get the total count for search results
	countQuery := `
		SELECT COUNT(*) 
		FROM users 
		WHERE name ILIKE $1 OR email ILIKE $1
	`
	searchTerm := "%" + strings.ToLower(query) + "%"
	row := us.db.Pool.QueryRow(ctx, countQuery, searchTerm)
	
	var total int
	if err := row.Scan(&total); err != nil {
		return nil, fmt.Errorf("failed to count search results for query '%s': %w", query, err)
	}

	// Get the paginated search results
	sqlQuery := `
		SELECT id, name, email, created_at, updated_at 
		FROM users 
		WHERE name ILIKE $1 OR email ILIKE $1 
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := us.db.Query(ctx, sqlQuery, searchTerm, params.PageSize, params.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search users with query '%s': %w", query, err)
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan search result: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating search results: %w", err)
	}

	result := NewPaginatedResult(users, params, total)
	return result, nil
}

// GetAllPaginated retrieves users with pagination
func (us *UserStore) GetAllPaginated(ctx context.Context, params PaginationParams) (*PaginatedResult[*User], error) {
	// First get the total count
	total, err := us.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count users for pagination: %w", err)
	}

	// Get the paginated data
	query := "SELECT id, name, email, created_at, updated_at FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2"
	rows, err := us.db.Query(ctx, query, params.PageSize, params.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query paginated users: %w", err)
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan paginated user row: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating paginated user rows: %w", err)
	}

	result := NewPaginatedResult(users, params, total)
	return result, nil
}

// Count returns the total number of users
func (us *UserStore) Count(ctx context.Context) (int, error) {
	query := "SELECT COUNT(*) FROM users"
	row := us.db.Pool.QueryRow(ctx, query)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	return count, nil
}

// CounterStore provides database operations for counter state
type CounterStore struct {
	db *DB
}

// NewCounterStore creates a new CounterStore
func NewCounterStore(db *DB) *CounterStore {
	return &CounterStore{db: db}
}

// Get retrieves the current counter value
func (cs *CounterStore) Get(ctx context.Context) (int, error) {
	query := "SELECT count FROM counter_state WHERE id = $1"
	row := cs.db.Pool.QueryRow(ctx, query, counterID)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get counter value: %w", err)
	}

	return count, nil
}


// Increment increases the counter by 1
func (cs *CounterStore) Increment(ctx context.Context) (int, error) {
	query := "UPDATE counter_state SET count = count + 1 WHERE id = $1 RETURNING count"
	row := cs.db.Pool.QueryRow(ctx, query, counterID)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to increment counter: %w", err)
	}

	return count, nil
}

// Decrement decreases the counter by 1
func (cs *CounterStore) Decrement(ctx context.Context) (int, error) {
	query := "UPDATE counter_state SET count = count - 1 WHERE id = $1 RETURNING count"
	row := cs.db.Pool.QueryRow(ctx, query, counterID)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to decrement counter: %w", err)
	}

	return count, nil
}

// Reset sets the counter to 0
func (cs *CounterStore) Reset(ctx context.Context) (int, error) {
	query := "UPDATE counter_state SET count = 0 WHERE id = $1 RETURNING count"
	row := cs.db.Pool.QueryRow(ctx, query, counterID)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to reset counter: %w", err)
	}

	return count, nil
}