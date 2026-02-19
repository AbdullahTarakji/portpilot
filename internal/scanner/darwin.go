//go:build darwin

package scanner

import (
	"bufio"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type darwinScanner struct{}

// New creates a macOS scanner using lsof.
func New() (Scanner, error) {
	return &darwinScanner{}, nil
}

// Scan uses lsof to discover listening TCP and UDP ports on macOS.
func (d *darwinScanner) Scan() ([]PortInfo, error) {
	out, err := exec.Command("lsof", "-iTCP", "-iUDP", "-nP", "-sTCP:LISTEN").Output()
	if err != nil {
		// lsof may exit non-zero if some files can't be accessed (permission)
		if out == nil || len(out) == 0 {
			return nil, fmt.Errorf("running lsof: %w", err)
		}
	}

	ports, err := parseLsofOutput(string(out))
	if err != nil {
		return nil, fmt.Errorf("parsing lsof output: %w", err)
	}

	enrichWithProcessStats(ports)
	return ports, nil
}

// parseLsofOutput parses the output of `lsof -iTCP -iUDP -nP -sTCP:LISTEN`.
// Example line:
// rapportd    496 mike   4u  IPv4 0x1234   0t0  TCP *:49153 (LISTEN)
// rapportd    496 mike   5u  IPv6 0x5678   0t0  UDP *:5353
func parseLsofOutput(output string) ([]PortInfo, error) {
	var ports []PortInfo
	seen := make(map[string]bool) // deduplicate by port+proto+pid

	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 9 {
			continue
		}

		// Skip header
		if fields[0] == "COMMAND" {
			continue
		}

		processName := fields[0]
		pid, err := strconv.Atoi(fields[1])
		if err != nil {
			continue
		}
		user := fields[2]
		proto := fields[7] // TCP or UDP

		// Parse the address field (e.g., "*:49153", "127.0.0.1:8080", "[::1]:3000")
		addrField := fields[8]
		port, err := parsePortFromAddr(addrField)
		if err != nil {
			continue
		}

		state := "LISTEN"
		if len(fields) > 9 {
			state = strings.Trim(fields[9], "()")
		}
		if proto == "UDP" {
			state = "LISTEN" // UDP doesn't have LISTEN state but we show it as listening
		}

		key := fmt.Sprintf("%d:%s:%d", port, proto, pid)
		if seen[key] {
			continue
		}
		seen[key] = true

		ports = append(ports, PortInfo{
			Port:        port,
			Protocol:    proto,
			PID:         pid,
			ProcessName: processName,
			User:        user,
			State:       state,
		})
	}

	return ports, scanner.Err()
}

// parsePortFromAddr extracts the port number from lsof address field.
// Handles formats like: *:8080, 127.0.0.1:3000, [::1]:443, [::]:80
func parsePortFromAddr(addr string) (int, error) {
	idx := strings.LastIndex(addr, ":")
	if idx < 0 {
		return 0, fmt.Errorf("no port in address: %s", addr)
	}
	portStr := addr[idx+1:]
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 0, fmt.Errorf("parsing port %q: %w", portStr, err)
	}
	return port, nil
}
