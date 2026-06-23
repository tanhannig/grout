package runner

import (
	"fmt"
	"strings"

	"github.com/grout-dev/grout/config"
)

// DepsError is returned when dependency ordering fails.
type DepsError struct {
	Cycle []string
}

func (e *DepsError) Error() string {
	return fmt.Sprintf("dependency cycle detected: %s", strings.Join(e.Cycle, " -> "))
}

// OrderServices returns service names in dependency-resolved order using
// a topological sort (Kahn's algorithm). Returns DepsError if a cycle exists.
func OrderServices(services []config.Service) ([]string, error) {
	// Build adjacency and in-degree maps.
	inDegree := make(map[string]int, len(services))
	deps := make(map[string][]string, len(services))
	names := make(map[string]bool, len(services))

	for _, svc := range services {
		names[svc.Name] = true
		inDegree[svc.Name] = 0
	}

	for _, svc := range services {
		for _, dep := range svc.DependsOn {
			if !names[dep] {
				return nil, fmt.Errorf("service %q depends on unknown service %q", svc.Name, dep)
			}
			deps[dep] = append(deps[dep], svc.Name)
			inDegree[svc.Name]++
		}
	}

	// Collect nodes with no incoming edges.
	queue := []string{}
	for _, svc := range services {
		if inDegree[svc.Name] == 0 {
			queue = append(queue, svc.Name)
		}
	}

	ordered := make([]string, 0, len(services))
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		ordered = append(ordered, node)
		for _, dependent := range deps[node] {
			inDegree[dependent]--
			if inDegree[dependent] == 0 {
				queue = append(queue, dependent)
			}
		}
	}

	if len(ordered) != len(services) {
		// Detect cycle members.
		cycleMembers := []string{}
		for name, deg := range inDegree {
			if deg > 0 {
				cycleMembers = append(cycleMembers, name)
			}
		}
		return nil, &DepsError{Cycle: cycleMembers}
	}

	return ordered, nil
}
