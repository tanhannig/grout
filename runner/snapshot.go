package runner

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// ServiceSnapshot captures the state of a service at a point in time.
type ServiceSnapshot struct {
	Name      string        `json:"name"`
	Status    string        `json:"status"`
	PID       int           `json:"pid,omitempty"`
	Uptime    time.Duration `json:"uptime_ns,omitempty"`
	Restarts  int           `json:"restarts"`
	LastError string        `json:"last_error,omitempty"`
	CapturedAt time.Time    `json:"captured_at"`
}

// StateSnapshot holds snapshots for all services.
type StateSnapshot struct {
	CapturedAt time.Time                  `json:"captured_at"`
	Services   map[string]ServiceSnapshot `json:"services"`
}

// SnapshotStore persists and retrieves state snapshots.
type SnapshotStore struct {
	mu   sync.RWMutex
	path string
}

// NewSnapshotStore creates a SnapshotStore backed by the given file path.
func NewSnapshotStore(path string) *SnapshotStore {
	return &SnapshotStore{path: path}
}

// Save writes the snapshot to disk atomically.
func (s *SnapshotStore) Save(snap StateSnapshot) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	snap.CapturedAt = time.Now()
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}

// Load reads the most recent snapshot from disk.
func (s *SnapshotStore) Load() (StateSnapshot, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var snap StateSnapshot
	data, err := os.ReadFile(s.path)
	if err != nil {
		return snap, err
	}
	err = json.Unmarshal(data, &snap)
	return snap, err
}

// CaptureFrom builds a StateSnapshot from the runner's current status.
func CaptureFrom(r *Runner) StateSnapshot {
	snap := StateSnapshot{
		CapturedAt: time.Now(),
		Services:   make(map[string]ServiceSnapshot),
	}
	for _, svc := range r.Status() {
		errStr := ""
		if svc.Err != nil {
			errStr = svc.Err.Error()
		}
		snap.Services[svc.Name] = ServiceSnapshot{
			Name:       svc.Name,
			Status:     string(svc.State),
			PID:        svc.PID,
			Uptime:     svc.Uptime,
			Restarts:   svc.Restarts,
			LastError:  errStr,
			CapturedAt: time.Now(),
		}
	}
	return snap
}
