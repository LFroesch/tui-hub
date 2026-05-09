# tui-hub

Launcher for the `tui-suite` app collection. `tui-hub` keeps the default view focused on what you already have installed, and gives you a second page for browsing the curated install catalog.

## Install

Supported platforms: Linux and macOS. On Windows, use WSL.

Recommended:

```bash
curl -fsSL https://raw.githubusercontent.com/LFroesch/tui-hub/main/install.sh | bash
```

Other options:

```bash
go install github.com/LFroesch/tui-hub@latest
make install
```

Run:

```bash
tui-hub
tui-hub --version
```

## Pages

| Page | Purpose |
|------|---------|
| Installed | Launch apps already on your `PATH`, sorted by frecency |
| Available | Browse curated suite apps that are not installed yet |

Install and update actions call each app's own `install.sh` release installer.

## Projects In This Suite

| App | Description | In `tui-hub` catalog |
|-----|-------------|----------------------|
| [backup-xd](https://github.com/LFroesch/backup-xd) | Backup manager for databases and filesystem targets. | Yes |
| [bobdb](https://github.com/LFroesch/bobdb) | Database browser and query runner for SQLite, Postgres, and MongoDB. | Yes |
| [dwight](https://github.com/LFroesch/dwight) | Terminal AI chat client for Ollama and Gemini. | Yes |
| [logdog](https://github.com/LFroesch/logdog) | Log discovery, inspection, live tailing, and CLI filtering. | Yes |
| [portmon](https://github.com/LFroesch/portmon) | Port monitor and lightweight system dashboard. | Yes |
| [runx](https://github.com/LFroesch/runx) | Script runner with schedules, live output, and history. | Yes |
| [sb](https://github.com/LFroesch/sb) | `WORK.md` control plane for task cleanup, dumps, and agents. | Yes |
| [scout](https://github.com/LFroesch/scout) | File explorer with preview, search, and bookmarks. | Yes |
| [seedbank](https://github.com/LFroesch/seedbank) | Fake-data generator for fixtures and demos. | Yes |
| [stickies](https://github.com/LFroesch/stickies) | Quick notes and daily journal with a small CLI. | Yes |
| [unrot](https://github.com/LFroesch/unrot) | Knowledge review and spaced-repetition study tool. | Yes |
| [zap](https://github.com/LFroesch/zap) | Personal file registry for fast reopening and editing. | Yes |

`chunes` is out of the current `tui-hub` v1 scope.

## Controls

| Key | Action |
|-----|--------|
| `tab`, `shift+tab`, `1`, `2` | Switch pages |
| `j/k`, `up/down` | Move |
| `enter` | Launch selected installed app |
| `i` | Install selected available app |
| `u` | Update selected installed app |
| `r` | Refresh release info for installed apps |
| `g`, `G` | Jump to top or bottom |
| `?` | Help |
| `q` | Quit |

## Config

Config is stored at `~/.config/tui-hub/config.json`.

It only keeps user state:

- last active page
- per-app launch count
- last-launched timestamps used for frecency sorting

The catalog itself is built into the app.

## Notes

- version checks are manual; `tui-hub` does not call GitHub on startup
- local builds with a `-dirty` suffix are treated as the matching release version for update checks
- games are out of scope for the current launcher

## License

[AGPL-3.0](LICENSE)
