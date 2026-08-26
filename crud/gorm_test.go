package crud

import (
	"testing"
)

func TestSanitizeIdentifier(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"created_at", "created_at"},
		{"name; DROP TABLE users;--", "nameDROPTABLEusers"},
		{"user_id'", "user_id"},
		{"id OR 1=1", "idOR11"},
		{"ORDER BY 1", "ORDERBY1"},
	}

	for _, tt := range tests {
		result := sanitizeIdentifier(tt.input)
		if result != tt.expected {
			t.Errorf("sanitizeIdentifier(%q) = %q; expected %q", tt.input, result, tt.expected)
		}
	}
}

func TestSanitizeLikePattern(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"normal_text", "normal\\_text"},
		{"100%", "100\\%"},
		{"user\\name", "user\\\\name"},
	}

	for _, tt := range tests {
		result := sanitizeLikePattern(tt.input)
		if result != tt.expected {
			t.Errorf("sanitizeLikePattern(%q) = %q; expected %q", tt.input, result, tt.expected)
		}
	}
}
