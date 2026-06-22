package runner

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestDefaultHealthCheck(t *testing.T) {
	hc := DefaultHealthCheck("http://localhost:8080/health")
	if hc.Retries != 5 {
		t.Errorf("expected 5 retries, got %d", hc.Retries)
	}
	if hc.Interval != 2*time.Second {
		t.Errorf("unexpected interval: %v", hc.Interval)
	}
}

func TestProbe_Healthy(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	hc := HealthCheck{URL: ts.URL, Interval: 10 * time.Millisecond, Timeout: 500 * time.Millisecond, Retries: 3}
	status := hc.Probe("web")
	if !status.Healthy {
		t.Errorf("expected healthy, got err: %v", status.Err)
	}
	if status.Service != "web" {
		t.Errorf("expected service name 'web', got %q", status.Service)
	}
}

func TestProbe_Unhealthy_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	hc := HealthCheck{URL: ts.URL, Interval: 5 * time.Millisecond, Timeout: 500 * time.Millisecond, Retries: 2}
	status := hc.Probe("api")
	if status.Healthy {
		t.Error("expected unhealthy for 500 response")
	}
}

func TestProbe_Unhealthy_NoServer(t *testing.T) {
	hc := HealthCheck{URL: "http://127.0.0.1:19999/health", Interval: 5 * time.Millisecond, Timeout: 50 * time.Millisecond, Retries: 2}
	status := hc.Probe("db")
	if status.Healthy {
		t.Error("expected unhealthy when server not running")
	}
	if status.Err == nil {
		t.Error("expected non-nil error")
	}
}
