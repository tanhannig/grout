package runner

// Logs returns the shared LogBuffer for all services managed by this Runner.
// It is safe to call concurrently and is intended for use by CLI commands.
func (r *Runner) Logs() *LogBuffer {
	return r.logBuffer
}

// logBufferField wires a LogBuffer into the Runner during construction.
// Called from New() in runner.go.
func initLogBuffer(r *Runner) {
	const maxEntries = 500
	r.logBuffer = NewLogBuffer(maxEntries)
}
