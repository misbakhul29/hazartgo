# ⚡ HazartGo

**HazartGo** adalah modern, high-performance, type-safe standalone Go web framework yang terinspirasi dari developer experience (DX) ala FastAPI & NestJS.

HazartGo menghadirkan **Zero-Boilerplate OpenAPI 3.1 & Swagger UI generation**, **AutoCRUD Model Engine**, **Built-in JWT & RBAC Engine**, **Response Standardizer**, serta **Hazart CLI Tool** untuk scaffolding project yang sangat cepat.

---

## ✨ Fitur Utama

- 🏎️ **High Performance Routing**: Radix Tree Router berbasis `net/http` standar Go (tanpa ketergantungan router eksternal).
- 🧬 **Type-Safe Generic Handlers**: Tidak perlu lagi decode/encode JSON (`json.NewDecoder`) secara manual.
- ⚡ **AutoCRUD Engine**: Cukup definisikan 1 struct Model DB, HazartGo otomatis membuatkan 5 REST API CRUD endpoints lengkap beserta OpenAPI spec-nya!
- 📖 **Comprehensive OpenAPI 3.1 & Swagger UI**: Generasi dokumentasi API otomatis di `/docs` & `/openapi.json` lengkap dengan `Contact`, `License`, `TermsOfService`, dan `SecuritySchemes`.
- 🔑 **Built-in JWT Package (`hazart/jwt`)**: Penandatanganan (`Sign`) dan verifikasi (`Verify`) token JWT berbasis HMAC-SHA256 bawaan tanpa dependensi pihak ketiga.
- 🛡️ **Role & Permission Access Control (RBAC)**: Guard middleware declaratif `middleware.RequireRole` & `middleware.RequirePermission` dengan Context Store (`ctx.Set` / `ctx.Get`).
- 🎯 **Standardized Response Helper**: Envelope JSON sukses dan error yang konsisten (`ctx.Success` & `ctx.Error`).
- 🚨 **RFC 7807 Problem Details**: Respon error terstruktur (`hazart.NotFound`, `hazart.BadRequest`, `hazart.Unauthorized`, dll).
- 🛠️ **Hazart CLI Tool**: Binary CLI (`hazart`) untuk scaffolding project baru, controller, dan resource berbasis AutoCRUD (`hazart init`, `hazart make:resource`).
- 🏗️ **Modular Controller Pattern**: Pemisahan route handler terisolasi via struct controller (`app.MountController`).
- 🛡️ **Built-in Middlewares**: Termasuk `Logger`, `Recovery`, `CORS`, `BearerAuth`, `RequireRole`, dan `RequirePermission`.

---

## 🚀 Quick Start

### 1. Installation

```bash
go get github.com/misbakhul29/hazartgo
```

### 2. Basic Example

```go
package main

import (
	"log"

	hazart "github.com/misbakhul29/hazartgo"
	"github.com/misbakhul29/hazartgo/crud"
	"github.com/misbakhul29/hazartgo/middleware"
	"github.com/misbakhul29/hazartgo/openapi"
)

// 1. Definisikan Struct Model DB (Sekaligus OpenAPI Spec & Validation!)
type User struct {
	ID    string `json:"id" path:"id" doc:"User ID"`
	Name  string `json:"name" validate:"required" doc:"User Full Name"`
	Email string `json:"email" validate:"required,email" doc:"User Email Address"`
}

// 2. Controller dengan AutoCRUD Engine + RBAC Protection
type UserController struct{}

func (uc *UserController) RegisterRoutes(g *hazart.Group) {
	// Auth middleware injection
	g.Use(func(next hazart.HandlerFunc) hazart.HandlerFunc {
		return func(ctx *hazart.Context) {
			ctx.Set("user_roles", []string{"admin"})
			next(ctx)
		}
	})

	// Proteksi seluruh route controller ini hanya untuk role "admin" atau "superadmin"
	g.Use(middleware.RequireRole("admin", "superadmin"))

	// Otomatis generate 5 REST API Endpoints:
	// GET /users, GET /users/:id, POST /users, PUT /users/:id, DELETE /users/:id
	crud.AutoCRUD[User](g, "")
}

func main() {
	app := hazart.New(hazart.Config{
		Title:          "My HazartGo API",
		Description:    "High Performance & Developer-Friendly Go REST API Framework",
		TermsOfService: "https://hazartgo.dev/terms",
		Contact: &openapi.Contact{
			Name:  "HazartGo Support",
			Email: "support@hazartgo.dev",
		},
		License: &openapi.License{
			Name: "MIT",
		},
		Version: "1.0.0",
	})

	// Built-in Middlewares
	app.Use(middleware.Logger())
	app.Use(middleware.Recovery())
	app.Use(middleware.CORS())

	// Configure OpenAPI Security Scheme
	app.OpenAPI.AddSecurityScheme(hazart.SecurityBearerAuth, openapi.SecurityScheme{
		Type:         "http",
		Scheme:       "bearer",
		BearerFormat: "JWT",
	})

	// Mount Controller
	app.MountController("/api/v1/users", &UserController{})

	log.Println("⚡ Server running on http://localhost:8080")
	log.Println("📚 Swagger UI Docs on http://localhost:8080/docs")

	app.Listen(":8080")
}
```

