package runner

import (
	"testing"
	"time"

	"github.com/user/grout/config"
)

func makeConfig(services []config.Service) *config.Config {
	return &config.Config{Services: services}
}

func TestNew_CreatesServices(t *testing.T) {
	cfg := makeConfig([]config.Service{
		{Name: "web", Command: "echo web"},
		{Name: "api", Command: "echo api"},
	})
	r := New(cfg)
	if len(r.services) != 2 {
		t.Fatalf("expected 2 services, got %d", len(r.services))
	}
	if r.services[0].Config.Name != "web" {
		t.Errorf("expected first service name 'web', got %q", r.services[0].Config.Name)
	}
}

func TestStart_RunsAndCompletes(t *testing.T) {
	cfg := makeConfig([]config.Service{
		{Name: "hello", Command: "echo hello"},
	})
	r := New(cfg)
	if err := r.Start(); err != nil {
		t.Fatalf("Start returned error: %v", err)
	}

	svc := r.services[0]
	select {
	case <-svc.Stopped:
		// process exited cleanly
	case <-time.After(3 * time.Second):
		t.Fatal("service did not stop within timeout")
	}
}

func TestStart_InvalidCommand(t *testing.T) {
	cfg := makeConfig([]config.Service{
		{Name: "bad", Command: "nonexistent_binary_xyz_123"},
	})
	r := New(cfg)
	// sh -c will start but exit non-zero; no error from Start itself
	if err := r.Start(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStop_KillsProcesses(t *testing.T) {
	cfg := makeConfig([]config.Service{
		{Name: "sleep", Command: "sleep 60"},
	})
	r := New(cfg)
	if err := r.Start(); err != nil {
		t.Fatalf("Start returned error: %v", err)
	}

	time.Sleep(100 * time.Millisecond)
	r.Stop()

	svc := r.services[0]
	select {
	case <-svc.Stopped:
		// killed successfully
	case <-time.After(3 * time.Second):
		t.Fatal("service was not killed within timeout")
	}
}
