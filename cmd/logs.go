package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"grout/runner"
)

var (
	logsN       int
	logsService string
)

func init() {
	logsCmd := &cobra.Command{
		Use:   "logs [service]",
		Short: "Show buffered log output for services",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runLogs,
	}
	logsCmd.Flags().IntVarP(&logsN, "tail", "n", 50, "Number of recent lines to show (0 = all)")
	RootCmd.AddCommand(logsCmd)
}

func runLogs(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		logsService = args[0]
	}

	if globalRunner == nil {
		fmt.Fprintln(os.Stderr, "No services running. Start with 'grout up'.")
		return nil
	}

	entries := globalRunner.Logs().Entries()

	if logsService != "" {
		filtered := entries[:0]
		for _, e := range entries {
			if strings.EqualFold(e.Service, logsService) {
				filtered = append(filtered, e)
			}
		}
		entries = filtered
	}

	if logsN > 0 && len(entries) > logsN {
		entries = entries[len(entries)-logsN:]
	}

	if len(entries) == 0 {
		fmt.Println("No log entries found.")
		return nil
	}

	for _, e := range entries {
		level := "out"
		if e.IsStderr {
			level = "err"
		}
		fmt.Printf("%s  [%s] (%s) %s\n",
			e.Timestamp.Format("15:04:05.000"),
			e.Service,
			level,
			e.Line,
		)
	}
	return nil
}
