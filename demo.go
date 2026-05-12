package main

import (
	"os"
	"path/filepath"
)

func demoSeedRoot() string {
	if root := os.Getenv("TUI_HUB_DEMO_ROOT"); root != "" {
		return root
	}
	return filepath.Join("/home/demo", "seed-data")
}

func isDemoMode() bool {
	return os.Getenv("TUI_HUB_DEMO") == "1" ||
		os.Getenv("DEMO_ENV") == "1" ||
		os.Getenv("DEMO_READONLY") == "1"
}

func withDemoEnv(env []string) []string {
	if !isDemoMode() {
		return env
	}
	return append(env,
		"TUI_HUB_DEMO=1",
		"DEMO_ENV=1",
		"DEMO_READONLY=1",
	)
}

func demoHomeRoot() string {
	if root := os.Getenv("TUI_HUB_DEMO_HOME"); root != "" {
		return root
	}
	return filepath.Join(os.TempDir(), "tui-hub-demo-home")
}

func demoHomeDir(appID string) string {
	if !isDemoMode() {
		return ""
	}
	if appID == "" {
		return demoHomeRoot()
	}
	return filepath.Join(demoHomeRoot(), appID)
}

func withDemoAppEnv(env []string, appID string) []string {
	env = withDemoEnv(env)
	home := demoHomeDir(appID)
	if home == "" {
		return env
	}
	_ = os.MkdirAll(filepath.Join(home, ".local", "share"), 0o755)
	_ = os.MkdirAll(filepath.Join(home, ".config"), 0o755)
	return append(env,
		"HOME="+home,
		"XDG_DATA_HOME="+filepath.Join(home, ".local", "share"),
		"XDG_CONFIG_HOME="+filepath.Join(home, ".config"),
	)
}

func demoWorkingDir(appID string) string {
	if !isDemoMode() {
		return ""
	}
	dir := filepath.Join(demoSeedRoot(), appID)
	if info, err := os.Stat(dir); err == nil && info.IsDir() {
		return dir
	}
	return ""
}

func demoLaunchArgs(appID string) []string {
	if !isDemoMode() {
		return nil
	}
	switch appID {
	case "scout":
		if dir := demoWorkingDir(appID); dir != "" {
			return []string{"--root", dir}
		}
	}
	return nil
}
