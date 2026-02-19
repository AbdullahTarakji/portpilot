//go:build linux

package scanner

import (
	"bufio"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type linuxScanner struct{}

// New creates a Linux scanner using ss.
func New() (Scanner, error) {
	return &linuxScanner{}, nil
}

// Scan uses ss to discover listening TCP and UDP ports on Linux.
func (l *linuxScanner) Scan() ([]PortInfo, error) {
	out, err := exec.Command("ss", "-tulnp").Output()
	if err != nil {
		if out == nil || len(out) == 0 {
			return nil, fmt.Errorf("running ss: %w", err)
		}
	}

	ports, err := parseSSOutput(string(out))
	if err != nil {
		return nil, fmt.Errorf("parsing ss output: %w", err)
	}

	enrichWithProcessStats(ports)
	return ports, nil
}

// parseSSOutput parses the output of `ss -tulnp`.
// Example lines:
// Netid  State   Recv-Q  Send-Q   Local Address:Port   Peer Address:Port  Process
// tcp    LISTEN  0       128      0.0.0.0:22            0.0.0.0:*          users:(("sshd",pid=1234,fd=3))
// udp    UNCONN  0       0        0.0.0.0:5353          0.0.0.0:*          users:(("avahi-daemon",pid=567,fd=12))
func parseSSOutput(output string) ([]PortInfo, error) {
	var ports []PortInfo
	seen := make(map[string]bool)

	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}

		// Skip header
		proto := strings.ToUpper(fields[0])
		if proto != "TCP" && proto != "UDP" {
			continue
		}

		state := fields[1]
		if proto == "UDP" && state == "UNCONN" {
			state = "LISTEN"
		}

		// Local address is field 4
		localAddr := fields[4]
		port, err := parsePortFromAddr(localAddr)
		if err != nil {
			continue
		}

		// Parse process info from the last field if it contains users:(...)
		var pid int
		var processName string
		for _, f := range fields[5:] {
			if strings.HasPrefix(f, "users:") || strings.Contains(f, "pid=") {
				pid, processName = parseSSProcess(f)
			}
		}

		key := fmt.Sprintf("%d:%s:%d", port, proto, pid)
		if seen[key] {
			continue
		}
		seen[key] = true

		info := PortInfo{
			Port:        port,
			Protocol:    proto,
			PID:         pid,
			ProcessName: processName,
			State:       state,
		}

		// Try to get user from /proc if we have a PID
		if pid > 0 {
			if user, err := getProcessUser(pid); err == nil {
				info.User = user
			}
		}

		ports = append(ports, info)
	}

	return ports, scanner.Err()
}

// parseSSProcess extracts PID and process name from ss process field.
// Input format: users:(("sshd",pid=1234,fd=3))
func parseSSProcess(field string) (int, string) {
	// Extract process name
	var name string
	if start := strings.Index(field, "((\""); start >= 0 {
		rest := field[start+3:]
		if end := strings.Index(rest, "\""); end >= 0 {
			name = rest[:end]
		}
	}

	// Extract PID
	var pid int
	if pidIdx := strings.Index(field, "pid="); pidIdx >= 0 {
		rest := field[pidIdx+4:]
		if end := strings.IndexAny(rest, ",)"); end >= 0 {
			rest = rest[:end]
		}
		pid, _ = strconv.Atoi(rest)
	}

	return pid, name
}

// getProcessUser reads the owner of a process from /proc.
func getProcessUser(pid int) (string, error) {
	out, err := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "user=").Output()
	if err != nil {
		return "", fmt.Errorf("getting user for pid %d: %w", pid, err)
	}
	return strings.TrimSpace(string(out)), nil
}

// parsePortFromAddr extracts the port number from an ss address field.
// Handles formats like: 0.0.0.0:22, [::]:80, *:5353
func parsePortFromAddr(addr string) (int, error) {
	idx := strings.LastIndex(addr, ":")
	if idx < 0 {
		return 0, fmt.Errorf("no port in address: %s", addr)
	}
	portStr := addr[idx+1:]
	// ss may show * for wildcard
	if portStr == "*" {
		return 0, fmt.Errorf("wildcard port")
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 0, fmt.Errorf("parsing port %q: %w", portStr, err)
	}
	return port, nil
}
