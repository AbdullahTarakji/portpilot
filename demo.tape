Output demo.gif

Set FontSize 14
Set Width 1200
Set Height 600
Set Theme "Catppuccin Mocha"
Set TypingSpeed 50ms
Set Padding 20

# Title
Type "# 🚀 PortPilot — See what's running on your ports"
Enter
Sleep 1s
Ctrl+L

# Show the list command
Type "portpilot list"
Enter
Sleep 3s

# Show filtering by process
Type "portpilot list --process node"
Enter
Sleep 2s

# Check a specific port
Type "portpilot check 5432"
Enter
Sleep 2s

# Check a free port
Type "portpilot check 9999"
Enter
Sleep 2s

# JSON output
Type "portpilot list --json | head -20"
Enter
Sleep 3s

# Launch TUI
Type "portpilot"
Enter
Sleep 3s

# Search filter
Type "/"
Sleep 500ms
Type "node"
Sleep 2s

# Exit search
Escape
Sleep 1s

# Quit
Type "q"
Sleep 1s
