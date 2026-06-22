package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/grout-dev/grout/config"
	"github.com/grout-dev/grout/runner"
)

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop all running services",
	Long:  `Stop all services defined in grout.yaml that are currently running.`,
	RunE:  runDown,
}

func init() {
	rootCmd.AddCommand(downCmd)
}

// runDown loads the config, initializes the runner, and stops all services.
func runDown(cmd *cobra.Command, args []string) error {
	cfgPath := filepath.Join(".", "grout.yaml")

	// Check that config file exists before proceeding
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		return fmt.Errorf("grout.yaml not found in current directory — run 'grout init' to create one")
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	r, err := runner.New(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize runner: %w", err)
	}

	fmt.Println("Stopping services...")

	stopped, errs := r.Stop()

	if len(stopped) > 0 {
		for _, name := range stopped {
			fmt.Printf("  ✓ stopped %s\n", name)
		}
	}

	if len(errs) > 0 {
		for name, stopErr := range errs {
			fmt.Fprintf(os.Stderr, "  ✗ error stopping %s: %v\n", name, stopErr)
		}
		return fmt.Errorf("one or more services failed to stop cleanly")
	}

	if len(stopped) == 0 {
		fmt.Println("No running services to stop.")
		return nil
	}

	fmt.Println("All services stopped.")
	return nil
}
