package scanner

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Scanner defines the interface for port scanning across platforms.
type Scanner interface {
	// Scan returns all listening ports and their associated process information.
	Scan() ([]PortInfo, error)
}

// New creates a platform-appropriate Scanner.
// Implemented in platform-specific files (darwin.go, linux.go).

// enrichWithProcessStats augments port entries with CPU, memory, and command info from ps.
func enrichWithProcessStats(ports []PortInfo) {
	pids := make(map[int]bool)
	for _, p := range ports {
		if p.PID > 0 {
			pids[p.PID] = true
		}
	}
	if len(pids) == 0 {
		return
	}

	stats := make(map[int]processStats)
	for pid := range pids {
		s, err := getProcessStats(pid)
		if err == nil {
			stats[pid] = s
		}
	}

	for i := range ports {
		if s, ok := stats[ports[i].PID]; ok {
			ports[i].CPU = s.cpu
			ports[i].Mem = s.mem
			if ports[i].Command == "" {
				ports[i].Command = s.command
			}
			ports[i].StartTime = s.startTime
		}
	}
}

type processStats struct {
	cpu       float64
	mem       float64
	command   string
	startTime time.Time
}

func getProcessStats(pid int) (processStats, error) {
	out, err := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "%cpu,%mem,lstart,command").Output()
	if err != nil {
		return processStats{}, fmt.Errorf("ps for pid %d: %w", pid, err)
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) < 2 {
		return processStats{}, fmt.Errorf("unexpected ps output for pid %d", pid)
	}

	// The output line looks like:
	//  0.0  0.1 Thu Jan  2 15:04:05 2025 /usr/bin/some-command --flag
	line := strings.TrimSpace(lines[1])
	return parseProcessStats(line)
}

// parseProcessStats parses a single line of `ps -o %cpu,%mem,lstart,command` output.
func parseProcessStats(line string) (processStats, error) {
	var s processStats

	fields := strings.Fields(line)
	if len(fields) < 7 {
		return s, fmt.Errorf("too few fields in ps output: %q", line)
	}

	cpu, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return s, fmt.Errorf("parsing cpu: %w", err)
	}
	s.cpu = cpu

	mem, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return s, fmt.Errorf("parsing mem: %w", err)
	}
	s.mem = mem

	// lstart format: "Day Mon DD HH:MM:SS YYYY" (5 fields starting at index 2)
	if len(fields) >= 7 {
		timeStr := strings.Join(fields[2:7], " ")
		t, err := time.Parse("Mon Jan 2 15:04:05 2006", timeStr)
		if err == nil {
			s.startTime = t
		}
	}

	if len(fields) > 7 {
		s.command = strings.Join(fields[7:], " ")
	}

	return s, nil
}
