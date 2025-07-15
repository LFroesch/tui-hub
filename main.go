package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	// Color scheme matching your other TUI apps
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#7C3AED")).
			Padding(0, 1)

	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#7C3AED")).
			Padding(0, 1)

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 1)

	commandStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#60A5FA")).
			Bold(true)

	keyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#34D399")).
			Bold(true)

	descStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9CA3AF"))
)

func main() {
	// Initialize the launcher
	launcher := NewLauncher()

	// Create the Bubble Tea program
	p := tea.NewProgram(
		launcher,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	// Run the program
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type Launcher struct {
	games          []GameEntry
	tuiApps        []GameEntry
	selectedGame   int
	selectedTUIApp int
	selectedMenu   int
	menuState      MenuState
	width          int
	height         int
	stats          GlobalStats
}

type MenuState int

const (
	MainMenu MenuState = iota
	TUIAppsMenu
	StatsMenu
	OptionsMenu
	CreditsMenu
)

type GameEntry struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Icon        string     `json:"icon"`
	Command     string     `json:"command"`
	Path        string     `json:"path"`
	Category    string     `json:"category"`
	Difficulty  string     `json:"difficulty"`
	Version     string     `json:"version"`
	Author      string     `json:"author"`
	Executable  bool       `json:"executable"`
	Config      GameConfig `json:"config"`
	HighScore   string     `json:"high_score,omitempty"`
	Stats       GameStats  `json:"stats,omitempty"`
}

type GameConfig struct {
	// Common config fields - exact structure depends on game
	Data map[string]interface{} `json:"-"`
}

type GameStats struct {
	TimesPlayed int    `json:"times_played"`
	TotalTime   string `json:"total_time"`
	HighScore   string `json:"high_score"`
	LastPlayed  string `json:"last_played"`
}

type Config struct {
	Launcher LauncherConfig `json:"launcher"`
	Games    []GameEntry    `json:"games"`
	TUIApps  []GameEntry    `json:"tui_apps"`
	Stats    StatsConfig    `json:"stats"`
	Settings SettingsConfig `json:"settings"`
}

type LauncherConfig struct {
	Title   string `json:"title"`
	Version string `json:"version"`
	Author  string `json:"author"`
	Theme   string `json:"theme"`
}

type StatsConfig struct {
	Global       GlobalStats   `json:"global"`
	Achievements []Achievement `json:"achievements"`
}

type GlobalStats struct {
	GamesPlayed          int    `json:"games_played"`
	TotalTimeSeconds     int    `json:"total_time_seconds"`
	AchievementsUnlocked int    `json:"achievements_unlocked"`
	FavoriteGame         string `json:"favorite_game"`
	LastPlayed           string `json:"last_played"`
}

type Achievement struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Unlocked    bool   `json:"unlocked"`
	UnlockDate  string `json:"unlock_date,omitempty"`
}

type SettingsConfig struct {
	Theme              string `json:"theme"`
	SoundEnabled       bool   `json:"sound_enabled"`
	AutoSave           bool   `json:"auto_save"`
	StatisticsTracking bool   `json:"statistics_tracking"`
	Notifications      bool   `json:"notifications"`
	BackupSaves        bool   `json:"backup_saves"`
	ControllerSupport  bool   `json:"controller_support"`
	TerminalSize       string `json:"terminal_size"`
}

func NewLauncher() *Launcher {
	config := loadConfig()
	launcher := &Launcher{
		games:          config.Games,
		tuiApps:        config.TUIApps,
		selectedGame:   0,
		selectedTUIApp: 0,
		selectedMenu:   0,
		menuState:      TUIAppsMenu,
		stats:          config.Stats.Global,
	}
	return launcher
}

func loadConfig() Config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Error finding home directory: %v", err)
		return getDefaultConfig()
	}
	configPath := filepath.Join(homeDir, ".config", "tui-hub", "config.json")

	// Try to open config file
	file, err := os.Open(configPath)
	if err != nil {
		// If not found, create default config and write it
		defaultConfig := getDefaultConfig()
		_ = os.MkdirAll(fmt.Sprintf("%s/.config/tui-hub", homeDir), 0755)
		out, err := os.Create(configPath)
		if err != nil {
			log.Printf("Error creating config.json: %v", err)
			return defaultConfig
		}
		defer out.Close()
		enc := json.NewEncoder(out)
		enc.SetIndent("", "  ")
		if err := enc.Encode(defaultConfig); err != nil {
			log.Printf("Error writing default config.json: %v", err)
		}
		return defaultConfig
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		log.Printf("Error parsing config.json: %v", err)
		return getDefaultConfig()
	}

	return config
}

