package runner

import (
	"fmt"
	"os"
	"os/exec"
	"sync"

	"github.com/user/grout/config"
)

// Service wraps a config.Service with its running process.
type Service struct {
	Config  config.Service
	Cmd     *exec.Cmd
	Stopped chan struct{}
}

// Runner manages the lifecycle of all services.
type Runner struct {
	services []*Service
	mu       sync.Mutex
}

// New creates a Runner from a loaded config.
func New(cfg *config.Config) *Runner {
	r := &Runner{}
	for _, svc := range cfg.Services {
		r.services = append(r.services, &Service{
			Config:  svc,
			Stopped: make(chan struct{}),
		})
	}
	return r
}

// Start launches all services concurrently.
func (r *Runner) Start() error {
	var wg sync.WaitGroup
	errCh := make(chan error, len(r.services))

	for _, svc := range r.services {
		wg.Add(1)
		go func(s *Service) {
			defer wg.Done()
			if err := r.startService(s); err != nil {
				errCh <- fmt.Errorf("service %q: %w", s.Config.Name, err)
			}
		}(svc)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Runner) startService(s *Service) error {
	cmd := exec.Command("sh", "-c", s.Config.Command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if s.Config.Dir != "" {
		cmd.Dir = s.Config.Dir
	}

	r.mu.Lock()
	s.Cmd = cmd
	r.mu.Unlock()

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start: %w", err)
	}

	go func() {
		_ = cmd.Wait()
		close(s.Stopped)
	}()

	return nil
}

// Stop terminates all running services.
func (r *Runner) Stop() {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, svc := range r.services {
		if svc.Cmd != nil && svc.Cmd.Process != nil {
			_ = svc.Cmd.Process.Kill()
		}
	}
}
