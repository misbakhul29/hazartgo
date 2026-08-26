package controllers

import (
	"fmt"
	"time"

	hazart "github.com/misbakhul29/hazartgo"
	"github.com/misbakhul29/hazartgo/crud"
	"github.com/misbakhul29/hazartgo/storage"
)

type FeatureController struct{}

func NewFeatureController() *FeatureController {
	return &FeatureController{}
}

func (fc *FeatureController) RegisterRoutes(g *hazart.Group) {
	// 1. Auto-Validation & Binding Example
	hazart.GroupPost(g, "/validate", fc.TestValidation, hazart.RouteMeta{
		Summary:     "Test Request Validation & Binding",
		Description: "Demonstrates hazart.BindAndValidate with struct tags",
		Tags:        []string{"Features"},
	})

	// 2. Realtime SSE Event Stream Endpoint
	hazart.GroupGet(g, "/sse", fc.TestSSE, hazart.RouteMeta{
		Summary:     "Test Server-Sent Events (SSE)",
		Description: "Streams live server events to connected client",
		Tags:        []string{"Features"},
	})

	// 3. Realtime WebSocket Endpoint
	hazart.GroupGet(g, "/ws", fc.TestWS, hazart.RouteMeta{
		Summary:     "Test WebSocket Connection",
		Description: "Hijacks HTTP connection to WebSocket stream",
		Tags:        []string{"Features"},
	})

	// 4. File Upload Endpoint
	hazart.GroupPost(g, "/upload", fc.TestUpload, hazart.RouteMeta{
		Summary:     "Test File Upload & Storage Driver",
		Description: "Uploads a file with size and content-type validation",
		Tags:        []string{"Features"},
	})
}

type ValidationTestReq struct {
	Name  string `json:"name" validate:"required,min=3" doc:"User Name"`
	Email string `json:"email" validate:"required,email" doc:"User Email"`
	Age   int    `json:"age" validate:"required,gte=18" doc:"User Age (min 18)"`
}

func (fc *FeatureController) TestValidation(ctx *hazart.Context, req *ValidationTestReq) (*ValidationTestReq, error) {
	// Auto validation executed by hazart.BindAndValidate
	return req, nil
}

func (fc *FeatureController) TestSSE(ctx *hazart.Context, req *crud.EmptyReq) (*hazart.Map, error) {
	_ = ctx.SSEvent("welcome", map[string]string{
		"message": "Welcome to HazartGo SSE Stream!",
		"time":    time.Now().Format(time.RFC3339),
	})
	return nil, nil
}

func (fc *FeatureController) TestWS(ctx *hazart.Context, req *crud.EmptyReq) (*hazart.Map, error) {
	return &hazart.Map{"message": "Connect via WebSocket client to /api/v1/features/ws"}, nil
}

func (fc *FeatureController) TestUpload(ctx *hazart.Context, req *crud.EmptyReq) (*storage.FileResult, error) {
	result, err := storage.SaveFile(ctx, "file", "./uploads", storage.UploadOptions{
		MaxSize:      5 * 1024 * 1024, // 5MB
		AllowedTypes: []string{"image/png", "image/jpeg", "text/plain"},
	})
	if err != nil {
		return nil, hazart.BadRequest(fmt.Sprintf("Upload failed: %v", err))
	}
	return result, nil
}
