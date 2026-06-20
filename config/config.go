package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const DefaultConfigFile = "grout.yaml"

// Service represents a single managed service definition.
type Service struct {
	Name    string            `yaml:"name"`
	Command string            `yaml:"command"`
	Dir     string            `yaml:"dir"`
	Env     map[string]string `yaml:"env"`
	Port    int               `yaml:"port"`
}

// Config is the top-level grout configuration.
type Config struct {
	Version  string    `yaml:"version"`
	Services []Service `yaml:"services"`
}

// Load reads and parses a grout config file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file %q: %w", path, err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Validate checks that the config contains required fields.
func (c *Config) Validate() error {
	if len(c.Services) == 0 {
		return fmt.Errorf("config must define at least one service")
	}
	seen := make(map[string]bool)
	for i, svc := range c.Services {
		if svc.Name == "" {
			return fmt.Errorf("service at index %d is missing a name", i)
		}
		if svc.Command == "" {
			return fmt.Errorf("service %q is missing a command", svc.Name)
		}
		if seen[svc.Name] {
			return fmt.Errorf("duplicate service name %q", svc.Name)
		}
		seen[svc.Name] = true
	}
	return nil
}