Jalankan server dan buka browser di:
- 🌐 **Swagger UI**: `http://localhost:8080/docs`
- 📑 **OpenAPI 3.1 Spec**: `http://localhost:8080/openapi.json`

---

## 🔑 Built-in JWT Package (`hazart/jwt`)

Paket JWT bawaan untuk mengelola pembuatan & verifikasi token secara ringkas:

```go
import "github.com/misbakhul29/hazartgo/jwt"

jwtManager := jwt.New("my-super-secret-key")

// 1. Generate JWT Token
token, err := jwtManager.Sign(jwt.MapClaims{
    "sub":   "user_123",
    "name":  "Budi Santoso",
    "roles": []string{"admin"},
}, 24*time.Hour)

// 2. Verifikasi Token
claims, err := jwtManager.Verify(token)
if err != nil {
    log.Println("Token tidak valid / expired")
}
```

---

## 🛡️ Role & Permission Access Control (RBAC)

Memproteksi route berdasarkan **Role** atau **Permission** pengguna menggunakan Context Store (`ctx.Set` & `ctx.Get`):

```go
func (c *ProductController) RegisterRoutes(g *hazart.Group) {
    // 1. Set role/permission dari Auth Middleware
    g.Use(func(next hazart.HandlerFunc) hazart.HandlerFunc {
        return func(ctx *hazart.Context) {
            ctx.Set("user_roles", []string{"editor"})
            ctx.Set("user_permissions", []string{"products:read", "products:write"})
            next(ctx)
        }
    })

    // 2. Proteksi berdasarkan Role
    g.Use(middleware.RequireRole("admin", "editor"))

    // 3. Atau Proteksi berdasarkan Specific Permission
    // g.Use(middleware.RequirePermission("products:write"))
}
```

---

## 🎯 Response Standardizer Helper

Gunakan `ctx.Success` dan `ctx.Error` untuk memastikan format JSON response konsisten di seluruh aplikasi:

```go
// Response Success
ctx.Success(http.StatusOK, userObj, "User data retrieved successfully")
// JSON Output: {"success": true, "message": "User data retrieved successfully", "data": {...}}

// Response Error
ctx.Error(http.StatusBadRequest, "Invalid parameter", map[string]string{"field": "email"})
// JSON Output: {"success": false, "error": "Invalid parameter", "details": {"field": "email"}}
```

---

## 🛠️ Hazart CLI Scaffolding Tool

Instal CLI Tool HazartGo ke dalam sistem Go kamu:

```bash
go install github.com/misbakhul29/hazartgo/cmd/hazart@latest
```

### Perintah CLI:

#### 1. Inisialisasi Project Baru
```bash
hazart init my-backend-app
```

#### 2. Generate Controller Baru
```bash
hazart make:controller Auth --group "/api/v1"
```

#### 3. Generate Complete Resource (Model + AutoCRUD Controller)
```bash
hazart make:resource Product --group "/api/v1"
```

---

## 📄 License

MIT License © 2026 Misbakhul Munir
