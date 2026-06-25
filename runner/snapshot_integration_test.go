package runner

import (
	"path/filepath"
	"testing"
	"time"

	"grout/config"
)

func TestCaptureFrom_ReflectsRunnerState(t *testing.T) {
	cfg := makeConfig("echo hello")
	r := New(cfg)

	if err := r.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	time.Sleep(150 * time.Millisecond)
	_ = r.Stop()

	snap := CaptureFrom(r)
	if len(snap.Services) == 0 {
		t.Fatal("expected at least one service in snapshot")
	}
	for name, svc := range snap.Services {
		if svc.Name != name {
			t.Errorf("snapshot name mismatch: key %q, svc.Name %q", name, svc.Name)
		}
	}
}

func TestCaptureFrom_PersistAndReload(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "grout-state.json")
	store := NewSnapshotStore(path)

	cfg := &config.Config{
		Services: []config.Service{
			{Name: "api", Command: "echo api"},
			{Name: "worker", Command: "echo worker"},
		},
	}
	r := New(cfg)

	snap := CaptureFrom(r)
	if err := store.Save(snap); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := store.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(loaded.Services) != 2 {
		t.Errorf("expected 2 services, got %d", len(loaded.Services))
	}
	if _, ok := loaded.Services["api"]; !ok {
		t.Error("expected 'api' in loaded snapshot")
	}
	if _, ok := loaded.Services["worker"]; !ok {
		t.Error("expected 'worker' in loaded snapshot")
	}
}
