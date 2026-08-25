package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	hazart "github.com/misbakhul29/hazartgo"
	"github.com/misbakhul29/hazartgo/_examples/controllers"
	"github.com/misbakhul29/hazartgo/_examples/models"
	"github.com/misbakhul29/hazartgo/_examples/repositories"
	"github.com/misbakhul29/hazartgo/crud"
	"github.com/misbakhul29/hazartgo/middleware"
	"github.com/misbakhul29/hazartgo/openapi"
)

func setupExampleApp() *hazart.App {
	app := hazart.New(hazart.Config{
		Title:          "HazartGo Showcase API",
		Description:    "High Performance & Developer-Friendly Go REST API Framework with Automatic OpenAPI Specs & Validation",
		TermsOfService: "https://hazartgo.dev/terms",
		Contact: &openapi.Contact{
			Name:  "HazartGo Support",
			URL:   "https://github.com/misbakhul29/hazartgo",
			Email: "support@hazartgo.dev",
		},
		License: &openapi.License{
			Name: "MIT License",
			URL:  "https://opensource.org/licenses/MIT",
		},
		Version: "1.0.0",
	})

	app.Use(middleware.Logger())
	app.Use(middleware.Recovery())
	app.Use(middleware.CORS())

	app.OpenAPI.AddSecurityScheme(hazart.SecurityBearerAuth, openapi.SecurityScheme{
		Type:         "http",
		Scheme:       "bearer",
		BearerFormat: "JWT",
	})

	userRepo := repositories.NewUserMemoryRepository()
	app.MountController("/api/v1/users", controllers.NewUserController(userRepo))
	app.MountController("/api/v1/products", controllers.NewProductController())

	return app
}

func TestMainExample_FullWorkflow(t *testing.T) {
	app := setupExampleApp()

	// 1. Test Docs & OpenAPI JSON endpoints
	t.Run("GET /docs", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/docs", nil)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200 for /docs, got %d", rec.Code)
		}
	})

	t.Run("GET /openapi.json", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/openapi.json", nil)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200 for /openapi.json, got %d", rec.Code)
		}

		var spec openapi.Spec
		if err := json.Unmarshal(rec.Body.Bytes(), &spec); err != nil {
			t.Fatalf("invalid json spec: %v", err)
		}

		if spec.Info.Title != "HazartGo Showcase API" {
			t.Errorf("unexpected spec title: %s", spec.Info.Title)
		}
		if spec.Info.Contact.Email != "support@hazartgo.dev" {
			t.Errorf("unexpected spec contact email: %s", spec.Info.Contact.Email)
		}
	})

	var createdID string

	// 2. Test POST /api/v1/users (Create User via AutoCRUD)
	t.Run("POST /api/v1/users", func(t *testing.T) {
		newUser := models.User{
			Name:  "Budi Santoso",
			Email: "budi@example.com",
		}
		body, _ := json.Marshal(newUser)

		req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d. Body: %s", rec.Code, rec.Body.String())
		}

		var res models.User
		if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if res.ID == "" {
			t.Error("expected generated ID, got empty string")
		}
		if res.Name != "Budi Santoso" || res.Email != "budi@example.com" {
			t.Errorf("unexpected user data: %+v", res)
		}

		createdID = res.ID
	})

	// 3. Test GET /api/v1/users (List All)
	t.Run("GET /api/v1/users", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/users", nil)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", rec.Code)
		}

		var users []models.User
		if err := json.Unmarshal(rec.Body.Bytes(), &users); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if len(users) != 1 {
			t.Errorf("expected 1 user, got %d", len(users))
		}
	})

	// 4. Test GET /api/v1/users/:id (Find By ID)
	t.Run("GET /api/v1/users/:id", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/users/"+createdID, nil)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", rec.Code)
		}

		var u models.User
		if err := json.Unmarshal(rec.Body.Bytes(), &u); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if u.ID != createdID || u.Name != "Budi Santoso" {
			t.Errorf("unexpected user record: %+v", u)
		}
	})

	// 5. Test PUT /api/v1/users/:id (Update)
	t.Run("PUT /api/v1/users/:id", func(t *testing.T) {
		updatedData := models.User{
			ID:    createdID,
			Name:  "Budi Santoso Updated",
			Email: "budi.updated@example.com",
		}
		body, _ := json.Marshal(updatedData)

		req := httptest.NewRequest("PUT", "/api/v1/users/"+createdID, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", rec.Code)
		}
	})

	// 6. Test DELETE /api/v1/users/:id (Delete)
	t.Run("DELETE /api/v1/users/:id", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/v1/users/"+createdID, nil)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", rec.Code)
		}
	})
}

func TestRequireRole_AccessForbidden(t *testing.T) {
	app := hazart.New(hazart.Config{Title: "Test RBAC"})

	g := app.Group("/api/v1/protected")
	g.Use(func(next hazart.HandlerFunc) hazart.HandlerFunc {
		return func(ctx *hazart.Context) {
			ctx.Set("user_roles", []string{"user"})
			next(ctx)
		}
	})
	g.Use(middleware.RequireRole("admin"))

	crud.AutoCRUD[models.User](g, "")

	req := httptest.NewRequest("GET", "/api/v1/protected", nil)
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected status 403 Forbidden for non-admin role, got %d", rec.Code)
	}
}
