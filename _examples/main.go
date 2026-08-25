package main

import (
	"log"

	hazart "github.com/misbakhul29/hazartgo"
	"github.com/misbakhul29/hazartgo/crud"
	"github.com/misbakhul29/hazartgo/middleware"
	"github.com/misbakhul29/hazartgo/openapi"
)

// DTO Requests & Responses

// Model Struct DB (Sekaligus OpenAPI Spec & Validation)
type User struct {
	ID    string `json:"id" path:"id" doc:"User ID"`
	Name  string `json:"name" validate:"required" doc:"User Full Name"`
	Email string `json:"email" validate:"required,email" doc:"User Email Address"`
}

// Controllers Pattern (Menggunakan AutoCRUD)
type UserController struct{}

func (uc *UserController) RegisterRoutes(g *hazart.Group) {
	// Otomatis membuat 5 endpoint CRUD lengkap (GET /users, GET /users/:id, POST /users, PUT /users/:id, DELETE /users/:id)
	crud.AutoCRUD[User](g, "")
}

func main() {
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

	// Register Built-in Middlewares
	app.Use(middleware.Logger())
	app.Use(middleware.Recovery())
	app.Use(middleware.CORS())

	// Configure OpenAPI Security Scheme (JWT Bearer Auth)
	app.OpenAPI.AddSecurityScheme(hazart.SecurityBearerAuth, openapi.SecurityScheme{
		Type:         "http",
		Scheme:       "bearer",
		BearerFormat: "JWT",
	})

	// Modular Controller Mount
	app.MountController("/api/v1/users", &UserController{})

	log.Println("⚡ HazartGo Framework running on http://localhost:8080")
	log.Println("📚 Swagger UI Documentation available on http://localhost:8080/docs")

	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
