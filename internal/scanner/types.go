package scanner

import "time"

// PortInfo holds information about a listening port and its associated process.
type PortInfo struct {
	Port        int       `json:"port"`
	Protocol    string    `json:"protocol"`
	PID         int       `json:"pid"`
	ProcessName string    `json:"process_name"`
	User        string    `json:"user"`
	State       string    `json:"state"`
	Command     string    `json:"command"`
	CPU         float64   `json:"cpu_percent"`
	Mem         float64   `json:"mem_percent"`
	StartTime   time.Time `json:"start_time"`
}

// ProcessInfo holds detailed information about a process.
type ProcessInfo struct {
	PID         int       `json:"pid"`
	Name        string    `json:"name"`
	User        string    `json:"user"`
	Command     string    `json:"command"`
	CPU         float64   `json:"cpu_percent"`
	Mem         float64   `json:"mem_percent"`
	StartTime   time.Time `json:"start_time"`
	ParentPID   int       `json:"parent_pid"`
	NumThreads  int       `json:"num_threads"`
	WorkingDir  string    `json:"working_dir"`
}
