package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type config struct {
	LastPage string `json:"last_page"`
}

func defaultConfig() config {
	return config{LastPage: pageInstalled}
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
	if cfg.LastPage != pageInstalled && cfg.LastPage != pageAvailable {
		cfg.LastPage = pageInstalled
	}
	return cfg
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
