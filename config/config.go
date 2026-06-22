package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Service represents a single runnable service.
type Service struct {
	Name        string            `yaml:"name"`
	Command     string            `yaml:"command"`
	Dir         string            `yaml:"dir"`
	Env         map[string]string `yaml:"env"`
	HealthCheck *HealthCheck      `yaml:"health_check"`
	Restart     *RestartPolicy    `yaml:"restart"`
}

// HealthCheck holds configuration for HTTP health probing.
type HealthCheck struct {
	URL      string `yaml:"url"`
	Interval int    `yaml:"interval"`
	Timeout  int    `yaml:"timeout"`
}

// RestartPolicy controls automatic restarts on failure.
type RestartPolicy struct {
	OnFailure  bool `yaml:"on_failure"`
	MaxRetries int  `yaml:"max_retries"`
}

// Config is the top-level grout configuration.
type Config struct {
	Services []Service `yaml:"services"`
}

const DefaultFile = "grout.yaml"

// Load reads and validates a grout config file.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}
	if err := validate(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func validate(cfg *Config) error {
	if len(cfg.Services) == 0 {
		return errors.New("config must define at least one service")
	}
	seen := make(map[string]bool, len(cfg.Services))
	for _, svc := range cfg.Services {
		if svc.Command == "" {
			return fmt.Errorf("service %q: command is required", svc.Name)
		}
		if seen[svc.Name] {
			return fmt.Errorf("duplicate service name: %q", svc.Name)
		}
		seen[svc.Name] = true
	}
	return nil
}
