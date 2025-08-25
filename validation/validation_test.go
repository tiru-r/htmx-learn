package validation

import (
	"strings"
	"testing"
)

func TestValidateUser(t *testing.T) {
	tests := []struct {
		name      string
		input     UserInput
		wantError bool
		errorMsg  string
	}{
		{
			name:      "valid user",
			input:     UserInput{Name: "John Doe", Email: "john@example.com"},
			wantError: false,
		},
		{
			name:      "empty name",
			input:     UserInput{Name: "", Email: "john@example.com"},
			wantError: true,
			errorMsg:  "name is required",
		},
		{
			name:      "empty email",
			input:     UserInput{Name: "John Doe", Email: ""},
			wantError: true,
			errorMsg:  "email is required",
		},
		{
			name:      "invalid email",
			input:     UserInput{Name: "John Doe", Email: "not-an-email"},
			wantError: true,
			errorMsg:  "email format is invalid",
		},
		{
			name:      "name too long",
			input:     UserInput{Name: strings.Repeat("a", 101), Email: "john@example.com"},
			wantError: true,
			errorMsg:  "name is too long",
		},
		{
			name:      "name with invalid characters",
			input:     UserInput{Name: "John<script>", Email: "john@example.com"},
			wantError: true,
			errorMsg:  "name contains invalid characters",
		},
		{
			name:      "email too long",
			input:     UserInput{Name: "John Doe", Email: strings.Repeat("a", 250) + "@example.com"},
			wantError: true,
			errorMsg:  "email is too long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUser(tt.input)
			
			if tt.wantError && err == nil {
				t.Errorf("ValidateUser() expected error, got nil")
			}
			
			if !tt.wantError && err != nil {
				t.Errorf("ValidateUser() unexpected error: %v", err)
			}
			
			if tt.wantError && err != nil && tt.errorMsg != "" {
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("ValidateUser() error = %v, expected to contain %q", err, tt.errorMsg)
				}
			}
		})
	}
}

func TestSanitizeInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "normal string",
			input:    "Hello World",
			expected: "Hello World",
		},
		{
			name:     "string with whitespace",
			input:    "  Hello World  ",
			expected: "Hello World",
		},
		{
			name:     "string with null bytes",
			input:    "Hello\x00World",
			expected: "HelloWorld",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only whitespace",
			input:    "   ",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeInput(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeInput() = %q, expected %q", result, tt.expected)
			}
		})
	}
}