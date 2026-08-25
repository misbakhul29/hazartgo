package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const cliVersion = "v1.2.0"

func main() {
	if len(os.Args) < 2 {
		checkLatestVersion()
		printUsage()
		return
	}

	command := os.Args[1]

	if command == "update" {
		updateHazart()
		return
	}

	checkLatestVersion()

	switch command {
	case "init":
		if len(os.Args) < 3 {
			log.Fatal("Usage: hazart init <project_name>")
		}
		initProject(os.Args[2])

	case "make:controller", "controller":
		if len(os.Args) < 3 {
			log.Fatal("Usage: hazart make:controller <Name> [--group <prefix>]")
		}
		name := os.Args[2]
		group := parseFlag("--group", "/api/v1")
		makeController(name, group)

	case "make:resource", "resource":
		if len(os.Args) < 3 {
			log.Fatal("Usage: hazart make:resource <Name> [--group <prefix>]")
		}
		name := os.Args[2]
		group := parseFlag("--group", "/api/v1")
		makeResource(name, group)

	case "make:model", "model":
		if len(os.Args) < 3 {
			log.Fatal("Usage: hazart make:model <Name>")
		}
		makeModel(os.Args[2])

	case "make:repository", "repository":
		if len(os.Args) < 3 {
			log.Fatal("Usage: hazart make:repository <Name>")
		}
		makeRepository(os.Args[2])

	case "make:middleware", "middleware":
		if len(os.Args) < 3 {
			log.Fatal("Usage: hazart make:middleware <Name>")
		}
		makeMiddleware(os.Args[2])

	case "make:auth", "auth":
		makeAuth()

	case "version", "-v", "--version":
		fmt.Printf("⚡ HazartGo CLI Scaffolding Tool %s\n", cliVersion)

	case "help", "-h", "--help":
		printUsage()

	default:
		printUsage()
	}
}

func checkLatestVersion() {
	latest := fetchLatestRemoteVersion()
	if isNewerVersion(latest, cliVersion) {
		fmt.Printf("💡 Notice: A new HazartGo release (%s) is available (Current installed: %s).\n", latest, cliVersion)
		fmt.Println("👉 Run 'hazart update' to upgrade HazartGo CLI & project library.")
		fmt.Println()
	}
}

func isNewerVersion(latest, current string) bool {
	if latest == "" || strings.HasPrefix(latest, "v0.0.0-") {
		return false
	}
	latestClean := strings.TrimPrefix(latest, "v")
	currentClean := strings.TrimPrefix(current, "v")

	lParts := strings.Split(latestClean, ".")
	cParts := strings.Split(currentClean, ".")

	if len(lParts) < 3 || len(cParts) < 3 {
		return false
	}

	for i := 0; i < 3; i++ {
		var lNum, cNum int
		fmt.Sscanf(lParts[i], "%d", &lNum)
		fmt.Sscanf(cParts[i], "%d", &cNum)
		if lNum > cNum {
			return true
		} else if lNum < cNum {
			return false
		}
	}
	return false
}

