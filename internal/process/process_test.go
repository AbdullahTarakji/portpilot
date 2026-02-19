package process

import (
	"os"
	"syscall"
	"testing"
)

func TestKillInvalidPID(t *testing.T) {
	err := Kill(-1, os.Kill)
	if err == nil {
		t.Error("expected error for invalid PID")
	}
}

func TestKillNonExistent(t *testing.T) {
	err := Kill(99999999, os.Kill)
	if err == nil {
		t.Error("expected error for non-existent PID")
	}
}

func TestGetDetailsInvalidPID(t *testing.T) {
	_, err := GetDetails(99999999)
	if err == nil {
		t.Error("expected error for non-existent PID")
	}
}

func TestGetDetailsSelf(t *testing.T) {
	info, err := GetDetails(os.Getpid())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.PID != os.Getpid() {
		t.Errorf("PID: got %d, want %d", info.PID, os.Getpid())
	}
	if info.Command == "" {
		t.Error("Command should not be empty for current process")
	}
}

func TestIsRunning(t *testing.T) {
	if !IsRunning(os.Getpid()) {
		t.Error("current process should be running")
	}
	if IsRunning(99999999) {
		t.Error("PID 99999999 should not be running")
	}
}

func TestParseSignal(t *testing.T) {
	tests := []struct {
		input string
		want  syscall.Signal
		err   bool
	}{
		{"SIGTERM", syscall.SIGTERM, false},
		{"TERM", syscall.SIGTERM, false},
		{"SIGKILL", syscall.SIGKILL, false},
		{"kill", syscall.SIGKILL, false},
		{"15", syscall.Signal(15), false},
		{"bogus", 0, true},
	}

	for _, tt := range tests {
		sig, err := ParseSignal(tt.input)
		if tt.err {
			if err == nil {
				t.Errorf("ParseSignal(%q): expected error", tt.input)
			}
			continue
		}
		if err != nil {
			t.Errorf("ParseSignal(%q): unexpected error: %v", tt.input, err)
			continue
		}
		if sig != tt.want {
			t.Errorf("ParseSignal(%q): got %v, want %v", tt.input, sig, tt.want)
		}
	}
}

func TestExtractProcessName(t *testing.T) {
	tests := []struct {
		cmd  string
		want string
	}{
		{"/usr/local/bin/node server.js", "node"},
		{"postgres -D /var/lib/pg", "postgres"},
		{"", ""},
	}
	for _, tt := range tests {
		got := extractProcessName(tt.cmd)
		if got != tt.want {
			t.Errorf("extractProcessName(%q): got %q, want %q", tt.cmd, got, tt.want)
		}
	}
}
