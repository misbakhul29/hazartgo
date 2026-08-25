package main

import (
	"log"

	hazart "github.com/misbakhul29/hazartgo"
	"github.com/misbakhul29/hazartgo/_examples/controllers"
	"github.com/misbakhul29/hazartgo/_examples/repositories"
	"github.com/misbakhul29/hazartgo/middleware"
	"github.com/misbakhul29/hazartgo/openapi"
)

func main() {
	// 1. Inisialisasi App Framework
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

	// 2. Global Middlewares
	app.Use(middleware.Logger())
	app.Use(middleware.Recovery())
	app.Use(middleware.CORS())

	// 3. Configure OpenAPI Security Scheme
	app.OpenAPI.AddSecurityScheme(hazart.SecurityBearerAuth, openapi.SecurityScheme{
		Type:         "http",
		Scheme:       "bearer",
		BearerFormat: "JWT",
	})

	// 4. Inisialisasi Dependency (Repository & Controllers)
	userRepo := repositories.NewUserMemoryRepository()

	userController := controllers.NewUserController(userRepo)
	productController := controllers.NewProductController()

	// 5. Mount Controllers ke Routing Group Prefix
	app.MountController("/api/v1/users", userController)
	app.MountController("/api/v1/products", productController)

	log.Println("⚡ HazartGo Framework running on http://localhost:8080")
	log.Println("📚 Swagger UI Documentation available on http://localhost:8080/docs")

	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
