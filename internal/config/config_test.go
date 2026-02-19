package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParse(t *testing.T) {
	data := []byte(`
groups:
  frontend:
    ports: [3000, 3001, 5173]
    color: blue
  backend:
    ports: [4000, 8000]
    color: green
refresh_interval: 5
show_system_ports: true
`)

	cfg, err := Parse(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.RefreshInterval != 5 {
		t.Errorf("refresh_interval: got %d, want 5", cfg.RefreshInterval)
	}
	if !cfg.ShowSystemPorts {
		t.Error("show_system_ports: got false, want true")
	}
	if len(cfg.Groups) != 2 {
		t.Fatalf("groups: got %d, want 2", len(cfg.Groups))
	}
	fe := cfg.Groups["frontend"]
	if len(fe.Ports) != 3 {
		t.Errorf("frontend ports: got %d, want 3", len(fe.Ports))
	}
	if fe.Color != "blue" {
		t.Errorf("frontend color: got %s, want blue", fe.Color)
	}
}

func TestParseEmpty(t *testing.T) {
	cfg, err := Parse([]byte(""))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.RefreshInterval != 2 {
		t.Errorf("default refresh_interval: got %d, want 2", cfg.RefreshInterval)
	}
}

func TestParseInvalidRefreshInterval(t *testing.T) {
	cfg, err := Parse([]byte("refresh_interval: 0"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.RefreshInterval != 2 {
		t.Errorf("should default to 2, got %d", cfg.RefreshInterval)
	}
}

func TestGroupForPort(t *testing.T) {
	cfg, _ := Parse([]byte(`
groups:
  db:
    ports: [5432, 3306]
    color: yellow
`))

	if g := cfg.GroupForPort(5432); g != "db" {
		t.Errorf("GroupForPort(5432): got %q, want %q", g, "db")
	}
	if g := cfg.GroupForPort(9999); g != "" {
		t.Errorf("GroupForPort(9999): got %q, want empty", g)
	}
}

func TestGroupColor(t *testing.T) {
	cfg, _ := Parse([]byte(`
groups:
  db:
    ports: [5432]
    color: yellow
`))

	if c := cfg.GroupColor("db"); c != "yellow" {
		t.Errorf("GroupColor(db): got %q, want yellow", c)
	}
	if c := cfg.GroupColor("nope"); c != "" {
		t.Errorf("GroupColor(nope): got %q, want empty", c)
	}
}

func TestLoadFromMissing(t *testing.T) {
	cfg, err := LoadFrom("/nonexistent/path/config.yaml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.RefreshInterval != 2 {
		t.Errorf("expected default config, got refresh_interval=%d", cfg.RefreshInterval)
	}
}

func TestLoadFromFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	os.WriteFile(path, []byte("refresh_interval: 10\n"), 0644)

	cfg, err := LoadFrom(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.RefreshInterval != 10 {
		t.Errorf("refresh_interval: got %d, want 10", cfg.RefreshInterval)
	}
}
