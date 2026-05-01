package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type config struct {
	LastPage string              `json:"last_page"`
	AppState map[string]appState `json:"app_state,omitempty"`
}

type appState struct {
	LaunchCount  int    `json:"launch_count"`
	LastLaunched string `json:"last_launched,omitempty"`
}

func defaultConfig() config {
	return config{
		LastPage: pageInstalled,
		AppState: map[string]appState{},
	}
}

func configPath() string {
	base, err := os.UserConfigDir()
	if err != nil {
		home, homeErr := os.UserHomeDir()
		if homeErr != nil {
			return filepath.Join(".config", "tui-hub", "config.json")
		}
		base = filepath.Join(home, ".config")
	}
	return filepath.Join(base, "tui-hub", "config.json")
}

func loadConfig() config {
	cfg := defaultConfig()
	data, err := os.ReadFile(configPath())
	if err != nil {
		return cfg
	}
	_ = json.Unmarshal(data, &cfg)
	return normalizeConfig(cfg)
}

func saveConfig(cfg config) error {
	path := configPath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func normalizeConfig(cfg config) config {
	if cfg.LastPage != pageInstalled && cfg.LastPage != pageAvailable {
		cfg.LastPage = pageInstalled
	}
	if cfg.AppState == nil {
		cfg.AppState = map[string]appState{}
	}
	for appID, state := range cfg.AppState {
		if state.LaunchCount < 0 {
			state.LaunchCount = 0
		}
		if state.LastLaunched != "" {
			if _, err := time.Parse(time.RFC3339, state.LastLaunched); err != nil {
				state.LastLaunched = ""
			}
		}
		cfg.AppState[appID] = state
	}
	return cfg
}
