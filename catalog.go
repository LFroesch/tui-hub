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
		{ID: "runx", Name: "runx", Binary: "runx", Repo: "LFroesch/runx", Description: "Saved script runner with schedules, prompts, and live output.", Icon: "📜", Color: "117"},
		{ID: "scout", Name: "scout", Binary: "scout", Repo: "LFroesch/scout", Description: "File explorer with preview, search, bookmarks, and shell-friendly navigation.", Icon: "🔎", Color: "81"},
		{ID: "portmon", Name: "portmon", Binary: "portmon", Repo: "LFroesch/portmon", Description: "Live port monitor with process ownership and lightweight system stats.", Icon: "📡", Color: "214"},
		{ID: "backup-xd", Name: "backup-xd", Binary: "backup-xd", Repo: "LFroesch/backup-xd", Description: "Backup manager for database dumps, file copies, and restores.", Icon: "💾", Color: "141"},
		{ID: "seedbank", Name: "seedbank", Binary: "seedbank", Repo: "LFroesch/seedbank", Description: "Fake-data generator for fixtures, demos, and seed scripts.", Icon: "🌱", Color: "78"},
		{ID: "zap", Name: "zap", Binary: "zap", Repo: "LFroesch/zap", Description: "Personal file registry for fast preview, reopen, and editing.", Icon: "⚡", Color: "220"},
		{ID: "bobdb", Name: "bobdb", Binary: "bobdb", Repo: "LFroesch/bobdb", Description: "Database browser and query runner for SQLite, Postgres, and MongoDB.", Icon: "🧮", Color: "111"},
		{ID: "logdog", Name: "logdog", Binary: "logdog", Repo: "LFroesch/logdog", Description: "Log discovery, tailing, filtering, and terminal inspection.", Icon: "🐶", Color: "203"},
		{ID: "unrot", Name: "unrot", Binary: "unrot", Repo: "LFroesch/unrot", Description: "Knowledge review and spaced-repetition study app for your own notes.", Icon: "🧠", Color: "177"},
		{ID: "sb", Name: "sb", Binary: "sb", Repo: "LFroesch/sb", Description: "WORK.md control plane for cleanup, dumps, and agent-backed runs.", Icon: "📓", Color: "149"},
		{ID: "dwight", Name: "dwight", Binary: "dwight", Repo: "LFroesch/dwight", Description: "Terminal AI chat client for Ollama, Gemini, and local file context.", Icon: "🤖", Color: "51"},
		{ID: "stickies", Name: "stickies", Binary: "stickies", Repo: "LFroesch/stickies", Description: "Quick notes and daily journaling with a small pipe-friendly CLI.", Icon: "📝", Color: "229"},
	}
}
