package runner

import (
	"net"
	"testing"
	"time"
)

func TestDefaultPortChecker_Fields(t *testing.T) {
	pc := DefaultPortChecker()
	if pc.DialTimeout != 2*time.Second {
		t.Errorf("expected 2s dial timeout, got %v", pc.DialTimeout)
	}
}

func TestCheck_OpenPort(t *testing.T) {
	// Start a temporary listener so the port is genuinely open.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to open listener: %v", err)
	}
	defer ln.Close()

	port := ln.Addr().(*net.TCPAddr).Port
	pc := DefaultPortChecker()
	status := pc.Check("svc", port)

	if !status.Open {
		t.Errorf("expected port %d to be open, got error: %s", port, status.Error)
	}
	if status.Service != "svc" {
		t.Errorf("expected service name 'svc', got %q", status.Service)
	}
	if status.Port != port {
		t.Errorf("expected port %d, got %d", port, status.Port)
	}
	if status.Latency <= 0 {
		t.Error("expected positive latency")
	}
}

func TestCheck_ClosedPort(t *testing.T) {
	pc := &PortChecker{DialTimeout: 200 * time.Millisecond}
	// Port 1 is almost certainly closed in test environments.
	status := pc.Check("db", 1)

	if status.Open {
		t.Error("expected port 1 to be closed")
	}
	if status.Error == "" {
		t.Error("expected non-empty error string for closed port")
	}
}

func TestCheckAll_ReturnsAllServices(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to open listener: %v", err)
	}
	defer ln.Close()

	openPort := ln.Addr().(*net.TCPAddr).Port
	ports := map[string]int{
		"web": openPort,
		"db":  1, // closed
	}

	pc := &PortChecker{DialTimeout: 200 * time.Millisecond}
	results := pc.CheckAll(ports)

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	openCount := 0
	for _, r := range results {
		if r.Open {
			openCount++
		}
	}
	if openCount != 1 {
		t.Errorf("expected 1 open port, got %d", openCount)
	}
}
