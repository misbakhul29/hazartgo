package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]

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

	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Println(`⚡ HazartGo CLI Scaffolding Tool

Usage:
  hazart init <project_name>                Initialize a new HazartGo backend project
  hazart make:controller <Name> [--group]   Generate a new Controller struct & routes
  hazart make:resource <Name> [--group]     Generate complete Model + AutoCRUD Controller

Flags:
  --group <prefix>    Specify route group prefix (default: /api/v1)`)
}

func parseFlag(flagName string, defaultValue string) string {
	for i, arg := range os.Args {
		if arg == flagName && i+1 < len(os.Args) {
			return os.Args[i+1]
		}
	}
	return defaultValue
}

func initProject(projectName string) {
	fmt.Printf("🚀 Initializing new HazartGo project '%s'...\n", projectName)

	os.MkdirAll(filepath.Join(projectName, "controllers"), 0755)
	os.MkdirAll(filepath.Join(projectName, "models"), 0755)

	mainContent := fmt.Sprintf(`package main

import (
	"log"

	hazart "github.com/misbakhul29/hazartgo"
	"github.com/misbakhul29/hazartgo/middleware"
)

func main() {
	app := hazart.New(hazart.Config{
		Title:   "%s API",
		Version: "1.0.0",
	})

	app.Use(middleware.Logger())
	app.Use(middleware.Recovery())
	app.Use(middleware.CORS())

	log.Println("⚡ HazartGo Server running on http://localhost:8080")
	log.Println("📚 Swagger UI Docs available at http://localhost:8080/docs")

	app.Listen(":8080")
}
`, projectName)

	os.WriteFile(filepath.Join(projectName, "main.go"), []byte(mainContent), 0644)
	fmt.Printf("✅ Project '%s' successfully initialized!\n", projectName)
}

func capitalize(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func makeController(name string, group string) {
	name = capitalize(name)
	filename := fmt.Sprintf("controllers/%s_controller.go", strings.ToLower(name))

	content := fmt.Sprintf(`package controllers

import (
	hazart "github.com/misbakhul29/hazartgo"
	"github.com/misbakhul29/hazartgo/middleware"
)

type %sController struct{}

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
`, name, name, name, name, name, name)

	os.MkdirAll("controllers", 0755)
	os.WriteFile(filename, []byte(content), 0644)
	fmt.Printf("✅ Controller created at '%s' under group '%s'!\n", filename, group)
}

func makeResource(name string, group string) {
	name = capitalize(name)
	modelFile := fmt.Sprintf("models/%s.go", strings.ToLower(name))
	controllerFile := fmt.Sprintf("controllers/%s_controller.go", strings.ToLower(name))

	modelContent := fmt.Sprintf(`package models

type %s struct {
	ID        string ` + "`json:\"id\" path:\"id\" doc:\"%s ID\"`" + `
	Name      string ` + "`json:\"name\" validate:\"required\" doc:\"%s Name\"`" + `
}
`, name, name, name)

	controllerContent := fmt.Sprintf(`package controllers

import (
	hazart "github.com/misbakhul29/hazartgo"
	"github.com/misbakhul29/hazartgo/crud"
	"github.com/misbakhul29/hazartgo/middleware"
	"%s/models"
)

type %sController struct{}

func (c *%sController) RegisterRoutes(g *hazart.Group) {
	// Optional RBAC Middleware:
	// g.Use(middleware.RequireRole("admin"))

	// Auto CRUD endpoints generated for %s Model!
	crud.AutoCRUD[models.%s](g, "/%ss")
}
`, getModuleName(), name, name, name, name, strings.ToLower(name))

	os.MkdirAll("models", 0755)
	os.MkdirAll("controllers", 0755)

	os.WriteFile(modelFile, []byte(modelContent), 0644)
	os.WriteFile(controllerFile, []byte(controllerContent), 0644)

	fmt.Printf("✅ Model created at '%s'\n", modelFile)
	fmt.Printf("✅ AutoCRUD Controller created at '%s'\n", controllerFile)
}

func getModuleName() string {
	data, err := os.ReadFile("go.mod")
	if err != nil {
		return "app"
	}
	lines := strings.Split(string(data), "\n")
	if len(lines) > 0 && strings.HasPrefix(lines[0], "module ") {
		return strings.TrimSpace(strings.TrimPrefix(lines[0], "module "))
	}
	return "app"
}
