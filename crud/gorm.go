package crud

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/misbakhul29/hazartgo"
	"gorm.io/gorm"
)

// GormRepository is a generic repository implementation backed by GORM with pagination, sorting & search.
type GormRepository[T any] struct {
	db *gorm.DB
}

func NewGormRepository[T any](db *gorm.DB) *GormRepository[T] {
	return &GormRepository[T]{db: db}
}

func (g *GormRepository[T]) FindAll(ctx *hazart.Context) ([]T, error) {
	var items []T
	if g.db == nil {
		return nil, fmt.Errorf("invalid gorm db instance")
	}

	query := g.db

	// Pagination
	page, _ := strconv.Atoi(ctx.Query("page"))
	limit, _ := strconv.Atoi(ctx.Query("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit
	query = query.Limit(limit).Offset(offset)

	// Sorting (Sanitize column identifier against SQL injection)
	sortField := sanitizeIdentifier(ctx.Query("sort"))
	sortOrder := strings.ToUpper(ctx.Query("order"))
	if sortField != "" {
		if sortOrder != "DESC" {
			sortOrder = "ASC"
		}
		query = query.Order(fmt.Sprintf("%s %s", sortField, sortOrder))
	}

	// Search (Parameterized query + LIKE wildcard escaping based on struct fields)
	searchKey := ctx.Query("search")
	if searchKey != "" {
		escapedPattern := "%" + sanitizeLikePattern(searchKey) + "%"
		var elem T
		typ := reflect.TypeOf(elem)
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}

		hasName := false
		hasTitle := false
		if typ.Kind() == reflect.Struct {
			_, hasName = typ.FieldByName("Name")
			_, hasTitle = typ.FieldByName("Title")
		}

		if hasName && hasTitle {
			query = query.Where("name LIKE ? OR title LIKE ?", escapedPattern, escapedPattern)
		} else if hasName {
			query = query.Where("name LIKE ?", escapedPattern)
		} else if hasTitle {
			query = query.Where("title LIKE ?", escapedPattern)
		}
	}

	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func sanitizeIdentifier(s string) string {
	var sb strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

func sanitizeLikePattern(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "%", "\\%")
	s = strings.ReplaceAll(s, "_", "\\_")
	return s
}

func (g *GormRepository[T]) FindByID(ctx *hazart.Context, id string) (*T, error) {
	var item T
	if g.db == nil {
		return nil, fmt.Errorf("invalid gorm db instance")
	}

	if err := g.db.First(&item, "id = ?", id).Error; err != nil {
		return nil, hazart.NotFound("Resource with specified ID not found")
	}
	return &item, nil
}

func (g *GormRepository[T]) Create(ctx *hazart.Context, entity *T) (*T, error) {
	if g.db == nil {
		return nil, fmt.Errorf("invalid gorm db instance")
	}

	if err := g.db.Create(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (g *GormRepository[T]) Update(ctx *hazart.Context, id string, entity *T) (*T, error) {
	if g.db == nil {
		return nil, fmt.Errorf("invalid gorm db instance")
	}

	// Set ID field if setting string ID
	val := reflect.ValueOf(entity).Elem()
	idField := val.FieldByName("ID")
	if idField.IsValid() && idField.CanSet() && idField.Kind() == reflect.String {
		idField.SetString(id)
	}

	if err := g.db.Save(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (g *GormRepository[T]) Delete(ctx *hazart.Context, id string) error {
	var entity T
	if g.db == nil {
		return fmt.Errorf("invalid gorm db instance")
	}

	if err := g.db.Delete(&entity, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

// AutoCRUDGorm registers AutoCRUD endpoints backed by a GORM database instance
func AutoCRUDGorm[T any](g *hazart.Group, path string, db *gorm.DB) {
	repo := NewGormRepository[T](db)
	AutoCRUD[T](g, path, repo)
}
