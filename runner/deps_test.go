package runner

import (
	"testing"

	"github.com/grout-dev/grout/config"
)

func services(defs ...config.Service) []config.Service { return defs }

func svc(name string, deps ...string) config.Service {
	return config.Service{Name: name, Command: "echo " + name, DependsOn: deps}
}

func TestOrderServices_NoDeps(t *testing.T) {
	ordered, err := OrderServices(services(svc("a"), svc("b"), svc("c")))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ordered) != 3 {
		t.Fatalf("expected 3 services, got %d", len(ordered))
	}
}

func TestOrderServices_LinearChain(t *testing.T) {
	// c depends on b depends on a => order must be a, b, c
	ordered, err := OrderServices(services(svc("c", "b"), svc("b", "a"), svc("a")))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ordered[0] != "a" || ordered[1] != "b" || ordered[2] != "c" {
		t.Fatalf("unexpected order: %v", ordered)
	}
}

func TestOrderServices_DiamondDep(t *testing.T) {
	// d depends on b and c; b and c both depend on a
	ordered, err := OrderServices(services(
		svc("d", "b", "c"),
		svc("b", "a"),
		svc("c", "a"),
		svc("a"),
	))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// a must appear before b, c; b and c before d
	pos := func(name string) int {
		for i, n := range ordered {
			if n == name {
				return i
			}
		}
		return -1
	}
	if pos("a") >= pos("b") || pos("a") >= pos("c") || pos("b") >= pos("d") || pos("c") >= pos("d") {
		t.Fatalf("ordering violated diamond constraint: %v", ordered)
	}
}

func TestOrderServices_Cycle(t *testing.T) {
	_, err := OrderServices(services(svc("a", "b"), svc("b", "a")))
	if err == nil {
		t.Fatal("expected cycle error, got nil")
	}
	if _, ok := err.(*DepsError); !ok {
		t.Fatalf("expected *DepsError, got %T", err)
	}
}

func TestOrderServices_UnknownDep(t *testing.T) {
	_, err := OrderServices(services(svc("a", "ghost")))
	if err == nil {
		t.Fatal("expected error for unknown dependency, got nil")
	}
}
