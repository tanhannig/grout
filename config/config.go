package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// HealthCheckConfig holds optional readiness probe settings per service.
type HealthCheckConfig struct {
	URL      string `yaml:"url"`
	Retries  int    `yaml:"retries"`
	Interval string `yaml:"interval"`
	Timeout  string `yaml:"timeout"`
}

// Service represents a single managed process.
type Service struct {
	Name        string            `yaml:"name"`
	Command     string            `yaml:"command"`
	Dir         string            `yaml:"dir"`
	Env         map[string]string `yaml:"env"`
	HealthCheck *HealthCheckConfig `yaml:"health_check"`
}

// Config is the top-level grout configuration.
type Config struct {
	Services []Service `yaml:"services"`
}

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
	seen := make(map[string]bool)
	for _, svc := range cfg.Services {
		if svc.Command == "" {
			return fmt.Errorf("service %q missing command", svc.Name)
		}
		if seen[svc.Name] {
			return fmt.Errorf("duplicate service name: %q", svc.Name)
		}
		seen[svc.Name] = true
	}
	return nil
}
