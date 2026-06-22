package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/grout/config"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "grout.yaml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writing temp config: %v", err)
	}
	return path
}

func TestLoad_Valid(t *testing.T) {
	path := writeTemp(t, `
version: "1"
services:
  - name: api
    command: go run ./cmd/api
    dir: ./api
    port: 8080
  - name: worker
    command: go run ./cmd/worker
`)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Services) != 2 {
		t.Errorf("expected 2 services, got %d", len(cfg.Services))
	}
	if cfg.Services[0].Port != 8080 {
		t.Errorf("expected port 8080, got %d", cfg.Services[0].Port)
	}
}

func TestLoad_MissingCommand(t *testing.T) {
	path := writeTemp(t, `
version: "1"
services:
  - name: api
`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for missing command, got nil")
	}
}

func TestLoad_DuplicateNames(t *testing.T) {
	path := writeTemp(t, `
version: "1"
services:
  - name: api
    command: go run .
  - name: api
    command: go run .
`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for duplicate service names, got nil")
	}
}

func TestLoad_NoServices(t *testing.T) {
	path := writeTemp(t, `version: "1"
`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for empty services, got nil")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := config.Load("/nonexistent/path/grout.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
