# VHS Demo Tape for MLB CLI
# Run with: vhs demo.tape

Output demo.gif

Set FontSize 14
Set Width 1200
Set Height 600
Set Theme "Catppuccin Mocha"

Type "# MLB CLI - Interactive Terminal UI"
Enter
Sleep 1s

Type "./mlb"
Enter
Sleep 3s

# Navigate down through teams
Down
Sleep 500ms
Down
Sleep 500ms
Down
Sleep 500ms

# Filter teams
Type "/"
Sleep 500ms
Type "dodgers"
Sleep 1s
Enter
Sleep 500ms

# Clear filter and select team
Escape
Sleep 500ms

# Select a team (Dodgers)
Enter
Sleep 2s

# Navigate roster
Down
Sleep 300ms
Down
Sleep 300ms
Down
Sleep 300ms

# Select a player to see stats
Enter
Sleep 3s

# Go back
Backspace
Sleep 1s

# Switch to Standings tab
Tab
Sleep 2s

# Switch to Schedule tab
Tab
Sleep 2s

# Quit
Type "q"
Sleep 1s
