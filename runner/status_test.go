package runner

import (
	"strings"
	"testing"
)

func TestStatus_InitialState(t *testing.T) {
	cfg := makeConfig("svc-a", "echo hello")
	r, err := New(cfg)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	statuses := r.Status()
	if len(statuses) != 1 {
		t.Fatalf("expected 1 status, got %d", len(statuses))
	}
	s := statuses[0]
	if s.Name != "svc-a" {
		t.Errorf("expected name 'svc-a', got %q", s.Name)
	}
	if s.Running {
		t.Errorf("expected service to not be running initially")
	}
	if s.PID != 0 {
		t.Errorf("expected PID 0 initially, got %d", s.PID)
	}
}

func TestFormatStatus_Empty(t *testing.T) {
	out := FormatStatus(nil)
	if !strings.Contains(out, "No services") {
		t.Errorf("expected 'No services' message, got: %q", out)
	}
}

func TestFormatStatus_Stopped(t *testing.T) {
	statuses := []ServiceStatus{
		{Name: "web", Running: false, PID: 0},
	}
	out := FormatStatus(statuses)
	if !strings.Contains(out, "web") {
		t.Errorf("expected service name 'web' in output")
	}
	if !strings.Contains(out, "stopped") {
		t.Errorf("expected 'stopped' status in output")
	}
}

func TestFormatStatus_WithError(t *testing.T) {
	statuses := []ServiceStatus{
		{Name: "api", Running: false, Error: "exit status 1"},
	}
	out := FormatStatus(statuses)
	if !strings.Contains(out, "exit status 1") {
		t.Errorf("expected error message in output, got: %q", out)
	}
}

func TestFormatStatus_Headers(t *testing.T) {
	statuses := []ServiceStatus{
		{Name: "db", Running: true, PID: 1234},
	}
	out := FormatStatus(statuses)
	for _, header := range []string{"NAME", "STATUS", "PID", "UPTIME"} {
		if !strings.Contains(out, header) {
			t.Errorf("expected header %q in output", header)
		}
	}
}
