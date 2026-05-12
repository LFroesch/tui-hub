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
