package main

import "testing"

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
