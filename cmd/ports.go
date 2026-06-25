package cmd

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"grout/config"
	"grout/runner"
)

var portsCmd = &cobra.Command{
	Use:   "ports",
	Short: "Check TCP port availability for all configured services",
	RunE:  runPorts,
}

func init() {
	rootCmd.AddCommand(portsCmd)
}

func runPorts(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load("grout.yaml")
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	// Collect services that declare a port.
	portMap := make(map[string]int)
	for _, svc := range cfg.Services {
		if svc.Port > 0 {
			portMap[svc.Name] = svc.Port
		}
	}

	if len(portMap) == 0 {
		fmt.Println("No services have a port configured.")
		return nil
	}

	pc := runner.DefaultPortChecker()
	results := pc.CheckAll(portMap)

	// Sort by service name for stable output.
	sort.Slice(results, func(i, j int) bool {
		return results[i].Service < results[j].Service
	})

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "SERVICE\tPORT\tSTATUS\tLATENCY")
	for _, r := range results {
		status := "closed"
		if r.Open {
			status = "open"
		}
		fmt.Fprintf(w, "%s\t%d\t%s\t%v\n", r.Service, r.Port, status, r.Latency.Round(1000000))
	}
	return w.Flush()
}
