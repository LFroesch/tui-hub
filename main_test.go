package main

import (
	"encoding/json"
	"testing"
	"time"
)

func TestExtractVersion(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "plain", in: "runx 1.2.3", want: "1.2.3"},
		{name: "prefixed", in: "sb version v0.9.0", want: "v0.9.0"},
		{name: "missing", in: "development build", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractVersion(tt.in); got != tt.want {
				t.Fatalf("extractVersion(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestShouldOfferUpdate(t *testing.T) {
	tests := []struct {
		name   string
		local  string
		latest string
		want   bool
	}{
		{name: "same", local: "1.2.3", latest: "v1.2.3", want: false},
		{name: "dirty local build", local: "v0.9.1-dirty", latest: "v0.9.1", want: false},
		{name: "different", local: "1.2.3", latest: "v1.2.4", want: true},
		{name: "missing local", local: "", latest: "v1.2.4", want: true},
		{name: "missing latest", local: "1.2.3", latest: "", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldOfferUpdate(tt.local, tt.latest); got != tt.want {
				t.Fatalf("shouldOfferUpdate(%q, %q) = %v, want %v", tt.local, tt.latest, got, tt.want)
			}
		})
	}
}

func TestNormalizeVersion(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "trim v prefix", in: "v1.2.3", want: "1.2.3"},
		{name: "trim dirty suffix", in: "v0.9.1-dirty", want: "0.9.1"},
		{name: "keep prerelease", in: "v1.2.3-rc1", want: "1.2.3-rc1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeVersion(tt.in); got != tt.want {
				t.Fatalf("normalizeVersion(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestAppsForPage(t *testing.T) {
	m := model{
		apps: []suiteApp{
			{appCatalogEntry: appCatalogEntry{Name: "runx"}, Installed: true},
			{appCatalogEntry: appCatalogEntry{Name: "scout"}, Installed: false},
			{appCatalogEntry: appCatalogEntry{Name: "sb"}, Installed: true},
		},
	}

	if got := len(m.appsForPage(pageInstalled)); got != 2 {
		t.Fatalf("installed count = %d, want 2", got)
	}
	if got := len(m.appsForPage(pageAvailable)); got != 1 {
		t.Fatalf("available count = %d, want 1", got)
	}
}

func TestRefreshAppsIncludesBuiltInCatalog(t *testing.T) {
	apps := refreshApps("test")
	if got, want := len(apps), len(builtInCatalog()); got != want {
		t.Fatalf("refreshApps count = %d, want %d", got, want)
	}
}

func TestUpdateAppliesAppsRefreshedMsg(t *testing.T) {
	m := initialModel("dev")
	msg := appsRefreshedMsg{apps: refreshAppsWithConfig(m.cfg), status: "Ready"}
	next, cmd := m.Update(msg)
	if cmd == nil {
		t.Fatalf("expected follow-up local version scan cmd")
	}
	got := next.(model)
	if len(got.apps) != len(builtInCatalog()) {
		t.Fatalf("model apps count = %d, want %d", len(got.apps), len(builtInCatalog()))
	}
	if got.status != "Ready" {
		t.Fatalf("status = %q, want Ready", got.status)
	}
}

func TestNormalizeConfigRoundTripAppState(t *testing.T) {
	cfg := config{
		LastPage: pageAvailable,
		AppState: map[string]appState{
			"runx": {
				LaunchCount:  3,
				LastLaunched: "2026-04-30T12:00:00Z",
			},
		},
	}

	data, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("marshal config: %v", err)
	}

	var decoded config
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal config: %v", err)
	}

	decoded = normalizeConfig(decoded)
	state := decoded.AppState["runx"]
	if decoded.LastPage != pageAvailable {
		t.Fatalf("last page = %q, want %q", decoded.LastPage, pageAvailable)
	}
	if state.LaunchCount != 3 {
		t.Fatalf("launch count = %d, want 3", state.LaunchCount)
	}
	if state.LastLaunched != "2026-04-30T12:00:00Z" {
		t.Fatalf("last launched = %q", state.LastLaunched)
	}
}

func TestSortAppsInstalledUsesFrecency(t *testing.T) {
	now := time.Now().UTC()
	apps := []suiteApp{
		{appCatalogEntry: appCatalogEntry{ID: "runx", Name: "runx"}, Installed: true},
		{appCatalogEntry: appCatalogEntry{ID: "scout", Name: "scout"}, Installed: true},
		{appCatalogEntry: appCatalogEntry{ID: "bobdb", Name: "bobdb"}, Installed: true},
		{appCatalogEntry: appCatalogEntry{ID: "zap", Name: "zap"}, Installed: false},
		{appCatalogEntry: appCatalogEntry{ID: "dwight", Name: "dwight"}, Installed: false},
	}
	cfg := config{
		LastPage: pageInstalled,
		AppState: map[string]appState{
			"scout": {LaunchCount: 10, LastLaunched: now.Add(-72 * time.Hour).Format(time.RFC3339)},
			"bobdb": {LaunchCount: 3, LastLaunched: now.Add(-1 * time.Hour).Format(time.RFC3339)},
			"runx":  {LaunchCount: 1, LastLaunched: now.Add(-10 * time.Minute).Format(time.RFC3339)},
		},
	}

	sortApps(apps, cfg)

	if apps[0].ID != "scout" || apps[1].ID != "bobdb" || apps[2].ID != "runx" {
		t.Fatalf("installed order = %q, %q, %q", apps[0].ID, apps[1].ID, apps[2].ID)
	}
	if apps[3].ID != "dwight" || apps[4].ID != "zap" {
		t.Fatalf("available order = %q, %q", apps[3].ID, apps[4].ID)
	}
}

func TestRecordLaunchUpdatesConfigState(t *testing.T) {
	when := time.Date(2026, 4, 30, 14, 0, 0, 0, time.UTC)
	m := model{
		cfg: config{LastPage: pageInstalled, AppState: map[string]appState{}},
		apps: []suiteApp{
			{appCatalogEntry: appCatalogEntry{ID: "runx", Name: "runx"}, Installed: true},
		},
	}

	app, ok := m.recordLaunch("runx", when)
	if !ok {
		t.Fatalf("expected app lookup to succeed")
	}
	if app.Name != "runx" {
		t.Fatalf("app name = %q, want runx", app.Name)
	}
	state := m.cfg.AppState["runx"]
	if state.LaunchCount != 1 {
		t.Fatalf("launch count = %d, want 1", state.LaunchCount)
	}
	if state.LastLaunched != when.Format(time.RFC3339) {
		t.Fatalf("last launched = %q, want %q", state.LastLaunched, when.Format(time.RFC3339))
	}
}
