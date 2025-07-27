# TUI Hub (Apps + Games)

## Apps

Collection of terminal applications for productivity, development, and system utilities. Browse file managers, text editors, system monitors, and other TUI tools.

## Games

Terminal games including puzzles, arcade-style games, and interactive entertainment. All playable directly in your terminal.
## Installation

```bash
git clone <repository-url>
cd tui-hub
go build -o tui-hub main.go
./tui-hub
```

Or move to ~/.local/bin

## Controls

- **↑/↓**: Navigate
- **Enter**: Launch
- **←/→**: Switch menus  
- **q**: Quit

## Configuration

Config file: `~/.config/tui-hub/config.json`

Auto-created on first run. Add apps/games in JSON format with name, description, command, and path.