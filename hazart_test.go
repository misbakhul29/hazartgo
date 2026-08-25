package hazart_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/misbakhul29/hazartgo"
	"github.com/misbakhul29/hazartgo/middleware"
)

type TestReq struct {
	ID string `path:"id"`
}

type TestRes struct {
	Message string `json:"message"`
}

func TestHazartFramework(t *testing.T) {
	app := hazart.New(hazart.Config{
		Title:   "Test App",
		Version: "1.0.0",
	})

	hazart.Get(app, "/hello/:id", func(ctx *hazart.Context, req *TestReq) (*TestRes, error) {
		return &TestRes{
			Message: "Hello " + req.ID,
		}, nil
	})

	// Test GET Route
	req := httptest.NewRequest("GET", "/hello/world", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	expectedBody := `{"message":"Hello world"}`
	if strings.TrimSpace(w.Body.String()) != expectedBody {
		t.Fatalf("expected body %s, got %s", expectedBody, w.Body.String())
	}

	// Test OpenAPI Docs Route
	reqDocs := httptest.NewRequest("GET", "/openapi.json", nil)
	wDocs := httptest.NewRecorder()
	app.ServeHTTP(wDocs, reqDocs)

	if wDocs.Code != http.StatusOK {
		t.Fatalf("expected docs status 200, got %d", wDocs.Code)
	}

	// Test Middlewares (Recovery & CORS)
	app.Use(hazart.MiddlewareFunc(func(next hazart.HandlerFunc) hazart.HandlerFunc {
		return func(ctx *hazart.Context) {
			ctx.SetHeader("X-Custom-Middleware", "Active")
			next(ctx)
		}
	}))

	// Test Route Grouping & Group Middleware Execution
	v1 := app.Group("/api/v1")
	v1.Use(middleware.BearerAuth(func(token string) bool {
		return token == "valid-token"
	}))

	hazart.GroupGet(v1, "/ping", func(ctx *hazart.Context, req *TestReq) (*TestRes, error) {
		return &TestRes{Message: "pong"}, nil
	})

	// Test 1: Unauthenticated request (Missing Token) -> Expect 401
	reqUnauth := httptest.NewRequest("GET", "/api/v1/ping", nil)
	wUnauth := httptest.NewRecorder()
	app.ServeHTTP(wUnauth, reqUnauth)

	if wUnauth.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401 Unauthorized for missing token, got %d", wUnauth.Code)
	}

	// Test 2: Authenticated request (Valid Token) -> Expect 200
	reqAuth := httptest.NewRequest("GET", "/api/v1/ping", nil)
	reqAuth.Header.Set("Authorization", "Bearer valid-token")
	wAuth := httptest.NewRecorder()
	app.ServeHTTP(wAuth, reqAuth)

	if wAuth.Code != http.StatusOK {
		t.Fatalf("expected status 200 OK for valid token, got %d", wAuth.Code)
	}
}
