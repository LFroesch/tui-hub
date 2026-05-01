# TUI Hub

Launcher for the `tui-suite` app set. `tui-hub` keeps the main view focused on apps you already have installed, while a separate page lets you browse and install the rest of the curated suite.

## Quick Install

Supported platforms: Linux and macOS. On Windows, use WSL.

Recommended (installs to `~/.local/bin`):

```bash
curl -fsSL https://raw.githubusercontent.com/LFroesch/tui-hub/main/install.sh | bash
```

Or download a binary from [GitHub Releases](https://github.com/LFroesch/tui-hub/releases).

Or install with Go:

```bash
go install github.com/LFroesch/tui-hub@latest
```

Or build from source:

```bash
make install
```

Command:

```bash
tui-hub
tui-hub --version
```

## What It Shows

`tui-hub` ships with a built-in catalog of these suite apps:

- `runx`
- `scout`
- `portmon`
- `backup-xd`
- `seedbank`
- `zap`
- `bobdb`
- `logdog`
- `unrot`
- `sb`
- `dwight`

The launcher itself is not listed in the catalog.

## Pages

### Installed

Apps found on your `PATH`. This is the default page.

- Launch apps with `enter`
- See local version info when `<app> --version` is available
- Press `r` to manually check the latest GitHub release
- Press `u` to update an app when a newer release is found

### Available

Curated suite apps that are not currently installed.

- Press `i` to install the selected app

Install and update actions reuse each app's own `install.sh` release installer.

## Controls

- `tab`, `1`, `2` - switch between Installed and Available
- `j/k`, `up/down` - move selection
- `enter` - launch selected installed app
- `i` - install selected available app
- `u` - update selected installed app when an update is available
- `r` - manually check latest releases for installed apps
- `q` - quit

## Configuration

Config file: `~/.config/tui-hub/config.json`

The config is minimal and user-state only. Right now it stores the last active page. The app catalog, repo metadata, install behavior, and descriptions are built into `tui-hub`.

Unknown fields in an older config are ignored.

## Notes

- Version checks are manual only. `tui-hub` does not hit GitHub on startup.
- Games are intentionally out of scope for this version and can come back later.
- Future versions can grow into custom user-added app entries, but v1 stays curated and simple.

## License

[AGPL-3.0](LICENSE)
