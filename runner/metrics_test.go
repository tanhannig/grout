package runner

import (
	"testing"
	"time"
)

func TestMetricsStore_RecordStart(t *testing.T) {
	m := NewMetricsStore()
	m.RecordStart("web")
	m.RecordStart("web")

	snap, ok := m.Get("web")
	if !ok {
		t.Fatal("expected metrics for 'web'")
	}
	if snap.StartCount != 2 {
		t.Errorf("expected StartCount=2, got %d", snap.StartCount)
	}
	if snap.LastStarted.IsZero() {
		t.Error("expected LastStarted to be set")
	}
}

func TestMetricsStore_RecordStop_AccumulatesUptime(t *testing.T) {
	m := NewMetricsStore()
	m.RecordStart("api")
	time.Sleep(10 * time.Millisecond)
	m.RecordStop("api")

	snap, _ := m.Get("api")
	if snap.UptimeSecs <= 0 {
		t.Errorf("expected positive uptime, got %f", snap.UptimeSecs)
	}
	if snap.LastStopped.IsZero() {
		t.Error("expected LastStopped to be set")
	}
}

func TestMetricsStore_RecordRestart(t *testing.T) {
	m := NewMetricsStore()
	m.RecordRestart("worker")
	m.RecordRestart("worker")
	m.RecordRestart("worker")

	snap, ok := m.Get("worker")
	if !ok {
		t.Fatal("expected metrics for 'worker'")
	}
	if snap.RestartCount != 3 {
		t.Errorf("expected RestartCount=3, got %d", snap.RestartCount)
	}
}

func TestMetricsStore_Get_Unknown(t *testing.T) {
	m := NewMetricsStore()
	_, ok := m.Get("ghost")
	if ok {
		t.Error("expected ok=false for unknown service")
	}
}

func TestMetricsStore_All(t *testing.T) {
	m := NewMetricsStore()
	m.RecordStart("a")
	m.RecordStart("b")
	m.RecordStop("b")

	all := m.All()
	if len(all) != 2 {
		t.Errorf("expected 2 entries, got %d", len(all))
	}
	if _, ok := all["a"]; !ok {
		t.Error("expected entry for 'a'")
	}
	if _, ok := all["b"]; !ok {
		t.Error("expected entry for 'b'")
	}
}

func TestMetricsStore_StopWithoutStart(t *testing.T) {
	m := NewMetricsStore()
	// Should not panic
	m.RecordStop("orphan")
	snap, ok := m.Get("orphan")
	if !ok {
		t.Fatal("expected entry created on stop")
	}
	if snap.UptimeSecs != 0 {
		t.Errorf("expected zero uptime when no start recorded, got %f", snap.UptimeSecs)
	}
}
