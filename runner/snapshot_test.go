package runner

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSnapshotStore_SaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "state.json")
	store := NewSnapshotStore(path)

	snap := StateSnapshot{
		CapturedAt: time.Now(),
		Services: map[string]ServiceSnapshot{
			"web": {
				Name:    "web",
				Status:  "running",
				PID:     1234,
				Restarts: 1,
			},
		},
	}

	if err := store.Save(snap); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := store.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(loaded.Services) != 1 {
		t.Errorf("expected 1 service, got %d", len(loaded.Services))
	}
	svc, ok := loaded.Services["web"]
	if !ok {
		t.Fatal("expected 'web' service in snapshot")
	}
	if svc.PID != 1234 {
		t.Errorf("expected PID 1234, got %d", svc.PID)
	}
	if svc.Restarts != 1 {
		t.Errorf("expected 1 restart, got %d", svc.Restarts)
	}
}

func TestSnapshotStore_LoadMissing(t *testing.T) {
	store := NewSnapshotStore("/nonexistent/path/state.json")
	_, err := store.Load()
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestSnapshotStore_SaveAtomic(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "state.json")
	store := NewSnapshotStore(path)

	for i := 0; i < 5; i++ {
		snap := StateSnapshot{
			Services: map[string]ServiceSnapshot{
				"svc": {Name: "svc", Restarts: i},
			},
		}
		if err := store.Save(snap); err != nil {
			t.Fatalf("Save iteration %d failed: %v", i, err)
		}
	}

	// tmp file should not linger
	if _, err := os.Stat(path + ".tmp"); !os.IsNotExist(err) {
		t.Error("tmp file should not exist after atomic save")
	}

	loaded, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Services["svc"].Restarts != 4 {
		t.Errorf("expected 4 restarts, got %d", loaded.Services["svc"].Restarts)
	}
}
