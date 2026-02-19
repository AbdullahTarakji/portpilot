// Package process provides process management utilities for killing and inspecting processes.
package process

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// Details holds extended information about a running process.
type Details struct {
	PID        int       `json:"pid"`
	Name       string    `json:"name"`
	User       string    `json:"user"`
	Command    string    `json:"command"`
	CPU        float64   `json:"cpu_percent"`
	Mem        float64   `json:"mem_percent"`
	StartTime  time.Time `json:"start_time"`
	ParentPID  int       `json:"parent_pid"`
	NumThreads int       `json:"num_threads"`
}

// Kill sends the specified signal to a process.
func Kill(pid int, signal os.Signal) error {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("finding process %d: %w", pid, err)
	}

	if err := proc.Signal(signal); err != nil {
		return fmt.Errorf("sending signal %v to pid %d: %w", signal, pid, err)
	}

	return nil
}

// KillByPort finds and kills the process listening on the given port.
// Returns the PID that was killed.
func KillByPort(port int, signal os.Signal) (int, error) {
	pid, err := findPIDByPort(port)
	if err != nil {
		return 0, err
	}

	if err := Kill(pid, signal); err != nil {
		return pid, err
	}

	return pid, nil
}

// GetDetails returns detailed information about a process.
func GetDetails(pid int) (*Details, error) {
	out, err := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "pid,ppid,%cpu,%mem,user,lstart,command").Output()
	if err != nil {
		return nil, fmt.Errorf("getting details for pid %d: %w", pid, err)
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) < 2 {
		return nil, fmt.Errorf("process %d not found", pid)
	}

	return parseDetails(lines[1])
}

func parseDetails(line string) (*Details, error) {
	fields := strings.Fields(line)
	if len(fields) < 11 {
		return nil, fmt.Errorf("unexpected ps output: %q", line)
	}

	d := &Details{}

	pid, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, fmt.Errorf("parsing pid: %w", err)
	}
	d.PID = pid

	ppid, err := strconv.Atoi(fields[1])
	if err != nil {
		return nil, fmt.Errorf("parsing ppid: %w", err)
	}
	d.ParentPID = ppid

	cpu, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return nil, fmt.Errorf("parsing cpu: %w", err)
	}
	d.CPU = cpu

	mem, err := strconv.ParseFloat(fields[3], 64)
	if err != nil {
		return nil, fmt.Errorf("parsing mem: %w", err)
	}
	d.Mem = mem

	d.User = fields[4]

	// lstart: "Day Mon DD HH:MM:SS YYYY" (5 fields, index 5..9)
	timeStr := strings.Join(fields[5:10], " ")
	t, err := time.Parse("Mon Jan 2 15:04:05 2006", timeStr)
	if err == nil {
		d.StartTime = t
	}

	d.Command = strings.Join(fields[10:], " ")
	d.Name = extractProcessName(d.Command)

	return d, nil
}

func extractProcessName(command string) string {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return ""
	}
	name := parts[0]
	// Strip path
	if idx := strings.LastIndex(name, "/"); idx >= 0 {
		name = name[idx+1:]
	}
	return name
}

// IsRunning checks whether a process with the given PID is running.
func IsRunning(pid int) bool {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	// Signal 0 tests for process existence without actually signaling
	err = proc.Signal(syscall.Signal(0))
	return err == nil
}

func findPIDByPort(port int) (int, error) {
	out, err := exec.Command("lsof", "-iTCP:"+strconv.Itoa(port), "-iUDP:"+strconv.Itoa(port), "-nP", "-sTCP:LISTEN", "-t").Output()
	if err != nil {
		return 0, fmt.Errorf("no process found on port %d: %w", port, err)
	}

	pidStr := strings.TrimSpace(strings.Split(string(out), "\n")[0])
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return 0, fmt.Errorf("parsing pid for port %d: %w", port, err)
	}

	return pid, nil
}

// ParseSignal converts a signal name or number string to an os.Signal.
func ParseSignal(s string) (os.Signal, error) {
	signals := map[string]syscall.Signal{
		"SIGTERM": syscall.SIGTERM,
		"TERM":    syscall.SIGTERM,
		"SIGKILL": syscall.SIGKILL,
		"KILL":    syscall.SIGKILL,
		"SIGINT":  syscall.SIGINT,
		"INT":     syscall.SIGINT,
		"SIGHUP":  syscall.SIGHUP,
		"HUP":     syscall.SIGHUP,
		"SIGUSR1": syscall.SIGUSR1,
		"USR1":    syscall.SIGUSR1,
		"SIGUSR2": syscall.SIGUSR2,
		"USR2":    syscall.SIGUSR2,
	}

	upper := strings.ToUpper(s)
	if sig, ok := signals[upper]; ok {
		return sig, nil
	}

	// Try parsing as a number
	num, err := strconv.Atoi(s)
	if err == nil {
		return syscall.Signal(num), nil
	}

	return nil, fmt.Errorf("unknown signal: %s", s)
}
