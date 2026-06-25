package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"grout/config"
	"grout/runner"
)

var metricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "Show runtime metrics for all services",
	RunE:  runMetrics,
}

func init() {
	rootCmd.AddCommand(metricsCmd)
}

func runMetrics(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load("grout.yaml")
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	r, err := runner.New(cfg)
	if err != nil {
		return fmt.Errorf("creating runner: %w", err)
	}

	store := r.Metrics()
	all := store.All()

	if len(all) == 0 {
		fmt.Println("No metrics recorded yet. Have you run 'grout up'?")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "SERVICE\tSTARTS\tRESTARTS\tUPTIME (s)\tLAST STARTED")
	fmt.Fprintln(w, "-------\t------\t--------\t----------\t------------")

	for name, m := range all {
		lastStarted := "-"
		if !m.LastStarted.IsZero() {
			lastStarted = m.LastStarted.Format("15:04:05")
		}
		fmt.Fprintf(w, "%s\t%d\t%d\t%.1f\t%s\n",
			name, m.StartCount, m.RestartCount, m.UptimeSecs, lastStarted)
	}

	return w.Flush()
}
