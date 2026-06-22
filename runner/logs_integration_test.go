package runner_test

import (
	"strings"
	"testing"
	"time"

	"grout/config"
	"grout/runner"
)

func TestRunner_LogsPopulated(t *testing.T) {
	cfg := &config.Config{
		Services: []config.Service{
			{Name: "echo-svc", Command: "echo hello from grout"},
		},
	}

	r := runner.New(cfg)
	if err := r.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	// Give the process time to produce output and exit.
	time.Sleep(300 * time.Millisecond)
	r.Stop()

	entries := r.Logs().Entries()
	if len(entries) == 0 {
		t.Fatal("expected log entries, got none")
	}

	found := false
	for _, e := range entries {
		if e.Service == "echo-svc" && strings.Contains(e.Line, "hello from grout") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected 'hello from grout' in logs, entries: %+v", entries)
	}
}

func TestRunner_LogsFilterByService(t *testing.T) {
	cfg := &config.Config{
		Services: []config.Service{
			{Name: "svc-a", Command: "echo from-a"},
			{Name: "svc-b", Command: "echo from-b"},
		},
	}

	r := runner.New(cfg)
	if err := r.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	time.Sleep(300 * time.Millisecond)
	r.Stop()

	allEntries := r.Logs().Entries()
	for _, e := range allEntries {
		if e.Service != "svc-a" && e.Service != "svc-b" {
			t.Errorf("unexpected service in logs: %s", e.Service)
		}
	}
}
