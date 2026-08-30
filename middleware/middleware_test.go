package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/misbakhul29/hazartgo"
	"github.com/misbakhul29/hazartgo/middleware"
)

func TestCORS_Preflight(t *testing.T) {
	app := hazart.New(hazart.Config{})
	app.Use(middleware.CORS(middleware.CORSOptions{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
	}))

	req := httptest.NewRequest("OPTIONS", "/api/v1/items", nil)
	req.Header.Set("Access-Control-Request-Method", "GET")
	req.Header.Set("Origin", "http://localhost:3000")

	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected status 204 No Content for OPTIONS preflight, got %d", w.Code)
	}

	if origin := w.Header().Get("Access-Control-Allow-Origin"); origin != "http://localhost:3000" {
		t.Errorf("expected Access-Control-Allow-Origin 'http://localhost:3000', got '%s'", origin)
	}

	if methods := w.Header().Get("Access-Control-Allow-Methods"); methods != "GET, POST, PUT, DELETE, OPTIONS" {
		t.Errorf("expected Access-Control-Allow-Methods, got '%s'", methods)
	}

	if headers := w.Header().Get("Access-Control-Allow-Headers"); headers != "Content-Type, Authorization" {
		t.Errorf("expected Access-Control-Allow-Headers, got '%s'", headers)
	}
}
