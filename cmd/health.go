package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/yourusername/grout/config"
	"github.com/yourusername/grout/runner"
)

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Check readiness of services with configured health checks",
	RunE:  runHealth,
}

func init() {
	rootCmd.AddCommand(healthCmd)
}

func runHealth(cmd *cobra.Command, args []string) error {
	cfgPath, _ := cmd.Flags().GetString("config")
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	anyFailed := false
	for _, svc := range cfg.Services {
		if svc.HealthCheck == nil || svc.HealthCheck.URL == "" {
			fmt.Printf("  %-20s  skipped (no health check configured)\n", svc.Name)
			continue
		}

		hc := runner.DefaultHealthCheck(svc.HealthCheck.URL)
		if svc.HealthCheck.Retries > 0 {
			hc.Retries = svc.HealthCheck.Retries
		}
		if svc.HealthCheck.Interval != "" {
			if d, err := time.ParseDuration(svc.HealthCheck.Interval); err == nil {
				hc.Interval = d
			}
		}
		if svc.HealthCheck.Timeout != "" {
			if d, err := time.ParseDuration(svc.HealthCheck.Timeout); err == nil {
				hc.Timeout = d
			}
		}

		status := hc.Probe(svc.Name)
		if status.Healthy {
			fmt.Printf("  %-20s  healthy\n", svc.Name)
		} else {
			fmt.Fprintf(os.Stderr, "  %-20s  unhealthy: %v\n", svc.Name, status.Err)
			anyFailed = true
		}
	}

	if anyFailed {
		os.Exit(1)
	}
	return nil
}
