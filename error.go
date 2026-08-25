package hazart

import "net/http"

// APIError represents a standardized RFC 7807 problem details error response
type APIError struct {
	Type     string `json:"type,omitempty" doc:"URI reference identifying problem type"`
	Title    string `json:"title" doc:"Short human-readable summary of problem"`
	Status   int    `json:"status" doc:"HTTP status code"`
	Detail   string `json:"detail,omitempty" doc:"Human-readable explanation specific to this occurrence"`
	Instance string `json:"instance,omitempty" doc:"URI reference identifying the specific occurrence"`
}

func (e *APIError) Error() string {
	if e.Detail != "" {
		return e.Detail
	}
	return e.Title
}

// NewAPIError creates a new APIError instance
func NewAPIError(status int, title string, detail string) *APIError {
	return &APIError{
		Title:  title,
		Status: status,
		Detail: detail,
	}
}

// Helper functions for common HTTP Errors

func BadRequest(detail string) *APIError {
	return NewAPIError(http.StatusBadRequest, "Bad Request", detail)
}

func Unauthorized(detail string) *APIError {
	return NewAPIError(http.StatusUnauthorized, "Unauthorized", detail)
}

func Forbidden(detail string) *APIError {
	return NewAPIError(http.StatusForbidden, "Forbidden", detail)
}

func NotFound(detail string) *APIError {
	return NewAPIError(http.StatusNotFound, "Not Found", detail)
}

func InternalServerError(detail string) *APIError {
	return NewAPIError(http.StatusInternalServerError, "Internal Server Error", detail)
}
