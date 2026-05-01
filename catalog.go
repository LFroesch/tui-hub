package main

type appCatalogEntry struct {
	ID          string
	Name        string
	Binary      string
	Repo        string
	Description string
	Icon        string
	Color       string
}

func builtInCatalog() []appCatalogEntry {
	return []appCatalogEntry{
		{ID: "runx", Name: "runx", Binary: "runx", Repo: "LFroesch/runx", Description: "Run and schedule scripts from one TUI.", Icon: "▶", Color: "117"},
		{ID: "scout", Name: "scout", Binary: "scout", Repo: "LFroesch/scout", Description: "File explorer with preview, search, and bookmarks.", Icon: "🔎", Color: "81"},
		{ID: "portmon", Name: "portmon", Binary: "portmon", Repo: "LFroesch/portmon", Description: "Monitor ports and system stats from the terminal.", Icon: "📡", Color: "214"},
		{ID: "backup-xd", Name: "backup-xd", Binary: "backup-xd", Repo: "LFroesch/backup-xd", Description: "Manage local backup jobs and restores.", Icon: "💾", Color: "141"},
		{ID: "seedbank", Name: "seedbank", Binary: "seedbank", Repo: "LFroesch/seedbank", Description: "Generate fake data and export fixtures.", Icon: "🌱", Color: "78"},
		{ID: "zap", Name: "zap", Binary: "zap", Repo: "LFroesch/zap", Description: "Fast terminal note and text workflow.", Icon: "⚡", Color: "220"},
		{ID: "bobdb", Name: "bobdb", Binary: "bobdb", Repo: "LFroesch/bobdb", Description: "Browse databases and run queries from a TUI.", Icon: "🗄", Color: "111"},
		{ID: "logdog", Name: "logdog", Binary: "logdog", Repo: "LFroesch/logdog", Description: "Tail, filter, and inspect logs.", Icon: "🐶", Color: "203"},
		{ID: "unrot", Name: "unrot", Binary: "unrot", Repo: "LFroesch/unrot", Description: "Knowledge review and spaced repetition from the terminal.", Icon: "🧠", Color: "177"},
		{ID: "sb", Name: "sb", Binary: "sb", Repo: "LFroesch/sb", Description: "Second-brain control plane for WORK.md projects.", Icon: "📓", Color: "149"},
		{ID: "dwight", Name: "dwight", Binary: "dwight", Repo: "LFroesch/dwight", Description: "Chat and coding assistant terminal app.", Icon: "🤖", Color: "51"},
		{ID: "stickies", Name: "stickies", Binary: "stickies", Repo: "LFroesch/stickies", Description: "Quick notes and a daily journal in one TUI.", Icon: "📝", Color: "229"},
	}
}
