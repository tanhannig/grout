package runner

import (
	"fmt"
	"time"
)

// RestartPolicy defines how a service should be restarted.
type RestartPolicy struct {
	MaxRetries int
	Delay      time.Duration
}

// DefaultRestartPolicy returns a sensible default restart policy.
func DefaultRestartPolicy() RestartPolicy {
	return RestartPolicy{
		MaxRetries: 3,
		Delay:      2 * time.Second,
	}
}

// Restart stops and restarts a named service.
func (r *Runner) Restart(name string) error {
	r.mu.Lock()
	svc, ok := r.services[name]
	r.mu.Unlock()

	if !ok {
		return fmt.Errorf("service %q not found", name)
	}

	if err := svc.stop(); err != nil {
		return fmt.Errorf("stop %q: %w", name, err)
	}

	time.Sleep(300 * time.Millisecond)

	if err := svc.start(r.logBuffer); err != nil {
		return fmt.Errorf("start %q: %w", name, err)
	}

	return nil
}

// RestartWithPolicy restarts a service, retrying up to policy.MaxRetries times.
func (r *Runner) RestartWithPolicy(name string, policy RestartPolicy) error {
	var lastErr error
	for attempt := 0; attempt <= policy.MaxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(policy.Delay)
		}
		lastErr = r.Restart(name)
		if lastErr == nil {
			return nil
		}
	}
	return fmt.Errorf("restart %q failed after %d attempts: %w", name, policy.MaxRetries+1, lastErr)
}
