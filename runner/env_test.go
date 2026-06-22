package runner

import (
	"strings"
	"testing"

	"github.com/user/grout/config"
)

func envMap(env []string) map[string]string {
	m := make(map[string]string, len(env))
	for _, e := range env {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			m[parts[0]] = parts[1]
		}
	}
	return m
}

func TestResolve_NoOverrides(t *testing.T) {
	r := NewEnvResolver()
	svc := config.Service{Name: "web", Command: "echo hi"}
	result := r.Resolve(svc)
	if len(result) != len(r.base) {
		t.Errorf("expected %d entries, got %d", len(r.base), len(result))
	}
}

func TestResolve_AddsNewKey(t *testing.T) {
	r := &EnvResolver{base: []string{"PATH=/usr/bin"}}
	svc := config.Service{
		Name:    "api",
		Command: "go run .",
		Env:     map[string]string{"PORT": "8080"},
	}
	result := envMap(r.Resolve(svc))
	if result["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %q", result["PORT"])
	}
	if result["PATH"] != "/usr/bin" {
		t.Errorf("expected PATH to be preserved")
	}
}

func TestResolve_OverridesExistingKey(t *testing.T) {
	r := &EnvResolver{base: []string{"PORT=3000", "HOST=localhost"}}
	svc := config.Service{
		Name:    "worker",
		Command: "./worker",
		Env:     map[string]string{"PORT": "9090"},
	}
	result := envMap(r.Resolve(svc))
	if result["PORT"] != "9090" {
		t.Errorf("expected PORT=9090 after override, got %q", result["PORT"])
	}
	if result["HOST"] != "localhost" {
		t.Errorf("expected HOST to be preserved")
	}
}

func TestResolve_NoDuplicateKeys(t *testing.T) {
	r := &EnvResolver{base: []string{"FOO=bar", "BAZ=qux"}}
	svc := config.Service{
		Name:    "svc",
		Command: "echo",
		Env:     map[string]string{"FOO": "overridden"},
	}
	result := r.Resolve(svc)
	count := 0
	for _, e := range result {
		if strings.HasPrefix(e, "FOO=") {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected exactly 1 FOO entry, got %d", count)
	}
}
