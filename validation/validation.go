// Package validation provides input validation utilities for the HTMX learning application.
package validation

import (
	"errors"
	"net/mail"
	"strings"
	"unicode/utf8"
)

const (
	maxNameLength  = 100
	maxEmailLength = 254
	minNameLength  = 1
)

// ValidationError represents a validation error with field-specific information
type ValidationError struct {
	Field   string
	Message string
}

func (ve ValidationError) Error() string {
	return ve.Field + ": " + ve.Message
}

// ValidationErrors is a slice of validation errors
type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
	var messages []string
	for _, err := range ve {
		messages = append(messages, err.Error())
	}
	return strings.Join(messages, "; ")
}

// UserInput represents user input data for validation
type UserInput struct {
	Name  string
	Email string
}

// ValidateUser validates user input and returns any validation errors
func ValidateUser(input UserInput) error {
	var errors ValidationErrors

	// Validate name
	if nameErr := validateName(input.Name); nameErr != nil {
		errors = append(errors, ValidationError{Field: "name", Message: nameErr.Error()})
	}

	// Validate email
	if emailErr := validateEmail(input.Email); emailErr != nil {
		errors = append(errors, ValidationError{Field: "email", Message: emailErr.Error()})
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

// validateName validates the name field
func validateName(name string) error {
	name = strings.TrimSpace(name)

	if len(name) == 0 {
		return errors.New("name is required")
	}

	if utf8.RuneCountInString(name) < minNameLength {
		return errors.New("name is too short")
	}

	if utf8.RuneCountInString(name) > maxNameLength {
		return errors.New("name is too long (max 100 characters)")
	}

	// Check for potentially harmful characters (basic XSS prevention)
	if strings.ContainsAny(name, "<>\"'&") {
		return errors.New("name contains invalid characters")
	}

	return nil
}

// validateEmail validates the email field
func validateEmail(email string) error {
	email = strings.TrimSpace(email)

	if len(email) == 0 {
		return errors.New("email is required")
	}

	if len(email) > maxEmailLength {
		return errors.New("email is too long (max 254 characters)")
	}

	// Use Go's built-in email validation
	if _, err := mail.ParseAddress(email); err != nil {
		return errors.New("email format is invalid")
	}

	return nil
}

// SanitizeInput sanitizes string input by trimming whitespace and removing null bytes
func SanitizeInput(input string) string {
	// Remove null bytes that could cause issues
	input = strings.ReplaceAll(input, "\x00", "")
	// Trim whitespace
	return strings.TrimSpace(input)
}

