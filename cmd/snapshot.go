package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"grout/runner"
)

var snapshotFile string

func init() {
	snapshotCmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Save or display a state snapshot of running services",
	}

	saveCmd := &cobra.Command{
		Use:   "save",
		Short: "Capture and save current service state to disk",
		RunE:  runSnapshotSave,
	}
	saveCmd.Flags().StringVarP(&snapshotFile, "output", "o", ".grout-state.json", "snapshot output file")

	showCmd := &cobra.Command{
		Use:   "show",
		Short: "Display the last saved snapshot",
		RunE:  runSnapshotShow,
	}
	showCmd.Flags().StringVarP(&snapshotFile, "file", "f", ".grout-state.json", "snapshot file to read")

	snapshotCmd.AddCommand(saveCmd, showCmd)
	rootCmd.AddCommand(snapshotCmd)
}

func runSnapshotSave(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	r := runner.New(cfg)
	snap := runner.CaptureFrom(r)
	store := runner.NewSnapshotStore(snapshotFile)
	if err := store.Save(snap); err != nil {
		return fmt.Errorf("failed to save snapshot: %w", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Snapshot saved to %s\n", snapshotFile)
	return nil
}

func runSnapshotShow(cmd *cobra.Command, args []string) error {
	store := runner.NewSnapshotStore(snapshotFile)
	snap, err := store.Load()
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("no snapshot found at %s", snapshotFile)
		}
		return fmt.Errorf("failed to load snapshot: %w", err)
	}

	if jsonOut, _ := cmd.Flags().GetBool("json"); jsonOut {
		return json.NewEncoder(cmd.OutOrStdout()).Encode(snap)
	}

	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "SERVICE\tSTATUS\tPID\tRESTARTS\tCAPTURED")
	for _, svc := range snap.Services {
		fmt.Fprintf(w, "%s\t%s\t%d\t%d\t%s\n",
			svc.Name, svc.Status, svc.PID, svc.Restarts,
			svc.CapturedAt.Format(time.RFC3339))
	}
	return w.Flush()
}
