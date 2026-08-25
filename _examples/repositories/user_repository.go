package repositories

import (
	hazart "github.com/misbakhul29/hazartgo"
	"github.com/misbakhul29/hazartgo/_examples/models"
)

// UserRepository defines custom repository interface for User
type UserRepository interface {
	FindAll(ctx *hazart.Context) ([]models.User, error)
	FindByID(ctx *hazart.Context, id string) (*models.User, error)
	Create(ctx *hazart.Context, u *models.User) (*models.User, error)
	Update(ctx *hazart.Context, id string, u *models.User) (*models.User, error)
	Delete(ctx *hazart.Context, id string) error
}

// UserMemoryRepository implements UserRepository with in-memory storage (dapat diganti GORM/SQL)
type UserMemoryRepository struct {
	items map[string]*models.User
}

func NewUserMemoryRepository() *UserMemoryRepository {
	return &UserMemoryRepository{
		items: make(map[string]*models.User),
	}
}

func (r *UserMemoryRepository) FindAll(ctx *hazart.Context) ([]models.User, error) {
	users := make([]models.User, 0, len(r.items))
	for _, u := range r.items {
		users = append(users, *u)
	}
	return users, nil
}

func (r *UserMemoryRepository) FindByID(ctx *hazart.Context, id string) (*models.User, error) {
	u, exists := r.items[id]
	if !exists {
		return nil, hazart.NotFound("User tidak ditemukan")
	}
	return u, nil
}

func (r *UserMemoryRepository) Create(ctx *hazart.Context, u *models.User) (*models.User, error) {
	if u.ID == "" {
		u.ID = hazart.GenerateID()
	}
	r.items[u.ID] = u
	return u, nil
}

func (r *UserMemoryRepository) Update(ctx *hazart.Context, id string, u *models.User) (*models.User, error) {
	if _, exists := r.items[id]; !exists {
		return nil, hazart.NotFound("User tidak ditemukan untuk di-update")
	}
	u.ID = id
	r.items[id] = u
	return u, nil
}

func (r *UserMemoryRepository) Delete(ctx *hazart.Context, id string) error {
	if _, exists := r.items[id]; !exists {
		return hazart.NotFound("User tidak ditemukan untuk dihapus")
	}
	delete(r.items, id)
	return nil
}
