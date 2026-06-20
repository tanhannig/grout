package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/user/grout/config"
	"github.com/user/grout/runner"
)

func init() {
	registerCommand("up", runUp)
}

// runUp loads the grout config and starts all defined services.
// It blocks until an interrupt signal is received, then stops all services.
func runUp(args []string) error {
	cfgPath := "grout.yaml"
	if len(args) > 0 {
		cfgPath = args[0]
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	if len(cfg.Services) == 0 {
		return fmt.Errorf("no services defined in %s", cfgPath)
	}

	fmt.Printf("Starting %d service(s)...\n", len(cfg.Services))
	for _, svc := range cfg.Services {
		fmt.Printf("  → %s: %s\n", svc.Name, svc.Command)
	}

	r := runner.New(cfg)
	if err := r.Start(); err != nil {
		return fmt.Errorf("starting services: %w", err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	<-sigCh
	fmt.Println("\nShutting down services...")
	r.Stop()
	fmt.Println("All services stopped.")
	return nil
}
