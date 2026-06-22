package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/grout/config"
	"github.com/yourorg/grout/runner"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the current status of all services",
	RunE:  runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func runStatus(cmd *cobra.Command, args []string) error {
	cfgPath, err := cmd.Flags().GetString("config")
	if err != nil || cfgPath == "" {
		cfgPath = "grout.yaml"
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading config: %v\n", err)
		return err
	}

	r, err := runner.New(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating runner: %v\n", err)
		return err
	}

	statuses := r.Status()
	fmt.Print(runner.FormatStatus(statuses))
	return nil
}
