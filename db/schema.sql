-- HTMX Learn Database Schema
-- SQLite database for the HTMX + Go application

-- Users table for the dynamic content examples
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Counter state table for persistence
CREATE TABLE IF NOT EXISTS counter_state (
    id INTEGER PRIMARY KEY CHECK (id = 1), -- Single row constraint
    count INTEGER NOT NULL DEFAULT 0,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Insert initial counter state
INSERT OR IGNORE INTO counter_state (id, count) VALUES (1, 0);

-- Insert some sample users
INSERT OR IGNORE INTO users (name, email) VALUES
    ('John Doe', 'john@example.com'),
    ('Jane Smith', 'jane@example.com'),
    ('Bob Johnson', 'bob@example.com');

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_name ON users(name);

-- Trigger to update updated_at timestamp
CREATE TRIGGER IF NOT EXISTS update_users_timestamp 
    AFTER UPDATE ON users
BEGIN
    UPDATE users SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_counter_timestamp 
    AFTER UPDATE ON counter_state
BEGIN
    UPDATE counter_state SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;