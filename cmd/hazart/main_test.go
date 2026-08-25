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
