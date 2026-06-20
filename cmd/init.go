package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var initTemplate = `version: "1"
services:
  - name: example
    command: echo "hello from grout"
    dir: .
    env:
      ENV: development
    port: 8080
`

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a default grout.yaml in the current directory",
	RunE:  runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	const filename = "grout.yaml"

	if _, err := os.Stat(filename); err == nil {
		return fmt.Errorf("%s already exists; remove it before running init", filename)
	}

	if err := os.WriteFile(filename, []byte(initTemplate), 0644); err != nil {
		return fmt.Errorf("creating %s: %w", filename, err)
	}

	fmt.Printf("Created %s — edit it to define your services.\n", filename)
	return nil
}