func getDefaultConfig() Config {
	return Config{
		Launcher: LauncherConfig{
			Title:   "Terminal Gaming Suite",
			Version: "1.0.0",
			Author:  "Your Name",
			Theme:   "retro",
		},
		Games:   []GameEntry{},
		TUIApps: []GameEntry{},
		Stats: StatsConfig{
			Global: GlobalStats{
				GamesPlayed:          0,
				TotalTimeSeconds:     0,
				AchievementsUnlocked: 0,
				FavoriteGame:         "",
				LastPlayed:           "",
			},
			Achievements: []Achievement{},
		},
		Settings: SettingsConfig{
			Theme:              "retro",
			SoundEnabled:       true,
			AutoSave:           true,
			StatisticsTracking: true,
			Notifications:      true,
			BackupSaves:        true,
			ControllerSupport:  false,
			TerminalSize:       "auto",
		},
	}
}

func (m *Launcher) Init() tea.Cmd {
	return nil
}

func (m *Launcher) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	}

	return m, nil
}

func (m *Launcher) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.menuState {
	case MainMenu:
		return m.handleMainMenuKeys(msg)
	case TUIAppsMenu:
		return m.handleTUIAppsMenuKeys(msg)
	case StatsMenu:
		return m.handleStatsMenuKeys(msg)
	case OptionsMenu:
		return m.handleOptionsMenuKeys(msg)
	case CreditsMenu:
		return m.handleCreditsMenuKeys(msg)
	}
	return m, nil
}

func (m *Launcher) handleMainMenuKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "up", "k":
		if m.selectedGame > 0 {
			m.selectedGame--
		}
	case "down", "j":
		if m.selectedGame < len(m.games)-1 {
			m.selectedGame++
		}
	case "left":
		m.menuState = TUIAppsMenu
	case "right":
		m.menuState = TUIAppsMenu
	case "enter", " ":
		return m.launchGame(m.games[m.selectedGame])
	case "t":
		m.menuState = TUIAppsMenu
	case "s":
		m.menuState = StatsMenu
	case "o":
		m.menuState = OptionsMenu
	case "c":
		m.menuState = CreditsMenu
	}
	return m, nil
}

func (m *Launcher) handleTUIAppsMenuKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "up", "k":
		if m.selectedTUIApp > 0 {
			m.selectedTUIApp--
		}
	case "down", "j":
		if m.selectedTUIApp < len(m.tuiApps)-1 {
			m.selectedTUIApp++
		}
	case "left":
		m.menuState = MainMenu
	case "right":
		m.menuState = MainMenu
	case "enter", " ":
		return m.launchGame(m.tuiApps[m.selectedTUIApp])
	}
	return m, nil
}

func (m *Launcher) handleStatsMenuKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc", "backspace":
		m.menuState = MainMenu
	}
	return m, nil
}

func (m *Launcher) handleOptionsMenuKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc", "backspace":
		m.menuState = MainMenu
	}
	return m, nil
}

func (m *Launcher) handleCreditsMenuKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc", "backspace":
		m.menuState = MainMenu
	}
	return m, nil
}

func (m *Launcher) launchGame(game GameEntry) (tea.Model, tea.Cmd) {
	// Get the hardcoded tui-hub directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return m, tea.Quit
	}

	tuiHubDir := filepath.Join(homeDir, "projects", "active", "daily_use", "tui-hub")

	// Create the command
	var cmd *exec.Cmd
	if game.Path != "" && game.Path != "./" {
		// Use the hardcoded tui-hub directory as base
		fullPath := filepath.Join(tuiHubDir, game.Path)
		cmd = exec.Command("sh", "-c", fmt.Sprintf("cd %s && %s", fullPath, game.Command))
	} else {
		cmd = exec.Command("sh", "-c", game.Command)
	}

	// Ensure the command inherits the current environment (including PATH)
	cmd.Env = os.Environ()

	return m, tea.ExecProcess(
		cmd,
		func(err error) tea.Msg {
			if err != nil {
				return fmt.Sprintf("Error launching %s: %v", game.Name, err)
			}
			return fmt.Sprintf("Returned from %s", game.Name)
		},
	)
}

func (m *Launcher) View() string {
	switch m.menuState {
	case MainMenu:
		return m.renderMainMenu()
	case TUIAppsMenu:
		return m.renderTUIAppsMenu()
	case StatsMenu:
		return m.renderStatsMenu()
	case OptionsMenu:
		return m.renderOptionsMenu()
	case CreditsMenu:
		return m.renderCreditsMenu()
	}
	return ""
}

