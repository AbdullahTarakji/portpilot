package tui

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/AbdullahTarakji/portpilot/internal/config"
	"github.com/AbdullahTarakji/portpilot/internal/scanner"
)

// mockScanner returns canned port data for testing.
type mockScanner struct {
	ports []scanner.PortInfo
	err   error
}

func (m *mockScanner) Scan() ([]scanner.PortInfo, error) {
	return m.ports, m.err
}

func testPorts() []scanner.PortInfo {
	return []scanner.PortInfo{
		{Port: 3000, Protocol: "TCP", PID: 100, ProcessName: "node", User: "mike", State: "LISTEN", CPU: 2.1, Mem: 1.3},
		{Port: 5432, Protocol: "TCP", PID: 200, ProcessName: "postgres", User: "mike", State: "LISTEN", CPU: 0.1, Mem: 0.5},
		{Port: 6379, Protocol: "TCP", PID: 300, ProcessName: "redis-ser", User: "mike", State: "LISTEN", CPU: 0.0, Mem: 0.1},
		{Port: 8080, Protocol: "TCP", PID: 400, ProcessName: "Python", User: "mike", State: "LISTEN", CPU: 55.0, Mem: 12.0},
	}
}

func newTestModel() Model {
	s := &mockScanner{ports: testPorts()}
	cfg := config.DefaultConfig()
	m := New(s, cfg)
	m.ports = testPorts()
	m.width = 120
	m.height = 40
	m.lastRefresh = time.Now()
	return m
}

func TestNewModelInitialization(t *testing.T) {
	s := &mockScanner{ports: testPorts()}
	cfg := config.DefaultConfig()
	m := New(s, cfg)

	if m.cursor != 0 {
		t.Errorf("cursor: got %d, want 0", m.cursor)
	}
	if m.filterMode {
		t.Error("filterMode should be false initially")
	}
	if m.view != viewTable {
		t.Errorf("view: got %d, want viewTable(%d)", m.view, viewTable)
	}
	if m.showGroups {
		t.Error("showGroups should be false initially")
	}
	if m.sortCol.column != 0 {
		t.Errorf("sortCol: got %d, want 0", m.sortCol.column)
	}
	if !m.sortCol.asc {
		t.Error("sort should be ascending initially")
	}
}

func TestQuitOnQ(t *testing.T) {
	m := newTestModel()
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	_ = updated

	if cmd == nil {
		t.Fatal("expected quit command, got nil")
	}
	// Execute the command and check it returns tea.QuitMsg
	msg := cmd()
	if _, ok := msg.(tea.QuitMsg); !ok {
		t.Errorf("expected tea.QuitMsg, got %T", msg)
	}
}

func TestQuitOnCtrlC(t *testing.T) {
	m := newTestModel()
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})

	if cmd == nil {
		t.Fatal("expected quit command on Ctrl+C, got nil")
	}
	msg := cmd()
	if _, ok := msg.(tea.QuitMsg); !ok {
		t.Errorf("expected tea.QuitMsg, got %T", msg)
	}
}

func TestFilterModeToggle(t *testing.T) {
	m := newTestModel()

	// Press '/' to enter filter mode
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("/")})
	m = updated.(Model)

	if !m.filterMode {
		t.Error("expected filterMode=true after pressing '/'")
	}

	// Press Esc to exit filter mode
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = updated.(Model)

	if m.filterMode {
		t.Error("expected filterMode=false after pressing Esc")
	}
}

func TestFilterTextInput(t *testing.T) {
	m := newTestModel()

	// Enter filter mode
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("/")})
	m = updated.(Model)

	// Type "node"
	for _, ch := range "node" {
		updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
		m = updated.(Model)
	}

	if m.filter != "node" {
		t.Errorf("filter: got %q, want %q", m.filter, "node")
	}

	// Backspace
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	m = updated.(Model)

	if m.filter != "nod" {
		t.Errorf("filter after backspace: got %q, want %q", m.filter, "nod")
	}
}

func TestHelpOverlayToggle(t *testing.T) {
	m := newTestModel()

	// Press '?' to open help
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("?")})
	m = updated.(Model)

	if m.view != viewHelp {
		t.Errorf("view: got %d, want viewHelp(%d)", m.view, viewHelp)
	}

	// Press '?' again to close
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("?")})
	m = updated.(Model)

	if m.view != viewTable {
		t.Errorf("view: got %d, want viewTable(%d)", m.view, viewTable)
	}
}

func TestCursorNavigation(t *testing.T) {
	m := newTestModel()

	// Move down
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updated.(Model)
	if m.cursor != 1 {
		t.Errorf("cursor after down: got %d, want 1", m.cursor)
	}

	// Move down again
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updated.(Model)
	if m.cursor != 2 {
		t.Errorf("cursor after 2x down: got %d, want 2", m.cursor)
	}

	// Move up
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = updated.(Model)
	if m.cursor != 1 {
		t.Errorf("cursor after up: got %d, want 1", m.cursor)
	}

	// Don't go below 0
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = updated.(Model)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = updated.(Model)
	if m.cursor != 0 {
		t.Errorf("cursor should not go below 0: got %d", m.cursor)
	}
}

