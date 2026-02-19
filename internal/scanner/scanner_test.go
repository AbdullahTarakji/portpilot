package scanner

import (
	"testing"
	"time"
)

func TestParseProcessStats(t *testing.T) {
	line := "1.5 0.8 Thu Feb 19 04:00:00 2026 /usr/local/bin/node server.js"
	s, err := parseProcessStats(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.cpu != 1.5 {
		t.Errorf("cpu: got %f, want 1.5", s.cpu)
	}
	if s.mem != 0.8 {
		t.Errorf("mem: got %f, want 0.8", s.mem)
	}
	expected := time.Date(2026, 2, 19, 4, 0, 0, 0, time.UTC)
	if !s.startTime.Equal(expected) {
		t.Errorf("startTime: got %v, want %v", s.startTime, expected)
	}
	if s.command != "/usr/local/bin/node server.js" {
		t.Errorf("command: got %q, want %q", s.command, "/usr/local/bin/node server.js")
	}
}

func TestParseProcessStatsTooFewFields(t *testing.T) {
	_, err := parseProcessStats("1.5 0.8")
	if err == nil {
		t.Error("expected error for too few fields")
	}
}
