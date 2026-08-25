# ⚡ HazartGo

**HazartGo** adalah modern, high-performance, type-safe standalone Go web framework yang terinspirasi dari developer experience (DX) ala FastAPI & NestJS.

HazartGo menghadirkan **Zero-Boilerplate OpenAPI 3.1 & Swagger UI generation**, **AutoCRUD Model Engine**, **Automatic Request Binding & Validation**, serta **Hazart CLI Tool** untuk scaffolding project yang sangat cepat.

---

## ✨ Fitur Utama

- 🏎️ **High Performance Routing**: Radix Tree Router berbasis `net/http` standar Go (tanpa ketergantungan router eksternal).
- 🧬 **Type-Safe Generic Handlers**: Tidak perlu lagi decode/encode JSON (`json.NewDecoder`) secara manual.
- ⚡ **AutoCRUD Engine**: Cukup definisikan 1 struct Model DB, HazartGo otomatis membuatkan 5 REST API CRUD endpoints lengkap beserta OpenAPI spec-nya!
- 📖 **Zero-Boilerplate OpenAPI 3.1 & Swagger UI**: Generasi dokumentasi API otomatis di `/docs` & `/openapi.json` tanpa perlu menulis file YAML manual.
- 🔒 **Type-Safe Security Schemes**: Dukungan `hazart.SecurityBearerAuth` untuk tombol *Authorize* interaktif di Swagger UI & `middleware.BearerAuth` untuk runtime security guard.
- 🚨 **RFC 7807 Problem Details**: Respon error terstruktur yang konsisten (`hazart.NotFound`, `hazart.BadRequest`, dll).
- 🛠️ **Hazart CLI Tool**: Binary CLI (`hazart`) untuk scaffolding project baru, controller, dan resource berbasis AutoCRUD (`hazart init`, `hazart make:resource`).
- 🏗️ **Modular Controller Pattern**: Pemisahan route handler terisolasi via struct controller (`app.MountController`).
- 🛡️ **Built-in Middlewares**: Termasuk `Logger`, `Recovery`, dan `CORS`.

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
)

// 1. Definisikan Struct Model DB (Sekaligus OpenAPI Spec & Validation!)
type User struct {
	ID    string `json:"id" path:"id" doc:"User ID"`
	Name  string `json:"name" validate:"required" doc:"User Full Name"`
	Email string `json:"email" validate:"required,email" doc:"User Email Address"`
}

// 2. Controller dengan AutoCRUD Engine
type UserController struct{}

func (uc *UserController) RegisterRoutes(g *hazart.Group) {
	// Otomatis generate 5 REST API Endpoints:
	// GET /users, GET /users/:id, POST /users, PUT /users/:id, DELETE /users/:id
	crud.AutoCRUD[User](g, "")
}

func main() {
	app := hazart.New(hazart.Config{
		Title:   "My HazartGo API",
		Version: "1.0.0",
	})

	// Built-in Middlewares
	app.Use(middleware.Logger())
	app.Use(middleware.Recovery())
	app.Use(middleware.CORS())

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

## 🔒 Security & Bearer Auth

```go
// 1. Registrasi Security Scheme di App
app.OpenAPI.AddSecurityScheme(hazart.SecurityBearerAuth, openapi.SecurityScheme{
    Type:         "http",
    Scheme:       "bearer",
    BearerFormat: "JWT",
})

// 2. Pasang Runtime Middleware & OpenAPI Tag di Group/Controller
func (c *UserController) RegisterRoutes(g *hazart.Group) {
    g.Use(middleware.BearerAuth(func(token string) bool {
        return token == "secret-jwt-token"
    }))

    hazart.GroupGet(g, "/profile", GetProfile, hazart.RouteMeta{
        Summary:  "Get User Profile",
        Security: hazart.SecurityBearerAuth,
    })
}
```

---

## 🚨 Structured Error Response (RFC 7807)

HazartGo menyediakan helper error standar terstruktur:

```go
func GetUser(ctx *hazart.Context, req *UserReq) (*UserRes, error) {
    if req.ID == "404" {
        return nil, hazart.NotFound("User ID not found")
    }
    return &UserRes{ID: req.ID}, nil
}
```

**JSON Output (`404 Not Found`)**:
```json
{
  "title": "Not Found",
  "status": 404,
  "detail": "User ID not found",
  "instance": "/api/v1/users/404"
}
```

---

## 📄 License

MIT License © 2026 Misbakhul Munir
