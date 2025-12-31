package svc

import (
	"mime"
	"path/filepath"
	"strings"
)

// Tools provides utility functions for file validation and processing
type Tools struct {
	MaxFileSize int64
	AllowedExts map[string]struct{}
}

func NewTools(maxSize int64, exts []string) *Tools {
	t := &Tools{
		MaxFileSize: maxSize,
		AllowedExts: make(map[string]struct{}),
	}
	for _, ext := range exts {
		t.AllowedExts[strings.ToLower(ext)] = struct{}{}
	}
	return t
}

func (t *Tools) CheckFileSize(size int64) bool {
	if t.MaxFileSize <= 0 {
		return true
	}
	return size <= t.MaxFileSize
}

func (t *Tools) CheckExtension(filename string) bool {
	if len(t.AllowedExts) == 0 {
		return true
	}
	ext := strings.ToLower(filepath.Ext(filename))
	// remove dot if present
	if strings.HasPrefix(ext, ".") {
		ext = ext[1:]
	}
	_, ok := t.AllowedExts[ext]
	return ok
}

func (t *Tools) GetContentType(filename string) string {
	ext := filepath.Ext(filename)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		return "application/octet-stream"
	}
	return mimeType
}
