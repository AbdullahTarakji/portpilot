// Package config handles loading and parsing of the portpilot configuration file.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the portpilot configuration.
type Config struct {
	Groups          map[string]Group `yaml:"groups"`
	RefreshInterval int              `yaml:"refresh_interval"`
	ShowSystemPorts bool             `yaml:"show_system_ports"`
}

// Group defines a named port group with associated color.
type Group struct {
	Ports []int  `yaml:"ports"`
	Color string `yaml:"color"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Groups:          make(map[string]Group),
		RefreshInterval: 2,
		ShowSystemPorts: false,
	}
}

// Load reads the config from ~/.portpilot.yaml.
// Returns default config if the file doesn't exist.
func Load() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return DefaultConfig(), nil
	}

	return LoadFrom(filepath.Join(home, ".portpilot.yaml"))
}

// LoadFrom reads config from a specific file path.
// Returns default config if the file doesn't exist.
func LoadFrom(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return nil, fmt.Errorf("reading config %s: %w", path, err)
	}

	return Parse(data)
}

// Parse parses config from YAML bytes.
func Parse(data []byte) (*Config, error) {
	cfg := DefaultConfig()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	if cfg.RefreshInterval < 1 {
		cfg.RefreshInterval = 2
	}

	if cfg.Groups == nil {
		cfg.Groups = make(map[string]Group)
	}

	return cfg, nil
}

// GroupForPort returns the group name for a port, or empty string if ungrouped.
func (c *Config) GroupForPort(port int) string {
	for name, g := range c.Groups {
		for _, p := range g.Ports {
			if p == port {
				return name
			}
		}
	}
	return ""
}

// GroupColor returns the color for a group name, or empty string.
func (c *Config) GroupColor(name string) string {
	if g, ok := c.Groups[name]; ok {
		return g.Color
	}
	return ""
}
