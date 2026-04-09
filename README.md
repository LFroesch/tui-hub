# TUI Hub (Apps + Games)

## Quick Install

Recommended (installs to `~/.local/bin`):

```bash
curl -fsSL https://raw.githubusercontent.com/LFroesch/tui-hub/main/install.sh | bash
```

Or download a binary from [GitHub Releases](https://github.com/LFroesch/tui-hub/releases).

Or build from source:

```bash
make install
```

Command:

```bash
tui-hub
```
## Apps

Collection of terminal applications for productivity, development, and system utilities. Browse file managers, text editors, system monitors, and other TUI tools.

## Games

A collection of terminal-based games built with Go and BubbleTea including Chess with full rule implementation, Snake with score tracking, Blackjack with card counting, Auto-battler with strategic gameplay, and Mini ASCII Roguelike with dungeon exploration and combat systems.

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