func TestGroupToggle(t *testing.T) {
	m := newTestModel()

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("g")})
	m = updated.(Model)

	if !m.showGroups {
		t.Error("expected showGroups=true after pressing 'g'")
	}

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("g")})
	m = updated.(Model)

	if m.showGroups {
		t.Error("expected showGroups=false after pressing 'g' again")
	}
}

func TestDetailViewToggle(t *testing.T) {
	m := newTestModel()

	// Press Enter to open detail
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(Model)

	if m.view != viewDetail {
		t.Errorf("view: got %d, want viewDetail(%d)", m.view, viewDetail)
	}

	// Press Esc to go back
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = updated.(Model)

	if m.view != viewTable {
		t.Errorf("view: got %d, want viewTable(%d)", m.view, viewTable)
	}
}

func TestConfirmKillView(t *testing.T) {
	m := newTestModel()

	// Press 'k' to open kill confirmation
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
	m = updated.(Model)

	if m.view != viewConfirmKill {
		t.Errorf("view: got %d, want viewConfirmKill(%d)", m.view, viewConfirmKill)
	}

	// Press 'n' to cancel
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("n")})
	m = updated.(Model)

	if m.view != viewTable {
		t.Errorf("view after cancel: got %d, want viewTable(%d)", m.view, viewTable)
	}
	if m.statusMsg != "Kill cancelled" {
		t.Errorf("statusMsg: got %q, want %q", m.statusMsg, "Kill cancelled")
	}
}

func TestSortColumnToggle(t *testing.T) {
	m := newTestModel()

	// Press '1' to sort by first column
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("1")})
	m = updated.(Model)

	if m.sortCol.column != 0 {
		t.Errorf("sortCol: got %d, want 0", m.sortCol.column)
	}
	// It was already col 0 asc, so should flip to desc
	if m.sortCol.asc {
		t.Error("expected sort desc after pressing same column")
	}

	// Press '3' to sort by third column
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("3")})
	m = updated.(Model)

	if m.sortCol.column != 2 {
		t.Errorf("sortCol: got %d, want 2", m.sortCol.column)
	}
	if !m.sortCol.asc {
		t.Error("expected sort asc for new column")
	}
}

func TestWindowSizeMsg(t *testing.T) {
	m := newTestModel()

	updated, _ := m.Update(tea.WindowSizeMsg{Width: 200, Height: 50})
	m = updated.(Model)

	if m.width != 200 {
		t.Errorf("width: got %d, want 200", m.width)
	}
	if m.height != 50 {
		t.Errorf("height: got %d, want 50", m.height)
	}
}

func TestScanResultMsg(t *testing.T) {
	m := newTestModel()
	m.ports = nil

	newPorts := testPorts()[:2]
	updated, _ := m.Update(scanResultMsg{ports: newPorts})
	m = updated.(Model)

	if len(m.ports) != 2 {
		t.Errorf("ports: got %d, want 2", len(m.ports))
	}
	if m.lastRefresh.IsZero() {
		t.Error("lastRefresh should be set after scan")
	}
}

func TestViewRendersWithoutPanic(t *testing.T) {
	m := newTestModel()

	// Table view
	output := m.View()
	if output == "" {
		t.Error("View() returned empty string for table view")
	}

	// Help view
	m.view = viewHelp
	output = m.View()
	if output == "" {
		t.Error("View() returned empty string for help view")
	}

	// Detail view
	m.view = viewDetail
	output = m.View()
	if output == "" {
		t.Error("View() returned empty string for detail view")
	}

	// Confirm kill view
	m.view = viewConfirmKill
	output = m.View()
	if output == "" {
		t.Error("View() returned empty string for confirm kill view")
	}
}

func TestViewRenderWithFilter(t *testing.T) {
	m := newTestModel()
	m.filter = "node"
	m.filterMode = true

	output := m.View()
	if output == "" {
		t.Error("View() returned empty string with active filter")
	}
}

func TestEmptyPortsList(t *testing.T) {
	s := &mockScanner{ports: nil}
	cfg := config.DefaultConfig()
	m := New(s, cfg)
	m.width = 120
	m.height = 40
	m.lastRefresh = time.Now()

	// Should not panic on empty ports
	output := m.View()
	if output == "" {
		t.Error("View() returned empty string with no ports")
	}

	// Navigation on empty list should be safe
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updated.(Model)
	if m.cursor != 0 {
		t.Errorf("cursor should stay at 0 with empty list, got %d", m.cursor)
	}

	// Kill on empty should not crash
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
	m = updated.(Model)
	if m.view != viewTable {
		t.Error("kill on empty list should stay on table view")
	}
}
