package crud

import (
	"reflect"

	"github.com/misbakhul29/hazartgo"
)

// Repository is a generic interface for database persistence operations
type Repository[T any] interface {
	FindAll(ctx *hazart.Context) ([]T, error)
	FindByID(ctx *hazart.Context, id string) (*T, error)
	Create(ctx *hazart.Context, entity *T) (*T, error)
	Update(ctx *hazart.Context, id string, entity *T) (*T, error)
	Delete(ctx *hazart.Context, id string) error
}

// MemoryRepository is a default in-memory repository implementation for quick prototyping
type MemoryRepository[T any] struct {
	items map[string]*T
}

func NewMemoryRepository[T any]() *MemoryRepository[T] {
	return &MemoryRepository[T]{
		items: make(map[string]*T),
	}
}

func (m *MemoryRepository[T]) FindAll(ctx *hazart.Context) ([]T, error) {
	res := make([]T, 0, len(m.items))
	for _, item := range m.items {
		res = append(res, *item)
	}
	return res, nil
}

func (m *MemoryRepository[T]) FindByID(ctx *hazart.Context, id string) (*T, error) {
	item, exists := m.items[id]
	if !exists {
		return nil, hazart.NotFound("Resource not found")
	}
	return item, nil
}

func (m *MemoryRepository[T]) Create(ctx *hazart.Context, entity *T) (*T, error) {
	// Simple ID reflection assignment if field ID exists
	val := reflect.ValueOf(entity).Elem()
	idField := val.FieldByName("ID")
	if idField.IsValid() && idField.CanSet() {
		if idField.Kind() == reflect.String && idField.String() == "" {
			idField.SetString(hazart.GenerateID())
		}
	}

	id := hazart.ExtractID(entity)
	m.items[id] = entity
	return entity, nil
}

func (m *MemoryRepository[T]) Update(ctx *hazart.Context, id string, entity *T) (*T, error) {
	if _, exists := m.items[id]; !exists {
		return nil, hazart.NotFound("Resource not found")
	}
	m.items[id] = entity
	return entity, nil
}

func (m *MemoryRepository[T]) Delete(ctx *hazart.Context, id string) error {
	if _, exists := m.items[id]; !exists {
		return hazart.NotFound("Resource not found")
	}
	delete(m.items, id)
	return nil
}

type IDParamReq struct {
	ID string `path:"id" doc:"Resource unique ID"`
}

type EmptyReq struct{}

// AutoCRUD automatically registers 5 standard REST CRUD routes for a given entity model type T
func AutoCRUD[T any](g *hazart.Group, path string, repo ...Repository[T]) {
	var repository Repository[T]
	if len(repo) > 0 {
		repository = repo[0]
	} else {
		repository = NewMemoryRepository[T]()
	}

	entityName := reflect.TypeOf((*T)(nil)).Elem().Name()

	// 1. GET / (List All)
	hazart.GroupGet(g, path, func(ctx *hazart.Context, req *EmptyReq) (*[]T, error) {
		items, err := repository.FindAll(ctx)
		if err != nil {
			return nil, err
		}
		return &items, nil
	}, hazart.RouteMeta{
		Summary:     "List All " + entityName + "s",
		Description: "Fetch array of " + entityName + " items",
		Tags:        []string{entityName},
	})

	// 2. GET /:id (Find by ID)
	hazart.GroupGet(g, path+"/:id", func(ctx *hazart.Context, req *IDParamReq) (*T, error) {
		return repository.FindByID(ctx, req.ID)
	}, hazart.RouteMeta{
		Summary:     "Get " + entityName + " By ID",
		Description: "Fetch single " + entityName + " by path ID",
		Tags:        []string{entityName},
	})

	// 3. POST / (Create)
	hazart.GroupPost(g, path, func(ctx *hazart.Context, req *T) (*T, error) {
		return repository.Create(ctx, req)
	}, hazart.RouteMeta{
		Summary:     "Create " + entityName,
		Description: "Create a new " + entityName + " record",
		Tags:        []string{entityName},
	})

	// 4. PUT /:id (Update)
	hazart.GroupPut(g, path+"/:id", func(ctx *hazart.Context, req *T) (*T, error) {
		id := ctx.Param("id")
		return repository.Update(ctx, id, req)
	}, hazart.RouteMeta{
		Summary:     "Update " + entityName,
		Description: "Update an existing " + entityName + " record by ID",
		Tags:        []string{entityName},
	})

	// 5. DELETE /:id (Delete)
	hazart.GroupDelete(g, path+"/:id", func(ctx *hazart.Context, req *IDParamReq) (*hazart.Map, error) {
		err := repository.Delete(ctx, req.ID)
		if err != nil {
			return nil, err
		}
		return &hazart.Map{"message": entityName + " deleted successfully"}, nil
	}, hazart.RouteMeta{
		Summary:     "Delete " + entityName,
		Description: "Delete " + entityName + " record by ID",
		Tags:        []string{entityName},
	})
}
