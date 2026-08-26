package crud

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/misbakhul29/hazartgo"
)

// GormRepository is a generic repository implementation backed by GORM with pagination, sorting & search.
type GormRepository[T any] struct {
	db any // interface for *gorm.DB to keep soft dependency
}

// GormDBInterface defines minimal required GORM DB methods
type GormDBInterface interface {
	Find(dest interface{}, conds ...interface{}) GormDBInterface
	First(dest interface{}, conds ...interface{}) GormDBInterface
	Create(value interface{}) GormDBInterface
	Save(value interface{}) GormDBInterface
	Delete(value interface{}, conds ...interface{}) GormDBInterface
	Where(query interface{}, args ...interface{}) GormDBInterface
	Order(value interface{}) GormDBInterface
	Limit(limit int) GormDBInterface
	Offset(offset int) GormDBInterface
	Error() error
}

func NewGormRepository[T any](db GormDBInterface) *GormRepository[T] {
	return &GormRepository[T]{db: db}
}

func (g *GormRepository[T]) FindAll(ctx *hazart.Context) ([]T, error) {
	var items []T
	db, ok := g.db.(GormDBInterface)
	if !ok {
		return nil, fmt.Errorf("invalid gorm db instance")
	}

	query := db

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

	// Search (Parameterized query + LIKE wildcard escaping)
	searchKey := ctx.Query("search")
	if searchKey != "" {
		escapedPattern := "%" + sanitizeLikePattern(searchKey) + "%"
		query = query.Where("name LIKE ? OR title LIKE ?", escapedPattern, escapedPattern)
	}

	if err := query.Find(&items).Error(); err != nil {
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
	db, ok := g.db.(GormDBInterface)
	if !ok {
		return nil, fmt.Errorf("invalid gorm db instance")
	}

	if err := db.First(&item, "id = ?", id).Error(); err != nil {
		return nil, hazart.NotFound("Resource with specified ID not found")
	}
	return &item, nil
}

func (g *GormRepository[T]) Create(ctx *hazart.Context, entity *T) (*T, error) {
	db, ok := g.db.(GormDBInterface)
	if !ok {
		return nil, fmt.Errorf("invalid gorm db instance")
	}

	if err := db.Create(entity).Error(); err != nil {
		return nil, err
	}
	return entity, nil
}

func (g *GormRepository[T]) Update(ctx *hazart.Context, id string, entity *T) (*T, error) {
	db, ok := g.db.(GormDBInterface)
	if !ok {
		return nil, fmt.Errorf("invalid gorm db instance")
	}

	// Set ID field if setting string ID
	val := reflect.ValueOf(entity).Elem()
	idField := val.FieldByName("ID")
	if idField.IsValid() && idField.CanSet() && idField.Kind() == reflect.String {
		idField.SetString(id)
	}

	if err := db.Save(entity).Error(); err != nil {
		return nil, err
	}
	return entity, nil
}

func (g *GormRepository[T]) Delete(ctx *hazart.Context, id string) error {
	var entity T
	db, ok := g.db.(GormDBInterface)
	if !ok {
		return fmt.Errorf("invalid gorm db instance")
	}

	if err := db.Delete(&entity, "id = ?", id).Error(); err != nil {
		return err
	}
	return nil
}

// AutoCRUDGorm registers AutoCRUD endpoints backed by a GORM database interface
func AutoCRUDGorm[T any](g *hazart.Group, path string, db GormDBInterface) {
	repo := NewGormRepository[T](db)
	AutoCRUD[T](g, path, repo)
}
