package scanner

import (
	"testing"
	"time"
)

func TestParseLsofOutput(t *testing.T) {
	input := `COMMAND     PID   USER   FD   TYPE             DEVICE SIZE/OFF NODE NAME
rapportd    496   mike    4u  IPv4 0x1234      0t0  TCP *:49153 (LISTEN)
node      12345   mike    5u  IPv6 0x5678      0t0  TCP [::1]:3000 (LISTEN)
postgres  54321   mike    6u  IPv4 0xabcd      0t0  TCP 127.0.0.1:5432 (LISTEN)
mDNSRespo   100   _mdns   7u  IPv4 0x9999      0t0  UDP *:5353
`

	ports, err := parseLsofOutput(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(ports) != 4 {
		t.Fatalf("expected 4 ports, got %d", len(ports))
	}

	tests := []struct {
		idx         int
		port        int
		proto       string
		pid         int
		processName string
		user        string
		state       string
	}{
		{0, 49153, "TCP", 496, "rapportd", "mike", "LISTEN"},
		{1, 3000, "TCP", 12345, "node", "mike", "LISTEN"},
		{2, 5432, "TCP", 54321, "postgres", "mike", "LISTEN"},
		{3, 5353, "UDP", 100, "mDNSRespo", "_mdns", "LISTEN"},
	}

	for _, tt := range tests {
		p := ports[tt.idx]
		if p.Port != tt.port {
			t.Errorf("[%d] port: got %d, want %d", tt.idx, p.Port, tt.port)
		}
		if p.Protocol != tt.proto {
			t.Errorf("[%d] proto: got %s, want %s", tt.idx, p.Protocol, tt.proto)
		}
		if p.PID != tt.pid {
			t.Errorf("[%d] pid: got %d, want %d", tt.idx, p.PID, tt.pid)
		}
		if p.ProcessName != tt.processName {
			t.Errorf("[%d] name: got %s, want %s", tt.idx, p.ProcessName, tt.processName)
		}
		if p.User != tt.user {
			t.Errorf("[%d] user: got %s, want %s", tt.idx, p.User, tt.user)
		}
		if p.State != tt.state {
			t.Errorf("[%d] state: got %s, want %s", tt.idx, p.State, tt.state)
		}
	}
}

func TestParseLsofOutputEmpty(t *testing.T) {
	ports, err := parseLsofOutput("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ports) != 0 {
		t.Fatalf("expected 0 ports, got %d", len(ports))
	}
}

func TestParseLsofOutputDedup(t *testing.T) {
	input := `COMMAND  PID USER FD TYPE DEVICE SIZE/OFF NODE NAME
node   123 mike 4u IPv4 0x1 0t0 TCP *:3000 (LISTEN)
node   123 mike 5u IPv6 0x2 0t0 TCP *:3000 (LISTEN)
`
	ports, err := parseLsofOutput(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ports) != 1 {
		t.Fatalf("expected 1 port (deduped), got %d", len(ports))
	}
}

func TestParsePortFromAddr(t *testing.T) {
	tests := []struct {
		addr string
		port int
		err  bool
	}{
		{"*:8080", 8080, false},
		{"127.0.0.1:3000", 3000, false},
		{"[::1]:443", 443, false},
		{"[::]:80", 80, false},
		{"noport", 0, true},
	}

	for _, tt := range tests {
		port, err := parsePortFromAddr(tt.addr)
		if tt.err && err == nil {
			t.Errorf("parsePortFromAddr(%q): expected error", tt.addr)
		}
		if !tt.err && err != nil {
			t.Errorf("parsePortFromAddr(%q): unexpected error: %v", tt.addr, err)
		}
		if port != tt.port {
			t.Errorf("parsePortFromAddr(%q): got %d, want %d", tt.addr, port, tt.port)
		}
	}
}

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