func fetchLatestRemoteVersion() string {
	client := http.Client{
		Timeout: 1500 * time.Millisecond,
	}
	resp, err := client.Get("https://proxy.golang.org/github.com/misbakhul29/hazartgo/@latest")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ""
	}

	var data struct {
		Version string `json:"Version"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return ""
	}

	return data.Version
}

func updateHazart() {
	fmt.Println("⚡ HazartGo Auto-Updater")
	fmt.Printf("🔍 Checking for updates... (Current version: %s)\n", cliVersion)

	latest := fetchLatestRemoteVersion()
	if latest != "" && latest != cliVersion {
		fmt.Printf("🎉 New version available: %s (Installed: %s)\n\n", latest, cliVersion)
	} else if latest != "" {
		fmt.Printf("✅ You are already using the latest CLI version (%s)!\n\n", cliVersion)
	}

	fmt.Println("📦 Upgrading HazartGo CLI tool (go install github.com/misbakhul29/hazartgo/cmd/hazart@latest)...")
	cmdCLI := exec.Command("go", "install", "github.com/misbakhul29/hazartgo/cmd/hazart@latest")
	cmdCLI.Stdout = os.Stdout
	cmdCLI.Stderr = os.Stderr
	if err := cmdCLI.Run(); err != nil {
		fmt.Printf("❌ Failed to upgrade HazartGo CLI tool: %v\n", err)
	} else {
		fmt.Println("✅ HazartGo CLI tool successfully upgraded to latest version!")
	}

	if _, err := os.Stat("go.mod"); err == nil {
		fmt.Println("\n📦 Upgrading hazartgo library in current project (go get -u github.com/misbakhul29/hazartgo@latest)...")
		cmdLib := exec.Command("go", "get", "-u", "github.com/misbakhul29/hazartgo@latest")
		cmdLib.Stdout = os.Stdout
		cmdLib.Stderr = os.Stderr
		if err := cmdLib.Run(); err == nil {
			exec.Command("go", "mod", "tidy").Run()
			fmt.Println("✅ HazartGo library in current project successfully upgraded!")
		} else {
			fmt.Printf("❌ Failed to upgrade hazartgo library in current project: %v\n", err)
		}
	}
}

func printUsage() {
	fmt.Printf(`⚡ HazartGo CLI Scaffolding Tool (%s)

Usage:
  hazart <command> [arguments] [flags]

Commands:
  init <project_name>                Initialize a new HazartGo backend project
  make:controller <Name> [--group]   Generate a Controller struct & route handlers
  make:resource <Name> [--group]     Generate complete Model + Repository + AutoCRUD Controller
  make:model <Name>                  Generate a new DB/OpenAPI Model struct
  make:repository <Name>             Generate a new Repository interface & memory store
  make:middleware <Name>             Generate a new custom Middleware handler
  make:auth                          Scaffold JWT Authentication Controller
  update                             Automatically upgrade HazartGo CLI & project library

Aliases:
  controller, resource, model, repository, middleware, auth

Flags:
  --group <prefix>    Specify route group prefix (default: /api/v1)
  --version, -v       Display HazartGo CLI version
  --help, -h          Display CLI help documentation
`, cliVersion)
}

func parseFlag(flagName string, defaultValue string) string {
	for i, arg := range os.Args {
		if arg == flagName && i+1 < len(os.Args) {
			return os.Args[i+1]
		}
	}
	return defaultValue
}

func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

func toPascalCase(s string) string {
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-' || r == ' '
	})
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	if len(parts) == 0 {
		return capitalize(s)
	}
	return strings.Join(parts, "")
}

func capitalize(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func initProject(projectName string) {
	fmt.Printf("🚀 Initializing new HazartGo project '%s'...\n", projectName)

	targetDir := projectName
	moduleName := projectName
	if projectName == "." {
		cwd, err := os.Getwd()
		if err == nil {
			moduleName = filepath.Base(cwd)
		} else {
			moduleName = "app"
		}
	} else {
		moduleName = filepath.Base(projectName)
		os.MkdirAll(targetDir, 0755)
	}

	os.MkdirAll(filepath.Join(targetDir, "controllers"), 0755)
	os.MkdirAll(filepath.Join(targetDir, "models"), 0755)
	os.MkdirAll(filepath.Join(targetDir, "repositories"), 0755)
	os.MkdirAll(filepath.Join(targetDir, "middleware"), 0755)

	goModPath := filepath.Join(targetDir, "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		goModContent := fmt.Sprintf("module %s\n\ngo 1.22\n", moduleName)
		os.WriteFile(goModPath, []byte(goModContent), 0644)
	}

	itemModelContent := `package models

type Item struct {
	ID          string  ` + "`json:\"id\" path:\"id\" doc:\"Item ID\"`" + `
	Name        string  ` + "`json:\"name\" validate:\"required\" doc:\"Item Name\"`" + `
	Price       float64 ` + "`json:\"price\" validate:\"required\" doc:\"Item Price\"`" + `
	Description string  ` + "`json:\"description\" doc:\"Item Description\"`" + `
}
`
	os.WriteFile(filepath.Join(targetDir, "models", "item.go"), []byte(itemModelContent), 0644)

	itemRepoContent := fmt.Sprintf(`package repositories

import (
	"sync"

	hazart "github.com/misbakhul29/hazartgo"
	"%s/models"
)

type ItemRepository interface {
	FindAll(ctx *hazart.Context) ([]models.Item, error)
	FindByID(ctx *hazart.Context, id string) (*models.Item, error)
	Create(ctx *hazart.Context, entity *models.Item) (*models.Item, error)
	Update(ctx *hazart.Context, id string, entity *models.Item) (*models.Item, error)
	Delete(ctx *hazart.Context, id string) error
}

type itemMemoryRepository struct {
	mu    sync.RWMutex
	items map[string]*models.Item
}

func NewItemMemoryRepository() ItemRepository {
	repo := &itemMemoryRepository{
		items: make(map[string]*models.Item),
	}
	repo.items["item_1"] = &models.Item{ID: "item_1", Name: "HazartGo Pro Keyboard", Price: 149.99, Description: "Mechanical RGB Keyboard"}
	repo.items["item_2"] = &models.Item{ID: "item_2", Name: "HazartGo Wireless Mouse", Price: 79.50, Description: "Ergonomic Optical Mouse"}
	return repo
}

func (r *itemMemoryRepository) FindAll(ctx *hazart.Context) ([]models.Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]models.Item, 0, len(r.items))
	for _, item := range r.items {
		result = append(result, *item)
	}
	return result, nil
}

func (r *itemMemoryRepository) FindByID(ctx *hazart.Context, id string) (*models.Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, exists := r.items[id]
	if !exists {
		return nil, hazart.NotFound("Item not found")
	}
	return item, nil
}

func (r *itemMemoryRepository) Create(ctx *hazart.Context, entity *models.Item) (*models.Item, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if entity.ID == "" {
		entity.ID = hazart.GenerateID()
	}
	r.items[entity.ID] = entity
	return entity, nil
}

func (r *itemMemoryRepository) Update(ctx *hazart.Context, id string, entity *models.Item) (*models.Item, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.items[id]; !exists {
		return nil, hazart.NotFound("Item not found")
	}
	entity.ID = id
	r.items[id] = entity
	return entity, nil
}

func (r *itemMemoryRepository) Delete(ctx *hazart.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.items[id]; !exists {
		return hazart.NotFound("Item not found")
	}
	delete(r.items, id)
	return nil
}
`, moduleName)
	os.WriteFile(filepath.Join(targetDir, "repositories", "item_repository.go"), []byte(itemRepoContent), 0644)

	itemControllerContent := fmt.Sprintf(`package controllers

import (
	hazart "github.com/misbakhul29/hazartgo"
	"github.com/misbakhul29/hazartgo/crud"
	"%s/models"
	"%s/repositories"
)

type ItemController struct {
	repo repositories.ItemRepository
}

func NewItemController(repo repositories.ItemRepository) *ItemController {
	return &ItemController{repo: repo}
}

func (c *ItemController) RegisterRoutes(g *hazart.Group) {
	// Auto CRUD endpoints generated for Item Model!
	// Endpoints: GET /items, GET /items/:id, POST /items, PUT /items/:id, DELETE /items/:id
	crud.AutoCRUD[models.Item](g, "", c.repo)
}
`, moduleName, moduleName)
	os.WriteFile(filepath.Join(targetDir, "controllers", "item_controller.go"), []byte(itemControllerContent), 0644)

	authMiddlewareContent := `package middleware

import (
	hazart "github.com/misbakhul29/hazartgo"
)

// AuthGuard middleware example injecting context user state
func AuthGuard() hazart.MiddlewareFunc {
	return func(next hazart.HandlerFunc) hazart.HandlerFunc {
		return func(ctx *hazart.Context) {
			ctx.Set("user_id", "user_123")
			ctx.Set("user_roles", []string{"admin"})
			next(ctx)
		}
	}
}
`
	os.WriteFile(filepath.Join(targetDir, "middleware", "auth.go"), []byte(authMiddlewareContent), 0644)

	mainContent := fmt.Sprintf(`package main

import (
	"log"

	hazart "github.com/misbakhul29/hazartgo"
	"github.com/misbakhul29/hazartgo/middleware"
	"github.com/misbakhul29/hazartgo/openapi"

	"%s/controllers"
	"%s/repositories"
)

func main() {
	// 1. Initialize HazartGo Framework App
	app := hazart.New(hazart.Config{
		Title:       "%s API",
		Description: "High Performance & Developer-Friendly Go REST API Framework",
		Version:     "1.0.0",
		Contact: &openapi.Contact{
			Name:  "API Support",
			Email: "support@example.com",
		},
		License: &openapi.License{
			Name: "MIT",
		},
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

	// 4. Initialize Dependencies & Controllers
	itemRepo := repositories.NewItemMemoryRepository()
	itemController := controllers.NewItemController(itemRepo)

	// 5. Mount Controllers under API Route Group Prefix
	app.MountController("/api/v1/items", itemController)

	log.Println("⚡ HazartGo Server running on http://localhost:8080")
	log.Println("📚 Swagger UI Docs available at http://localhost:8080/docs")

	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("Server error: %%v", err)
	}
}
`, moduleName, moduleName, moduleName)

	os.WriteFile(filepath.Join(targetDir, "main.go"), []byte(mainContent), 0644)

	fmt.Println()
	fmt.Printf("✅ Project '%s' successfully initialized!\n", projectName)
	fmt.Println("👉 Next steps:")
	if projectName != "." {
		fmt.Printf("   cd %s\n", projectName)
	}
	fmt.Println("   go mod tidy")
	fmt.Println("   go run main.go")
}

func makeController(name string, group string) {
	pascalName := toPascalCase(name)
	snakeName := toSnakeCase(pascalName)
	filename := fmt.Sprintf("controllers/%s_controller.go", snakeName)

	content := fmt.Sprintf(`package controllers

import (
	hazart "github.com/misbakhul29/hazartgo"
	"github.com/misbakhul29/hazartgo/middleware"
)

type %sController struct{}

func New%sController() *%sController {
	return &%sController{}
}

func (c *%sController) RegisterRoutes(g *hazart.Group) {
	// Example Role Authorization Middleware:
	// g.Use(middleware.RequireRole("admin"))

	hazart.GroupGet(g, "/ping", c.Ping, hazart.RouteMeta{
		Summary: "Ping %s",
		Tags:    []string{"%s"},
	})
}

type PingRes struct {
	Message string ` + "`json:\"message\" doc:\"Ping response\"`" + `
}

func (c *%sController) Ping(ctx *hazart.Context, req *struct{}) (*PingRes, error) {
	return &PingRes{Message: "pong from %sController"}, nil
}
`, pascalName, pascalName, pascalName, pascalName, pascalName, pascalName, pascalName, pascalName, pascalName)

	os.MkdirAll("controllers", 0755)
	os.WriteFile(filename, []byte(content), 0644)
	fmt.Printf("✅ Controller created at '%s'!\n", filename)
	fmt.Printf("📌 Mount it in main.go with: app.MountController(\"%s\", controllers.New%sController())\n", group, pascalName)
}

func makeResource(name string, group string) {
	pascalName := toPascalCase(name)
	snakeName := toSnakeCase(pascalName)
	modName := getModuleName()

	modelFile := fmt.Sprintf("models/%s.go", snakeName)
	repoFile := fmt.Sprintf("repositories/%s_repository.go", snakeName)
	controllerFile := fmt.Sprintf("controllers/%s_controller.go", snakeName)

	modelContent := fmt.Sprintf(`package models

type %s struct {
	ID        string ` + "`json:\"id\" path:\"id\" doc:\"%s ID\"`" + `
	Name      string ` + "`json:\"name\" validate:\"required\" doc:\"%s Name\"`" + `
}
`, pascalName, pascalName, pascalName)

	lowerFirst := strings.ToLower(pascalName[:1]) + pascalName[1:]

	repoContent := fmt.Sprintf(`package repositories

import (
	"sync"

	hazart "github.com/misbakhul29/hazartgo"
	"%s/models"
)

type %sRepository interface {
	FindAll(ctx *hazart.Context) ([]models.%s, error)
	FindByID(ctx *hazart.Context, id string) (*models.%s, error)
	Create(ctx *hazart.Context, entity *models.%s) (*models.%s, error)
	Update(ctx *hazart.Context, id string, entity *models.%s) (*models.%s, error)
	Delete(ctx *hazart.Context, id string) error
}

type %sMemoryRepository struct {
	mu    sync.RWMutex
	items map[string]*models.%s
}

func New%sMemoryRepository() %sRepository {
	return &%sMemoryRepository{
		items: make(map[string]*models.%s),
	}
}

func (r *%sMemoryRepository) FindAll(ctx *hazart.Context) ([]models.%s, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	res := make([]models.%s, 0, len(r.items))
	for _, item := range r.items {
		res = append(res, *item)
	}
	return res, nil
}

func (r *%sMemoryRepository) FindByID(ctx *hazart.Context, id string) (*models.%s, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, exists := r.items[id]
	if !exists {
		return nil, hazart.NotFound("%s not found")
	}
	return item, nil
}

func (r *%sMemoryRepository) Create(ctx *hazart.Context, entity *models.%s) (*models.%s, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if entity.ID == "" {
		entity.ID = hazart.GenerateID()
	}
	r.items[entity.ID] = entity
	return entity, nil
}

func (r *%sMemoryRepository) Update(ctx *hazart.Context, id string, entity *models.%s) (*models.%s, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.items[id]; !exists {
		return nil, hazart.NotFound("%s not found")
	}
	entity.ID = id
	r.items[id] = entity
	return entity, nil
}

func (r *%sMemoryRepository) Delete(ctx *hazart.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.items[id]; !exists {
		return hazart.NotFound("%s not found")
	}
	delete(r.items, id)
	return nil
}
`, modName, pascalName, pascalName, pascalName, pascalName, pascalName, pascalName, pascalName,
		lowerFirst, pascalName, pascalName, pascalName, lowerFirst, pascalName,
		lowerFirst, pascalName, pascalName,
		lowerFirst, pascalName, pascalName,
		lowerFirst, pascalName, pascalName,
		lowerFirst, pascalName, pascalName, pascalName,
		lowerFirst, pascalName)

	controllerContent := fmt.Sprintf(`package controllers

import (
	hazart "github.com/misbakhul29/hazartgo"
	"github.com/misbakhul29/hazartgo/crud"
	"github.com/misbakhul29/hazartgo/middleware"
	"%s/models"
	"%s/repositories"
)

type %sController struct {
	repo repositories.%sRepository
}

func New%sController(repo repositories.%sRepository) *%sController {
	return &%sController{repo: repo}
}

func (c *%sController) RegisterRoutes(g *hazart.Group) {
	// Optional RBAC Middleware:
	// g.Use(middleware.RequireRole("admin"))

	// Auto CRUD endpoints generated for %s Model!
	crud.AutoCRUD[models.%s](g, "", c.repo)
}
`, modName, modName, pascalName, pascalName, pascalName, pascalName, pascalName, pascalName, pascalName, pascalName, pascalName)

	os.MkdirAll("models", 0755)
	os.MkdirAll("repositories", 0755)
	os.MkdirAll("controllers", 0755)

	os.WriteFile(modelFile, []byte(modelContent), 0644)
	os.WriteFile(repoFile, []byte(repoContent), 0644)
	os.WriteFile(controllerFile, []byte(controllerContent), 0644)

	fmt.Printf("✅ Model created at '%s'\n", modelFile)
	fmt.Printf("✅ Repository created at '%s'\n", repoFile)
	fmt.Printf("✅ AutoCRUD Controller created at '%s'\n", controllerFile)
	fmt.Printf("📌 Mount it in main.go with: app.MountController(\"%s/%ss\", controllers.New%sController(repositories.New%sMemoryRepository()))\n", group, strings.ToLower(pascalName), pascalName, pascalName)
}

func makeModel(name string) {
	pascalName := toPascalCase(name)
	snakeName := toSnakeCase(pascalName)
	filename := fmt.Sprintf("models/%s.go", snakeName)

	content := fmt.Sprintf(`package models

type %s struct {
	ID        string ` + "`json:\"id\" path:\"id\" doc:\"%s ID\"`" + `
	Name      string ` + "`json:\"name\" validate:\"required\" doc:\"%s Name\"`" + `
}
`, pascalName, pascalName, pascalName)

	os.MkdirAll("models", 0755)
	os.WriteFile(filename, []byte(content), 0644)
	fmt.Printf("✅ Model created at '%s'\n", filename)
}

func makeRepository(name string) {
	pascalName := toPascalCase(name)
	snakeName := toSnakeCase(pascalName)
	modName := getModuleName()
	filename := fmt.Sprintf("repositories/%s_repository.go", snakeName)
	lowerFirst := strings.ToLower(pascalName[:1]) + pascalName[1:]

	content := fmt.Sprintf(`package repositories

import (
	"sync"

	hazart "github.com/misbakhul29/hazartgo"
	"%s/models"
)

type %sRepository interface {
	FindAll(ctx *hazart.Context) ([]models.%s, error)
	FindByID(ctx *hazart.Context, id string) (*models.%s, error)
	Create(ctx *hazart.Context, entity *models.%s) (*models.%s, error)
	Update(ctx *hazart.Context, id string, entity *models.%s) (*models.%s, error)
	Delete(ctx *hazart.Context, id string) error
}

type %sMemoryRepository struct {
	mu    sync.RWMutex
	items map[string]*models.%s
}

func New%sMemoryRepository() %sRepository {
	return &%sMemoryRepository{
		items: make(map[string]*models.%s),
	}
}

func (r *%sMemoryRepository) FindAll(ctx *hazart.Context) ([]models.%s, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	res := make([]models.%s, 0, len(r.items))
	for _, item := range r.items {
		res = append(res, *item)
	}
	return res, nil
}

func (r *%sMemoryRepository) FindByID(ctx *hazart.Context, id string) (*models.%s, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, exists := r.items[id]
	if !exists {
		return nil, hazart.NotFound("%s not found")
	}
	return item, nil
}

func (r *%sMemoryRepository) Create(ctx *hazart.Context, entity *models.%s) (*models.%s, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if entity.ID == "" {
		entity.ID = hazart.GenerateID()
	}
	r.items[entity.ID] = entity
	return entity, nil
}

func (r *%sMemoryRepository) Update(ctx *hazart.Context, id string, entity *models.%s) (*models.%s, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.items[id]; !exists {
		return nil, hazart.NotFound("%s not found")
	}
	entity.ID = id
	r.items[id] = entity
	return entity, nil
}

func (r *%sMemoryRepository) Delete(ctx *hazart.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.items[id]; !exists {
		return hazart.NotFound("%s not found")
	}
	delete(r.items, id)
	return nil
}
`, modName, pascalName, pascalName, pascalName, pascalName, pascalName, pascalName, pascalName,
		lowerFirst, pascalName, pascalName, pascalName, lowerFirst, pascalName,
		lowerFirst, pascalName, pascalName,
		lowerFirst, pascalName, pascalName,
		lowerFirst, pascalName, pascalName,
		lowerFirst, pascalName, pascalName, pascalName,
		lowerFirst, pascalName)

	os.MkdirAll("repositories", 0755)
	os.WriteFile(filename, []byte(content), 0644)
	fmt.Printf("✅ Repository created at '%s'\n", filename)
}

func makeMiddleware(name string) {
	pascalName := toPascalCase(name)
	snakeName := toSnakeCase(pascalName)
	filename := fmt.Sprintf("middleware/%s.go", snakeName)

	content := fmt.Sprintf(`package middleware

import (
	hazart "github.com/misbakhul29/hazartgo"
)

// %s custom middleware handler
func %s() hazart.MiddlewareFunc {
	return func(next hazart.HandlerFunc) hazart.HandlerFunc {
		return func(ctx *hazart.Context) {
			// Middleware logic here before handler
			next(ctx)
			// Middleware logic here after handler
		}
	}
}
`, pascalName, pascalName)

	os.MkdirAll("middleware", 0755)
	os.WriteFile(filename, []byte(content), 0644)
	fmt.Printf("✅ Middleware created at '%s'\n", filename)
}

func makeAuth() {
	authControllerFile := "controllers/auth_controller.go"

	content := `package controllers

import (
	"time"

	hazart "github.com/misbakhul29/hazartgo"
	"github.com/misbakhul29/hazartgo/jwt"
)

type AuthController struct {
	jwtManager *jwt.JWT
}

func NewAuthController(secretKey string) *AuthController {
	return &AuthController{
		jwtManager: jwt.New(secretKey),
	}
}

func (c *AuthController) RegisterRoutes(g *hazart.Group) {
	hazart.GroupPost(g, "/login", c.Login, hazart.RouteMeta{
		Summary:     "User Login",
		Description: "Authenticate user and issue JWT Bearer Token",
		Tags:        []string{"Auth"},
	})
}

type LoginReq struct {
	Username string ` + "`json:\"username\" validate:\"required\" doc:\"Username\"`" + `
	Password string ` + "`json:\"password\" validate:\"required\" doc:\"Password\"`" + `
}

type LoginRes struct {
	Token string ` + "`json:\"token\" doc:\"JWT Access Token\"`" + `
}

func (c *AuthController) Login(ctx *hazart.Context, req *LoginReq) (*LoginRes, error) {
	if req.Username != "admin" || req.Password != "password123" {
		return nil, hazart.Unauthorized("Invalid username or password")
	}

	token, err := c.jwtManager.Sign(jwt.MapClaims{
		"sub":   "user_123",
		"name":  req.Username,
		"roles": []string{"admin"},
	}, 24*time.Hour)

	if err != nil {
		return nil, hazart.InternalServerError("Failed to issue token")
	}

	return &LoginRes{Token: token}, nil
}
`

	os.MkdirAll("controllers", 0755)
	os.WriteFile(authControllerFile, []byte(content), 0644)
	fmt.Printf("✅ AuthController created at '%s'!\n", authControllerFile)
	fmt.Println("📌 Mount it in main.go with: app.MountController(\"/api/v1/auth\", controllers.NewAuthController(\"your-secret-key\"))")
}

func getModuleName() string {
	data, err := os.ReadFile("go.mod")
	if err != nil {
		return "app"
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module "))
		}
	}
	return "app"
}

