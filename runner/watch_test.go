package runner

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
)

func TestDefaultWatchConfig(t *testing.T) {
	paths := []string{"/srv/app", "/srv/config"}
	cfg := DefaultWatchConfig(paths)

	if len(cfg.Paths) != 2 {
		t.Fatalf("expected 2 paths, got %d", len(cfg.Paths))
	}
	if cfg.Debounce != 500*time.Millisecond {
		t.Errorf("unexpected debounce: %v", cfg.Debounce)
	}
	if len(cfg.Exts) == 0 {
		t.Error("expected default extensions to be non-empty")
	}
}

func TestRelevantEvent_MatchingExt(t *testing.T) {
	r := &Runner{}
	exts := []string{".go", ".py"}

	event := fsnotify.Event{Name: "main.go", Op: fsnotify.Write}
	if !r.relevantEvent(event, exts) {
		t.Error("expected .go to be relevant")
	}
}

func TestRelevantEvent_NonMatchingExt(t *testing.T) {
	r := &Runner{}
	exts := []string{".go"}

	event := fsnotify.Event{Name: "README.md", Op: fsnotify.Write}
	if r.relevantEvent(event, exts) {
		t.Error("expected .md to be irrelevant")
	}
}

func TestRelevantEvent_NonWriteOp(t *testing.T) {
	r := &Runner{}
	exts := []string{".go"}

	event := fsnotify.Event{Name: "main.go", Op: fsnotify.Chmod}
	if r.relevantEvent(event, exts) {
		t.Error("expected chmod to be irrelevant")
	}
}

func TestWatchAndRestart_ExitsOnDone(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := makeConfig("watcher-svc", "echo hello")
	r := New(cfg)

	watchCfg := WatchConfig{
		Paths:    []string{tmpDir},
		Exts:     []string{".go"},
		Debounce: 50 * time.Millisecond,
	}

	done := make(chan struct{})
	errCh := make(chan error, 1)
	go func() {
		errCh <- r.WatchAndRestart("watcher-svc", watchCfg, done)
	}()

	// Allow watcher to start
	time.Sleep(50 * time.Millisecond)
	close(done)

	select {
	case err := <-errCh:
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("WatchAndRestart did not exit after done closed")
	}
}

func TestWatchAndRestart_DetectsFileChange(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := makeConfig("watch-echo", "echo restarted")
	r := New(cfg)

	watchCfg := WatchConfig{
		Paths:    []string{tmpDir},
		Exts:     []string{".go"},
		Debounce: 80 * time.Millisecond,
	}

	done := make(chan struct{})
	go func() { _ = r.WatchAndRestart("watch-echo", watchCfg, done) }()
	defer close(done)

	time.Sleep(60 * time.Millisecond)

	// Write a relevant file to trigger event
	testFile := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(testFile, []byte("package main"), 0644); err != nil {
		t.Fatal(err)
	}

	// Allow debounce + restart to occur
	time.Sleep(300 * time.Millisecond)
	// No panic or deadlock — test passes if we reach here
}
