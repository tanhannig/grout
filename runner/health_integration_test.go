package runner_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/grout/runner"
)

func TestHealthCheck_IntegrationWithRunner(t *testing.T) {
	ready := make(chan struct{})
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-ready:
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	}))
	defer ts.Close()

	// Close the ready channel after a short delay to simulate startup.
	go func() {
		time.Sleep(30 * time.Millisecond)
		close(ready)
	}()

	hc := runner.HealthCheck{
		URL:      ts.URL,
		Interval: 20 * time.Millisecond,
		Timeout:  500 * time.Millisecond,
		Retries:  10,
	}
	status := hc.Probe("svc")
	if !status.Healthy {
		t.Errorf("expected service to become healthy, got err: %v", status.Err)
	}
}

// TestHealthCheck_NeverReady verifies that Probe returns an unhealthy status
// when the service never becomes ready within the configured retry budget.
func TestHealthCheck_NeverReady(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer ts.Close()

	hc := runner.HealthCheck{
		URL:      ts.URL,
		Interval: 10 * time.Millisecond,
		Timeout:  100 * time.Millisecond,
		Retries:  3,
	}
	status := hc.Probe("svc")
	if status.Healthy {
		t.Error("expected service to remain unhealthy, but got healthy status")
	}
	if status.Err == nil {
		t.Error("expected a non-nil error for an unhealthy service")
	}
}

func TestDefaultHealthCheck_Fields(t *testing.T) {
	hc := runner.DefaultHealthCheck("http://localhost:3000/ready")
	if hc.URL != "http://localhost:3000/ready" {
		t.Errorf("unexpected URL: %s", hc.URL)
	}
	if hc.Retries <= 0 {
		t.Error("expected positive retries")
	}
	if hc.Timeout <= 0 {
		t.Error("expected positive timeout")
	}
}
