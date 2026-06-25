package runner

import (
	"sync"
	"time"
)

// ServiceMetrics holds runtime statistics for a single service.
type ServiceMetrics struct {
	StartCount   int
	LastStarted  time.Time
	LastStopped  time.Time
	UptimeSecs   float64
	RestartCount int
}

// MetricsStore tracks metrics for all services.
type MetricsStore struct {
	mu      sync.RWMutex
	records map[string]*ServiceMetrics
}

// NewMetricsStore creates an empty MetricsStore.
func NewMetricsStore() *MetricsStore {
	return &MetricsStore{
		records: make(map[string]*ServiceMetrics),
	}
}

// RecordStart marks a service as started.
func (m *MetricsStore) RecordStart(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	svc := m.getOrCreate(name)
	svc.StartCount++
	svc.LastStarted = time.Now()
}

// RecordStop marks a service as stopped and accumulates uptime.
func (m *MetricsStore) RecordStop(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	svc := m.getOrCreate(name)
	now := time.Now()
	if !svc.LastStarted.IsZero() {
		svc.UptimeSecs += now.Sub(svc.LastStarted).Seconds()
	}
	svc.LastStopped = now
}

// RecordRestart increments the restart counter for a service.
func (m *MetricsStore) RecordRestart(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.getOrCreate(name).RestartCount++
}

// Get returns a copy of the metrics for the named service.
func (m *MetricsStore) Get(name string) (ServiceMetrics, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.records[name]
	if !ok {
		return ServiceMetrics{}, false
	}
	return *v, true
}

// All returns a snapshot of metrics for every tracked service.
func (m *MetricsStore) All() map[string]ServiceMetrics {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make(map[string]ServiceMetrics, len(m.records))
	for k, v := range m.records {
		out[k] = *v
	}
	return out
}

// getOrCreate is not goroutine-safe; callers must hold m.mu.
func (m *MetricsStore) getOrCreate(name string) *ServiceMetrics {
	if m.records[name] == nil {
		m.records[name] = &ServiceMetrics{}
	}
	return m.records[name]
}
