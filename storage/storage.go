package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/misbakhul29/hazartgo"
)

// UploadOptions configures file upload rules
type UploadOptions struct {
	MaxSize      int64    // Max file size in bytes
	AllowedTypes []string // E.g. []string{"image/png", "image/jpeg"}
}

// FileResult holds metadata of a saved file
type FileResult struct {
	Filename    string `json:"filename"`
	Path        string `json:"path"`
	Size        int64  `json:"size"`
	ContentType string `json:"content_type"`
}

// StorageEngine defines uniform interface for local / S3 storage drivers
type StorageEngine interface {
	Save(fileHeader *multipart.FileHeader, destDir string) (*FileResult, error)
}

// LocalStorage implements StorageEngine for filesystem storage
type LocalStorage struct{}

func NewLocalStorage() *LocalStorage {
	return &LocalStorage{}
}

func (l *LocalStorage) Save(fileHeader *multipart.FileHeader, destDir string) (*FileResult, error) {
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	src, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	destPath := filepath.Join(destDir, fileHeader.Filename)
	out, err := os.Create(destPath)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	size, err := io.Copy(out, src)
	if err != nil {
		return nil, err
	}

	return &FileResult{
		Filename:    fileHeader.Filename,
		Path:        destPath,
		Size:        size,
		ContentType: fileHeader.Header.Get("Content-Type"),
	}, nil
}

// SaveFile parses, validates and saves uploaded file from context form
func SaveFile(c *hazart.Context, fieldName string, destDir string, opts ...UploadOptions) (*FileResult, error) {
	var opt UploadOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	file, header, err := c.Request.FormFile(fieldName)
	if err != nil {
		return nil, fmt.Errorf("failed to get form file '%s': %w", fieldName, err)
	}
	defer file.Close()

	// 1. Validate Size
	if opt.MaxSize > 0 && header.Size > opt.MaxSize {
		return nil, fmt.Errorf("file size %d bytes exceeds maximum limit of %d bytes", header.Size, opt.MaxSize)
	}

	// 2. Validate Content Type
	if len(opt.AllowedTypes) > 0 {
		contentType := header.Header.Get("Content-Type")
		allowed := false
		for _, t := range opt.AllowedTypes {
			if strings.EqualFold(contentType, t) || strings.HasSuffix(header.Filename, t) {
				allowed = true
				break
			}
		}
		if !allowed {
			return nil, fmt.Errorf("content type '%s' is not in allowed types: %v", contentType, opt.AllowedTypes)
		}
	}

	driver := NewLocalStorage()
	return driver.Save(header, destDir)
}
