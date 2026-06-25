package runner

import (
	"fmt"
	"net"
	"time"
)

// PortStatus holds the result of a port availability check.
type PortStatus struct {
	Service string
	Port    int
	Open    bool
	Latency time.Duration
	Error   string
}

// PortChecker probes TCP ports for configured services.
type PortChecker struct {
	DialTimeout time.Duration
}

// DefaultPortChecker returns a PortChecker with sensible defaults.
func DefaultPortChecker() *PortChecker {
	return &PortChecker{
		DialTimeout: 2 * time.Second,
	}
}

// Check attempts a TCP dial to the given host:port and returns a PortStatus.
func (pc *PortChecker) Check(service string, port int) PortStatus {
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	start := time.Now()
	conn, err := net.DialTimeout("tcp", addr, pc.DialTimeout)
	latency := time.Since(start)

	if err != nil {
		return PortStatus{
			Service: service,
			Port:    port,
			Open:    false,
			Latency: latency,
			Error:   err.Error(),
		}
	}
	_ = conn.Close()
	return PortStatus{
		Service: service,
		Port:    port,
		Open:    true,
		Latency: latency,
	}
}

// CheckAll checks every service→port mapping and returns a slice of results.
func (pc *PortChecker) CheckAll(ports map[string]int) []PortStatus {
	results := make([]PortStatus, 0, len(ports))
	for svc, port := range ports {
		results = append(results, pc.Check(svc, port))
	}
	return results
}