func (m *Launcher) renderMainMenu() string {
	var s strings.Builder

	// Header
	s.WriteString(headerStyle.Render("🎮 Terminal Gaming Suite"))
	s.WriteString("\n\n")

	// Table header
	s.WriteString(normalStyle.Render("Name"))
	s.WriteString(strings.Repeat(" ", 30))
	s.WriteString(normalStyle.Render("Description"))
	s.WriteString("\n")
	s.WriteString(strings.Repeat("─", m.width))
	s.WriteString("\n")

	// Games list
	for i, game := range m.games {
		nameCol := fmt.Sprintf("%s %s", game.Icon, game.Name)
		if i == m.selectedGame {
			nameCol = selectedStyle.Render(nameCol)
			s.WriteString(nameCol)
			s.WriteString(strings.Repeat(" ", 30-len(game.Name)-3))
			s.WriteString(selectedStyle.Render(game.Description))
		} else {
			s.WriteString(normalStyle.Render(nameCol))
			s.WriteString(strings.Repeat(" ", 30-len(game.Name)-3))
			s.WriteString(descStyle.Render(game.Description))
		}
		s.WriteString("\n")
	}

	// Commands
	s.WriteString("\n")
	s.WriteString(keyStyle.Render("enter"))
	s.WriteString(commandStyle.Render(": launch • "))
	s.WriteString(keyStyle.Render("↑↓"))
	s.WriteString(commandStyle.Render(": navigate • "))
	s.WriteString(keyStyle.Render("←→"))
	s.WriteString(commandStyle.Render(": swap menu • "))
	s.WriteString(keyStyle.Render("q"))
	s.WriteString(commandStyle.Render(": quit"))

	return s.String()
}

func (m *Launcher) renderTUIAppsMenu() string {
	var s strings.Builder

	// Header
	s.WriteString(headerStyle.Render("🖥️  TUI Applications"))
	s.WriteString("\n\n")

	// Table header
	s.WriteString(normalStyle.Render("Name"))
	s.WriteString(strings.Repeat(" ", 30))
	s.WriteString(normalStyle.Render("Description"))
	s.WriteString("\n")
	s.WriteString(strings.Repeat("─", m.width))
	s.WriteString("\n")

	// Apps list
	for i, app := range m.tuiApps {
		nameCol := fmt.Sprintf("%s %s", app.Icon, app.Name)
		if i == m.selectedTUIApp {
			nameCol = selectedStyle.Render(nameCol)
			s.WriteString(nameCol)
			s.WriteString(strings.Repeat(" ", 30-len(app.Name)-3))
			s.WriteString(selectedStyle.Render(app.Description))
		} else {
			s.WriteString(normalStyle.Render(nameCol))
			s.WriteString(strings.Repeat(" ", 30-len(app.Name)-3))
			s.WriteString(descStyle.Render(app.Description))
		}
		s.WriteString("\n")
	}

	// Commands
	s.WriteString("\n")
	s.WriteString(keyStyle.Render("enter"))
	s.WriteString(commandStyle.Render(": launch • "))
	s.WriteString(keyStyle.Render("↑↓"))
	s.WriteString(commandStyle.Render(": navigate • "))
	s.WriteString(keyStyle.Render("←→"))
	s.WriteString(commandStyle.Render(": swap menu • "))
	s.WriteString(keyStyle.Render("q"))
	s.WriteString(commandStyle.Render(": quit"))

	return s.String()
}

func (m *Launcher) renderStatsMenu() string {
	var s strings.Builder

	// Header
	s.WriteString(headerStyle.Render("📊 Statistics"))
	s.WriteString("\n\n")

	// Stats content
	s.WriteString(normalStyle.Render("Games Played: "))
	s.WriteString(commandStyle.Render(fmt.Sprintf("%d", m.stats.GamesPlayed)))
	s.WriteString("\n")

	s.WriteString(normalStyle.Render("Total Time: "))
	s.WriteString(commandStyle.Render(fmt.Sprintf("%d seconds", m.stats.TotalTimeSeconds)))
	s.WriteString("\n")

	s.WriteString(normalStyle.Render("Achievements: "))
	s.WriteString(commandStyle.Render(fmt.Sprintf("%d", m.stats.AchievementsUnlocked)))
	s.WriteString("\n")

	s.WriteString(normalStyle.Render("Favorite Game: "))
	s.WriteString(commandStyle.Render(m.stats.FavoriteGame))
	s.WriteString("\n\n")

	// Commands
	s.WriteString(keyStyle.Render("esc"))
	s.WriteString(commandStyle.Render(": back"))

	return s.String()
}

func (m *Launcher) renderOptionsMenu() string {
	var s strings.Builder

	// Header
	s.WriteString(headerStyle.Render("⚙️  Options"))
	s.WriteString("\n\n")

	s.WriteString(descStyle.Render("[Coming Soon]"))
	s.WriteString("\n\n")

	// Commands
	s.WriteString(keyStyle.Render("esc"))
	s.WriteString(commandStyle.Render(": back"))

	return s.String()
}

func (m *Launcher) renderCreditsMenu() string {
	var s strings.Builder

	// Header
	s.WriteString(headerStyle.Render("👨‍💻 Credits"))
	s.WriteString("\n\n")

	s.WriteString(normalStyle.Render("Terminal Gaming Suite"))
	s.WriteString("\n")
	s.WriteString(descStyle.Render("Developed with love"))
	s.WriteString("\n\n")

	// Commands
	s.WriteString(keyStyle.Render("esc"))
	s.WriteString(commandStyle.Render(": back"))

	return s.String()
}
