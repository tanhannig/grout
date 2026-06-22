package runner

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/grout/config"
)

// EnvResolver builds the environment for a service process by merging
// the host environment with per-service overrides defined in the config.
type EnvResolver struct {
	base []string
}

// NewEnvResolver creates an EnvResolver seeded with the current process env.
func NewEnvResolver() *EnvResolver {
	return &EnvResolver{base: os.Environ()}
}

// Resolve returns the final env slice for the given service.
// Service-level variables override any matching host variable.
func (r *EnvResolver) Resolve(svc config.Service) []string {
	if len(svc.Env) == 0 {
		return r.base
	}

	// Index overrides by key for O(1) lookup.
	overrides := make(map[string]string, len(svc.Env))
	for k, v := range svc.Env {
		overrides[k] = v
	}

	// Start with base, replacing any keys present in overrides.
	result := make([]string, 0, len(r.base)+len(overrides))
	seen := make(map[string]bool, len(overrides))

	for _, entry := range r.base {
		parts := strings.SplitN(entry, "=", 2)
		key := parts[0]
		if val, ok := overrides[key]; ok {
			result = append(result, fmt.Sprintf("%s=%s", key, val))
			seen[key] = true
		} else {
			result = append(result, entry)
		}
	}

	// Append any override keys not already in base.
	for k, v := range overrides {
		if !seen[k] {
			result = append(result, fmt.Sprintf("%s=%s", k, v))
		}
	}

	return result
}
