package controllers

import (
	"log"

	hazart "github.com/misbakhul29/hazartgo"
	"github.com/misbakhul29/hazartgo/_examples/models"
	"github.com/misbakhul29/hazartgo/crud"
	"gorm.io/gorm"
)

// ProductController menangani endpoint Product secara manual/eksplisit
type ProductController struct {
	db *gorm.DB
}

func NewProductController(db ...*gorm.DB) *ProductController {
	c := &ProductController{}
	if len(db) > 0 {
		c.db = db[0]
	}
	return c
}

func (pc *ProductController) RegisterRoutes(g *hazart.Group) {
	if pc.db != nil {
		// Real DB GORM AutoCRUD (GET /api/v1/products, GET/:id, POST, PUT/:id, DELETE/:id)
		crud.AutoCRUDGorm[models.Product](g, "", pc.db)
		return
	}

	// Manual/Mock Handlers
	// GET /api/v1/products
	hazart.GroupGet(g, "", pc.ListProducts, hazart.RouteMeta{
		Summary:     "List All Products",
		Description: "Mengambil daftar seluruh produk yang tersedia",
		Tags:        []string{"Products"},
	})

	// POST /api/v1/products
	hazart.GroupPost(g, "", pc.CreateProduct, hazart.RouteMeta{
		Summary:     "Create Product",
		Description: "Menambahkan produk baru ke katalog",
		Tags:        []string{"Products"},
		Security:    hazart.SecurityBearerAuth,
	})
}

func (pc *ProductController) ListProducts(ctx *hazart.Context, req *crud.EmptyReq) (*[]models.Product, error) {
	products := []models.Product{
		{ID: "prod_1", Name: "MacBook Pro M3", Price: 1999.99},
		{ID: "prod_2", Name: "Mechanical Keyboard", Price: 149.50},
	}
	return &products, nil
}

func (pc *ProductController) CreateProduct(ctx *hazart.Context, req *models.Product) (*models.Product, error) {
	req.ID = hazart.GenerateID()
	log.Printf("[ProductController] Successfully created product ID %s: %s ($%.2f)", req.ID, req.Name, req.Price)
	return req, nil
}
