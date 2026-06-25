package runner

import (
	"testing"
	"time"
)

func TestNewHookRunner_NilMap(t *testing.T) {
	hr := NewHookRunner(nil)
	if hr == nil {
		t.Fatal("expected non-nil HookRunner")
	}
	if hr.hooks == nil {
		t.Fatal("expected hooks map to be initialized")
	}
}

func TestRun_NoHooksForService(t *testing.T) {
	hr := NewHookRunner(nil)
	if err := hr.Run("unknown", "pre-start"); err != nil {
		t.Errorf("expected nil error for missing service, got %v", err)
	}
}

func TestRun_UnknownStage(t *testing.T) {
	hr := NewHookRunner(map[string]*LifecycleHooks{
		"svc": {},
	})
	err := hr.Run("svc", "invalid-stage")
	if err == nil {
		t.Fatal("expected error for unknown stage")
	}
}

func TestRun_NilHookInStage(t *testing.T) {
	hr := NewHookRunner(map[string]*LifecycleHooks{
		"svc": {PreStart: nil},
	})
	if err := hr.Run("svc", "pre-start"); err != nil {
		t.Errorf("expected nil for nil hook, got %v", err)
	}
}

func TestRun_SuccessfulHook(t *testing.T) {
	hr := NewHookRunner(map[string]*LifecycleHooks{
		"svc": {
			PostStart: &Hook{
				Command: "echo grout-hook-ok",
				Timeout: 5 * time.Second,
			},
		},
	})
	if err := hr.Run("svc", "post-start"); err != nil {
		t.Errorf("expected successful hook run, got %v", err)
	}
}

func TestRun_FailingHook(t *testing.T) {
	hr := NewHookRunner(map[string]*LifecycleHooks{
		"svc": {
			PreStop: &Hook{
				Command: "false",
				Timeout: 5 * time.Second,
			},
		},
	})
	if err := hr.Run("svc", "pre-stop"); err == nil {
		t.Error("expected error from failing hook command")
	}
}

func TestSetHooks_UpdatesEntry(t *testing.T) {
	hr := NewHookRunner(nil)
	hr.SetHooks("api", &LifecycleHooks{
		PostStop: &Hook{Command: "echo done", Timeout: 2 * time.Second},
	})
	if err := hr.Run("api", "post-stop"); err != nil {
		t.Errorf("unexpected error after SetHooks: %v", err)
	}
}
