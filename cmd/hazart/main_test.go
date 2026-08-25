package main

import (
	"testing"
)

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"UserController", "user_controller"},
		{"UserProfile", "user_profile"},
		{"user_profile", "user_profile"},
		{"User", "user"},
		{"", ""},
	}

	for _, tt := range tests {
		got := toSnakeCase(tt.input)
		if got != tt.expected {
			t.Errorf("toSnakeCase(%q) = %q; want %q", tt.input, got, tt.expected)
		}
	}
}

func TestToPascalCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"user_controller", "UserController"},
		{"user-profile", "UserProfile"},
		{"user", "User"},
		{"UserProfile", "UserProfile"},
		{"", ""},
	}

	for _, tt := range tests {
		got := toPascalCase(tt.input)
		if got != tt.expected {
			t.Errorf("toPascalCase(%q) = %q; want %q", tt.input, got, tt.expected)
		}
	}
}

func TestCapitalize(t *testing.T) {
	if capitalize("user") != "User" {
		t.Errorf("capitalize('user') failed")
	}
	if capitalize("") != "" {
		t.Errorf("capitalize('') failed")
	}
}

func TestIsNewerVersion(t *testing.T) {
	tests := []struct {
		latest   string
		current  string
		expected bool
	}{
		{"v1.3.0", "v1.2.0", true},
		{"v1.2.1", "v1.2.0", true},
		{"v2.0.0", "v1.2.0", true},
		{"v1.2.0", "v1.2.0", false},
		{"v1.1.0", "v1.2.0", false},
		{"v0.0.0-20260825090345-53844ebf15b4", "v1.2.0", false},
		{"", "v1.2.0", false},
	}

	for _, tt := range tests {
		got := isNewerVersion(tt.latest, tt.current)
		if got != tt.expected {
			t.Errorf("isNewerVersion(%q, %q) = %v; want %v", tt.latest, tt.current, got, tt.expected)
		}
	}
}
