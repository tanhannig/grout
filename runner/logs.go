package runner

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// LogEntry represents a single log line from a service.
type LogEntry struct {
	Service   string
	Timestamp time.Time
	Line      string
	IsStderr  bool
}

// LogBuffer holds recent log entries for each service.
type LogBuffer struct {
	mu      sync.RWMutex
	entries []LogEntry
	max     int
}

// NewLogBuffer creates a LogBuffer that retains up to max entries.
func NewLogBuffer(max int) *LogBuffer {
	return &LogBuffer{max: max, entries: make([]LogEntry, 0, max)}
}

// Append adds a new log entry, evicting the oldest if at capacity.
func (b *LogBuffer) Append(e LogEntry) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if len(b.entries) >= b.max {
		b.entries = b.entries[1:]
	}
	b.entries = append(b.entries, e)
}

// Entries returns a copy of all buffered log entries.
func (b *LogBuffer) Entries() []LogEntry {
	b.mu.RLock()
	defer b.mu.RUnlock()
	out := make([]LogEntry, len(b.entries))
	copy(out, b.entries)
	return out
}

// ServiceWriter is an io.Writer that prefixes each line with the service name
// and appends entries to a LogBuffer.
type ServiceWriter struct {
	service  string
	buffer   *LogBuffer
	isStderr bool
	out      io.Writer
	remainder []byte
}

// NewServiceWriter creates a ServiceWriter that also mirrors output to out.
func NewServiceWriter(service string, buffer *LogBuffer, isStderr bool, out io.Writer) *ServiceWriter {
	if out == nil {
		out = os.Stdout
	}
	return &ServiceWriter{service: service, buffer: buffer, isStderr: isStderr, out: out}
}

func (w *ServiceWriter) Write(p []byte) (int, error) {
	data := append(w.remainder, p...)
	w.remainder = nil
	for {
		idx := indexByte(data, '\n')
		if idx < 0 {
			w.remainder = append([]byte(nil), data...)
			break
		}
		line := string(data[:idx])
		data = data[idx+1:]
		entry := LogEntry{Service: w.service, Timestamp: time.Now(), Line: line, IsStderr: w.isStderr}
		w.buffer.Append(entry)
		fmt.Fprintf(w.out, "[%s] %s\n", w.service, line)
	}
	return len(p), nil
}

func indexByte(b []byte, c byte) int {
	for i, v := range b {
		if v == c {
			return i
		}
	}
	return -1
}
