# DEVLOG

## 2026-04-30

- Rebuilt `tui-hub` from the old games-oriented launcher template into a curated `tui-suite` app launcher with two pages: `Installed` and `Available`.
- Replaced config-authored app definitions with a built-in catalog for `runx`, `scout`, `portmon`, `backup-xd`, `seedbank`, `zap`, `bobdb`, `logdog`, `unrot`, `sb`, and `dwight`.
- Removed hardcoded local repo launch paths. Launch now resolves installed binaries from `PATH`, local version info comes from `<app> --version`, and remote release checks are manual only.
- Added in-app install and update actions that reuse each target app's `install.sh` flow instead of inventing a second installer path.
- Split the code into smaller files for bootstrap, catalog, config, actions, model, and view logic, and added small tests around version parsing and page filtering.
- Adjusted startup so installed/available catalog population does not depend on version probing succeeding; local version scans now run separately with short timeouts.
- Reworked the main list from tall per-app cards into a denser table with capped visible height, row scrolling, and page-style movement keys so it fits the rest of the suite better.
- Restored per-app visual identity in the new compact table by adding built-in emoji and accent colors for each suite app instead of leaving every row visually identical.
- Added persisted per-app launch count and last-launched config state, and now sort Installed apps by a frecency score while leaving Available alphabetical.
