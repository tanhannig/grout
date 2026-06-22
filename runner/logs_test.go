package runner

import (
	"bytes"
	"strings"
	"testing"
)

func TestLogBuffer_AppendAndEntries(t *testing.T) {
	buf := NewLogBuffer(3)
	buf.Append(LogEntry{Service: "svc", Line: "a"})
	buf.Append(LogEntry{Service: "svc", Line: "b"})

	entries := buf.Entries()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Line != "a" || entries[1].Line != "b" {
		t.Errorf("unexpected entries: %+v", entries)
	}
}

func TestLogBuffer_Eviction(t *testing.T) {
	buf := NewLogBuffer(2)
	buf.Append(LogEntry{Line: "a"})
	buf.Append(LogEntry{Line: "b"})
	buf.Append(LogEntry{Line: "c"})

	entries := buf.Entries()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries after eviction, got %d", len(entries))
	}
	if entries[0].Line != "b" || entries[1].Line != "c" {
		t.Errorf("expected b,c got %s,%s", entries[0].Line, entries[1].Line)
	}
}

func TestServiceWriter_WritesLines(t *testing.T) {
	var out bytes.Buffer
	buf := NewLogBuffer(10)
	w := NewServiceWriter("api", buf, false, &out)

	w.Write([]byte("hello\nworld\n"))

	entries := buf.Entries()
	if len(entries) != 2 {
		t.Fatalf("expected 2 log entries, got %d", len(entries))
	}
	if entries[0].Line != "hello" || entries[1].Line != "world" {
		t.Errorf("unexpected lines: %+v", entries)
	}
	if !strings.Contains(out.String(), "[api] hello") {
		t.Errorf("output missing prefix: %s", out.String())
	}
}

func TestServiceWriter_PartialLine(t *testing.T) {
	var out bytes.Buffer
	buf := NewLogBuffer(10)
	w := NewServiceWriter("db", buf, false, &out)

	w.Write([]byte("par"))
	if len(buf.Entries()) != 0 {
		t.Error("expected no entries for partial line")
	}

	w.Write([]byte("tial\n"))
	entries := buf.Entries()
	if len(entries) != 1 || entries[0].Line != "partial" {
		t.Errorf("expected partial line buffered, got %+v", entries)
	}
}

func TestServiceWriter_StderrFlag(t *testing.T) {
	var out bytes.Buffer
	buf := NewLogBuffer(10)
	w := NewServiceWriter("web", buf, true, &out)
	w.Write([]byte("err\n"))

	entries := buf.Entries()
	if len(entries) != 1 || !entries[0].IsStderr {
		t.Error("expected IsStderr=true")
	}
}
