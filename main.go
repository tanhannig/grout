package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "grout",
	Short: "A lightweight CLI for managing multi-service local dev environments",
	Long: `grout manages multi-service local dev environments via a single
configuration file, making it easy to start, stop, and monitor services.`,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
