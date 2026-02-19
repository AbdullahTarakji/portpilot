// Package main is the entry point for the portpilot CLI.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"github.com/AbdullahTarakji/portpilot/internal/config"
	"github.com/AbdullahTarakji/portpilot/internal/process"
	"github.com/AbdullahTarakji/portpilot/internal/scanner"
	"github.com/AbdullahTarakji/portpilot/internal/tui"
)

// Build-time variables injected via ldflags.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if err := rootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}

func rootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "portpilot",
		Short: "PortPilot — manage ports and processes",
		Long:  "A CLI + TUI tool for discovering, inspecting, and managing listening ports and their processes.",
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := scanner.New()
			if err != nil {
				return err
			}
			cfg, err := config.Load()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: config error: %v\n", err)
				cfg = config.DefaultConfig()
			}
			return tui.Run(s, cfg)
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	root.AddCommand(
		listCmd(),
		killCmd(),
		checkCmd(),
		watchCmd(),
		versionCmd(),
	)

	return root
}

func listCmd() *cobra.Command {
	var (
		jsonOutput  bool
		portFilter  int
		procFilter  string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List listening ports",
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := scanner.New()
			if err != nil {
				return err
			}

			ports, err := s.Scan()
			if err != nil {
				return fmt.Errorf("scanning ports: %w", err)
			}

			ports = applyFilters(ports, portFilter, procFilter)

			if jsonOutput {
				return printJSON(ports)
			}
			printTable(ports)
			return nil
		},
	}

	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")
	cmd.Flags().IntVar(&portFilter, "port", 0, "Filter by port number")
	cmd.Flags().StringVar(&procFilter, "process", "", "Filter by process name")

	return cmd
}

func killCmd() *cobra.Command {
	var (
		force     bool
		signalStr string
	)

	cmd := &cobra.Command{
		Use:   "kill <port>",
		Short: "Kill the process on a specified port",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			port := 0
			if _, err := fmt.Sscanf(args[0], "%d", &port); err != nil {
				return fmt.Errorf("invalid port: %s", args[0])
			}

			sig, err := process.ParseSignal(signalStr)
			if err != nil {
				return fmt.Errorf("invalid signal: %w", err)
			}

			// Find what's on the port first
			s, scanErr := scanner.New()
			if scanErr != nil {
				return scanErr
			}
			ports, scanErr := s.Scan()
			if scanErr != nil {
				return fmt.Errorf("scanning: %w", scanErr)
			}

			var target *scanner.PortInfo
			for _, p := range ports {
				if p.Port == port {
					target = &p
					break
				}
			}

			if target == nil {
				fmt.Printf("No process found on port %d\n", port)
				return nil
			}

			if !force {
				fmt.Printf("Kill %q (PID %d) on port %d? [y/N] ", target.ProcessName, target.PID, target.Port)
				var answer string
				_, _ = fmt.Scanln(&answer)
				if strings.ToLower(answer) != "y" {
					fmt.Println("Cancelled.")
					return nil
				}
			}

			if err := process.Kill(target.PID, sig); err != nil {
				return fmt.Errorf("killing process: %w", err)
			}

			fmt.Printf("Sent %v to PID %d (%s) on port %d\n", sig, target.PID, target.ProcessName, target.Port)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Skip confirmation")
	cmd.Flags().StringVarP(&signalStr, "signal", "s", "SIGTERM", "Signal to send (default SIGTERM)")

	return cmd
}

func checkCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check <port>",
		Short: "Check if a port is in use",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			port := 0
			if _, err := fmt.Sscanf(args[0], "%d", &port); err != nil {
				return fmt.Errorf("invalid port: %s", args[0])
			}

			s, err := scanner.New()
			if err != nil {
				return err
			}

			ports, err := s.Scan()
			if err != nil {
				return fmt.Errorf("scanning: %w", err)
			}

			for _, p := range ports {
				if p.Port == port {
					fmt.Printf("Port %d is in use by %q (PID %d, %s)\n", port, p.ProcessName, p.PID, p.Protocol)
					os.Exit(1)
				}
			}

			fmt.Printf("Port %d is free\n", port)
			return nil
		},
	}
}

func watchCmd() *cobra.Command {
	var (
		portFilter int
		interval   int
	)

	cmd := &cobra.Command{
		Use:   "watch",
		Short: "Watch ports with streaming output",
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := scanner.New()
			if err != nil {
				return err
			}

			ticker := time.NewTicker(time.Duration(interval) * time.Second)
			defer ticker.Stop()

			// Print immediately, then on each tick
			for {
				ports, err := s.Scan()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Scan error: %v\n", err)
				} else {
					filtered := applyFilters(ports, portFilter, "")
					fmt.Print("\033[2J\033[H") // clear screen
					fmt.Printf("PortPilot Watch — %s — %d ports\n\n",
						time.Now().Format("15:04:05"), len(filtered))
					printTable(filtered)
					fmt.Printf("\nRefreshing every %ds... Press Ctrl+C to stop.\n", interval)
				}

				<-ticker.C
			}
		},
	}

	cmd.Flags().IntVar(&portFilter, "port", 0, "Watch a specific port")
	cmd.Flags().IntVar(&interval, "interval", 2, "Refresh interval in seconds")

	return cmd
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("portpilot %s\n", version)
			fmt.Printf("  commit: %s\n", commit)
			fmt.Printf("  built:  %s\n", date)
		},
	}
}

func applyFilters(ports []scanner.PortInfo, port int, proc string) []scanner.PortInfo {
	if port == 0 && proc == "" {
		return ports
	}

	var result []scanner.PortInfo
	for _, p := range ports {
		if port != 0 && p.Port != port {
			continue
		}
		if proc != "" && !strings.Contains(strings.ToLower(p.ProcessName), strings.ToLower(proc)) {
			continue
		}
		result = append(result, p)
	}
	return result
}

func printTable(ports []scanner.PortInfo) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "PORT\tPROTO\tPID\tPROCESS\tUSER\tCPU%\tMEM%\tSTATE")
	fmt.Fprintln(w, "----\t-----\t---\t-------\t----\t----\t----\t-----")
	for _, p := range ports {
		fmt.Fprintf(w, "%d\t%s\t%d\t%s\t%s\t%.1f\t%.1f\t%s\n",
			p.Port, p.Protocol, p.PID, p.ProcessName, p.User, p.CPU, p.Mem, p.State)
	}
	w.Flush()
}

func printJSON(ports []scanner.PortInfo) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(ports)
}