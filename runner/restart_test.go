package runner

import (
	"testing"
	"time"
)

func TestRestart_UnknownService(t *testing.T) {
	cfg := makeConfig("svc1", "echo hello")
	r := New(cfg)

	err := r.Restart("nonexistent")
	if err == nil {
		t.Fatal("expected error for unknown service")
	}
}

func TestRestart_StopsAndRestarts(t *testing.T) {
	cfg := makeConfig("svc1", "sleep 5")
	r := New(cfg)

	if err := r.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}
	t.Cleanup(func() { _ = r.Stop() })

	time.Sleep(100 * time.Millisecond)

	if err := r.Restart("svc1"); err != nil {
		t.Fatalf("Restart: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	status := r.Status()
	if len(status) == 0 {
		t.Fatal("expected status entries")
	}
	for _, s := range status {
		if s.Name == "svc1" && s.State != StateRunning {
			t.Errorf("expected svc1 running after restart, got %s", s.State)
		}
	}
}

func TestDefaultRestartPolicy(t *testing.T) {
	p := DefaultRestartPolicy()
	if p.MaxRetries != 3 {
		t.Errorf("expected MaxRetries=3, got %d", p.MaxRetries)
	}
	if p.Delay != 2*time.Second {
		t.Errorf("expected Delay=2s, got %v", p.Delay)
	}
}

func TestRestartWithPolicy_InvalidCommand(t *testing.T) {
	cfg := makeConfig("bad", "sleep 5")
	r := New(cfg)

	// Override command to something invalid after creation
	r.mu.Lock()
	r.services["bad"].cmd = "not-a-real-command-xyz"
	r.mu.Unlock()

	policy := RestartPolicy{MaxRetries: 1, Delay: 10 * time.Millisecond}
	err := r.RestartWithPolicy("bad", policy)
	if err == nil {
		t.Fatal("expected error for invalid command")
	}
}
