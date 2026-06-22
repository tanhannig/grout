package runner

import (
	"fmt"
	"net/http"
	"time"
)

// HealthCheck defines a readiness probe for a service.
type HealthCheck struct {
	URL      string
	Interval time.Duration
	Timeout  time.Duration
	Retries  int
}

// HealthStatus represents the result of a health check.
type HealthStatus struct {
	Service string
	Healthy bool
	Err     error
}

// DefaultHealthCheck returns a HealthCheck with sensible defaults.
func DefaultHealthCheck(url string) HealthCheck {
	return HealthCheck{
		URL:      url,
		Interval: 2 * time.Second,
		Timeout:  1 * time.Second,
		Retries:  5,
	}
}

// Probe attempts to reach the health check URL, retrying up to hc.Retries times.
func (hc HealthCheck) Probe(serviceName string) HealthStatus {
	client := &http.Client{Timeout: hc.Timeout}
	var lastErr error
	for i := 0; i < hc.Retries; i++ {
		resp, err := client.Get(hc.URL)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode < 500 {
				return HealthStatus{Service: serviceName, Healthy: true}
			}
			lastErr = fmt.Errorf("status %d", resp.StatusCode)
		} else {
			lastErr = err
		}
		time.Sleep(hc.Interval)
	}
	return HealthStatus{Service: serviceName, Healthy: false, Err: lastErr}
}
