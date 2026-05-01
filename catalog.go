package main

type appCatalogEntry struct {
	ID          string
	Name        string
	Binary      string
	Repo        string
	Description string
}

func builtInCatalog() []appCatalogEntry {
	return []appCatalogEntry{
		{ID: "runx", Name: "runx", Binary: "runx", Repo: "LFroesch/runx", Description: "Run and schedule scripts from one TUI."},
		{ID: "scout", Name: "scout", Binary: "scout", Repo: "LFroesch/scout", Description: "File explorer with preview, search, and bookmarks."},
		{ID: "portmon", Name: "portmon", Binary: "portmon", Repo: "LFroesch/portmon", Description: "Monitor ports and system stats from the terminal."},
		{ID: "backup-xd", Name: "backup-xd", Binary: "backup-xd", Repo: "LFroesch/backup-xd", Description: "Manage local backup jobs and restores."},
		{ID: "seedbank", Name: "seedbank", Binary: "seedbank", Repo: "LFroesch/seedbank", Description: "Generate fake data and export fixtures."},
		{ID: "zap", Name: "zap", Binary: "zap", Repo: "LFroesch/zap", Description: "Fast terminal note and text workflow."},
		{ID: "bobdb", Name: "bobdb", Binary: "bobdb", Repo: "LFroesch/bobdb", Description: "Browse databases and run queries from a TUI."},
		{ID: "logdog", Name: "logdog", Binary: "logdog", Repo: "LFroesch/logdog", Description: "Tail, filter, and inspect logs."},
		{ID: "unrot", Name: "unrot", Binary: "unrot", Repo: "LFroesch/unrot", Description: "Knowledge review and spaced repetition from the terminal."},
		{ID: "sb", Name: "sb", Binary: "sb", Repo: "LFroesch/sb", Description: "Second-brain control plane for WORK.md projects."},
		{ID: "dwight", Name: "dwight", Binary: "dwight", Repo: "LFroesch/dwight", Description: "Chat and coding assistant terminal app."},
	}
}
