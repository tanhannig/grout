package runner

// Exported helpers for snapshot integration testing.

// SnapshotPath is the default file name used when no path is specified.
const SnapshotPath = ".grout-state.json"

// NewDefaultSnapshotStore returns a SnapshotStore using the default path.
func NewDefaultSnapshotStore() *SnapshotStore {
	return NewSnapshotStore(SnapshotPath)
}
