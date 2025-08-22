package db

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
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
func (us *UserStore) GetAll() ([]*User, error) {
	query := "SELECT id, name, email, created_at, updated_at FROM users ORDER BY created_at DESC"
	rows, err := us.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
}

// GetByID retrieves a user by ID
func (us *UserStore) GetByID(id int) (*User, error) {
	query := "SELECT id, name, email, created_at, updated_at FROM users WHERE id = ?"
	row := us.db.QueryRow(query, id)

	user := &User{}
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Add creates a new user in the database
func (us *UserStore) Add(name, email string) (*User, error) {
	query := "INSERT INTO users (name, email) VALUES (?, ?) RETURNING id, name, email, created_at, updated_at"
	row := us.db.QueryRow(query, name, email)

	user := &User{}
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Update modifies an existing user
func (us *UserStore) Update(id int, name, email string) (*User, error) {
	query := "UPDATE users SET name = ?, email = ? WHERE id = ? RETURNING id, name, email, created_at, updated_at"
	row := us.db.QueryRow(query, name, email, id)

	user := &User{}
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Delete removes a user from the database
func (us *UserStore) Delete(id int) error {
	query := "DELETE FROM users WHERE id = ?"
	result, err := us.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Search finds users by name or email
func (us *UserStore) Search(query string) ([]*User, error) {
	sqlQuery := `
		SELECT id, name, email, created_at, updated_at 
		FROM users 
		WHERE name LIKE ? OR email LIKE ? 
		ORDER BY created_at DESC
	`
	searchTerm := "%" + strings.ToLower(query) + "%"
	rows, err := us.db.Query(sqlQuery, searchTerm, searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
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
func (cs *CounterStore) Get() (int, error) {
	query := "SELECT count FROM counter_state WHERE id = 1"
	row := cs.db.QueryRow(query)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// Set updates the counter value
func (cs *CounterStore) Set(count int) error {
	query := "UPDATE counter_state SET count = ? WHERE id = 1"
	result, err := cs.db.Exec(query, count)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("counter state not found")
	}

	return nil
}

// Increment increases the counter by 1
func (cs *CounterStore) Increment() (int, error) {
	query := "UPDATE counter_state SET count = count + 1 WHERE id = 1 RETURNING count"
	row := cs.db.QueryRow(query)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// Decrement decreases the counter by 1
func (cs *CounterStore) Decrement() (int, error) {
	query := "UPDATE counter_state SET count = count - 1 WHERE id = 1 RETURNING count"
	row := cs.db.QueryRow(query)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// Reset sets the counter to 0
func (cs *CounterStore) Reset() (int, error) {
	query := "UPDATE counter_state SET count = 0 WHERE id = 1 RETURNING count"
	row := cs.db.QueryRow(query)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}