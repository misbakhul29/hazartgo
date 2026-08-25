package controllers

import (
	hazart "github.com/misbakhul29/hazartgo"
	"github.com/misbakhul29/hazartgo/_examples/models"
	"github.com/misbakhul29/hazartgo/_examples/repositories"
	"github.com/misbakhul29/hazartgo/crud"
	"github.com/misbakhul29/hazartgo/middleware"
)

// UserController menangani endpoint User dengan AutoCRUD Engine
type UserController struct {
	repo repositories.UserRepository
}

func NewUserController(repo repositories.UserRepository) *UserController {
	return &UserController{repo: repo}
}

func (uc *UserController) RegisterRoutes(g *hazart.Group) {
	// 1. Auth middleware: Injeksi role & permission ke Context Store
	g.Use(func(next hazart.HandlerFunc) hazart.HandlerFunc {
		return func(ctx *hazart.Context) {
			ctx.Set("user_roles", []string{"admin"})
			ctx.Set("user_permissions", []string{"users:read", "users:write"})
			next(ctx)
		}
	})

	// 2. Role Protection: Hanya role "admin" atau "superadmin" yang bisa akses
	g.Use(middleware.RequireRole("admin", "superadmin"))

	// 3. AutoCRUD Engine: Otomatis generate 5 REST API CRUD Endpoints (GET, GET/:id, POST, PUT/:id, DELETE/:id)
	crud.AutoCRUD[models.User](g, "", uc.repo)
}
