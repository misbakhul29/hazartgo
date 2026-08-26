package main

import (
	"log"
	"time"

	hazart "github.com/misbakhul29/hazartgo"
	"github.com/misbakhul29/hazartgo/_examples/controllers"
	"github.com/misbakhul29/hazartgo/_examples/models"
	"github.com/misbakhul29/hazartgo/_examples/repositories"
	"github.com/misbakhul29/hazartgo/crud"
	"github.com/misbakhul29/hazartgo/middleware"
	"github.com/misbakhul29/hazartgo/openapi"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func main() {
	// 1. Inisialisasi App Framework
	app := hazart.New(hazart.Config{
		Title:          "HazartGo Showcase API",
		Description:    "High Performance & Developer-Friendly Go REST API Framework with Automatic OpenAPI Specs, Validation & Enterprise Features",
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
		Version: "1.2.2",
	})

	// 2. Global Middlewares (Tracing, Structured Logging, Recovery, CORS, Rate Limiting, Response Caching)
	app.Use(middleware.Tracing())             // Distributed Trace-ID (X-Trace-ID)
	app.Use(middleware.StructuredLogger(nil)) // Structured JSON Logger (log/slog)
	app.Use(middleware.Recovery())            // Panic recovery
	app.Use(middleware.CORS())                // Cross-Origin Resource Sharing
	app.Use(middleware.RateLimiter(middleware.RateLimitConfig{
		RequestsPerWindow: 120,
		Window:            time.Minute,
	})) // 120 requests/minute per client IP
	app.Use(middleware.Cache(30 * time.Second)) // In-memory GET response cacher (30s TTL)

	// 3. Configure OpenAPI Security Scheme
	app.OpenAPI.AddSecurityScheme(hazart.SecurityBearerAuth, openapi.SecurityScheme{
		Type:         "http",
		Scheme:       "bearer",
		BearerFormat: "JWT",
	})

	// 4. Inisialisasi Real GORM Database Connection & AutoCRUD Endpoint
	dbConfig := hazart.DBConfig{
		Driver:      "sqlite", // Beralih ke "postgres" atau "mysql" sesuai kebutuhan
		File:        "hazart_demo.db",
		AutoMigrate: []any{&models.Product{}},
	}

	db, err := gorm.Open(sqlite.Open(dbConfig.BuildDSN()), &gorm.Config{})
	if err != nil {
		log.Printf("⚠️ GORM DB Connection notice: %v", err)
	} else {
		// Auto Migrate Table Schema
		for _, model := range dbConfig.AutoMigrate {
			_ = db.AutoMigrate(model)
		}
		log.Println("✅ Real GORM Database connected & schema migrated successfully!")

		// Mount 5 Real DB REST API CRUD Endpoints (GET, GET/:id, POST, PUT/:id, DELETE/:id)
		apiV1 := app.Group("/api/v1")
		crud.AutoCRUDGorm[models.Product](apiV1, "/gorm-products", db)
		log.Println("🚀 Registered Real DB AutoCRUD Endpoints on /api/v1/gorm-products")
	}

	// 5. Inisialisasi Dependency (Repositories & Controllers)
	userRepo := repositories.NewUserMemoryRepository()

	userController := controllers.NewUserController(userRepo)
	productController := controllers.NewProductController()
	featureController := controllers.NewFeatureController()

	// 6. Mount Controllers ke Routing Group Prefix
	app.MountController("/api/v1/users", userController)
	app.MountController("/api/v1/products", productController)
	app.MountController("/api/v1/features", featureController)

	log.Println("⚡ HazartGo Framework running on http://localhost:8080")
	log.Println("📚 Swagger UI Documentation available on http://localhost:8080/docs")

	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
