package runner

import (
	"fmt"
	"sync"
	"time"
)

// Hook represents a shell command to run at a lifecycle stage.
type Hook struct {
	Command string
	Timeout time.Duration
}

// LifecycleHooks holds optional pre/post hooks for a service.
type LifecycleHooks struct {
	PreStart  *Hook
	PostStart *Hook
	PreStop   *Hook
	PostStop  *Hook
}

// HookRunner executes lifecycle hooks for services.
type HookRunner struct {
	mu    sync.Mutex
	hooks map[string]*LifecycleHooks
}

// NewHookRunner creates a HookRunner with the given hooks map.
func NewHookRunner(hooks map[string]*LifecycleHooks) *HookRunner {
	if hooks == nil {
		hooks = make(map[string]*LifecycleHooks)
	}
	return &HookRunner{hooks: hooks}
}

// Run executes a hook for the given service and stage.
// stage must be one of: pre-start, post-start, pre-stop, post-stop.
func (h *HookRunner) Run(service, stage string) error {
	h.mu.Lock()
	lc, ok := h.hooks[service]
	h.mu.Unlock()
	if !ok || lc == nil {
		return nil
	}

	var hook *Hook
	switch stage {
	case "pre-start":
		hook = lc.PreStart
	case "post-start":
		hook = lc.PostStart
	case "pre-stop":
		hook = lc.PreStop
	case "post-stop":
		hook = lc.PostStop
	default:
		return fmt.Errorf("unknown lifecycle stage: %s", stage)
	}

	if hook == nil || hook.Command == "" {
		return nil
	}

	timeout := hook.Timeout
	if timeout == 0 {
		timeout = 10 * time.Second
	}

	return runWithTimeout(hook.Command, timeout)
}

// SetHooks registers or replaces hooks for a service.
func (h *HookRunner) SetHooks(service string, lc *LifecycleHooks) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.hooks[service] = lc
}
