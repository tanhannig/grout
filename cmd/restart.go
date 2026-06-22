package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"grout/config"
	"grout/runner"
)

var restartCmd = &cobra.Command{
	Use:   "restart [service...]",
	Short: "Restart one or more running services",
	RunE:  runRestart,
}

func init() {
	rootCmd.AddCommand(restartCmd)
}

func runRestart(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load("grout.yaml")
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	r := runner.New(cfg)
	if err := r.Start(); err != nil {
		return fmt.Errorf("start runner: %w", err)
	}

	targets := args
	if len(targets) == 0 {
		for _, svc := range cfg.Services {
			targets = append(targets, svc.Name)
		}
	}

	policy := runner.DefaultRestartPolicy()
	var failed bool
	for _, name := range targets {
		fmt.Fprintf(os.Stdout, "Restarting %s...\n", name)
		if err := r.RestartWithPolicy(name, policy); err != nil {
			fmt.Fprintf(os.Stderr, "  error: %v\n", err)
			failed = true
		} else {
			fmt.Fprintf(os.Stdout, "  %s restarted.\n", name)
		}
	}

	if failed {
		return fmt.Errorf("one or more services failed to restart")
	}
	return nil
}
