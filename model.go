package main

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	pageInstalled = "installed"
	pageAvailable = "available"
)

type suiteApp struct {
	appCatalogEntry
	Installed       bool
	ResolvedPath    string
	LocalVersion    string
	LatestVersion   string
	UpdateAvailable bool
}

type model struct {
	version    string
	cfg        config
	page       string
	width      int
	height     int
	selected   map[string]int
	scroll     map[string]int
	apps       []suiteApp
	status     string
	errMsg     string
	busy       bool
	checking   bool
	installing string
}

type statusMsg struct {
	text string
}

type errMsg struct {
	err error
}

type appsRefreshedMsg struct {
	apps   []suiteApp
	status string
}

type versionsCheckedMsg struct {
	apps   []suiteApp
	status string
}

type localVersionsScannedMsg struct {
	apps []suiteApp
}

type launchFinishedMsg struct {
	appID string
	err   error
}

type installFinishedMsg struct {
	appID  string
	apps   []suiteApp
	status string
	err    error
}

func initialModel(appVersion string) model {
	cfg := loadConfig()
	m := model{
		version:  appVersion,
		cfg:      cfg,
		page:     cfg.LastPage,
		selected: map[string]int{pageInstalled: 0, pageAvailable: 0},
		scroll:   map[string]int{pageInstalled: 0, pageAvailable: 0},
		status:   "Scanning installed apps...",
	}
	if m.page != pageAvailable {
		m.page = pageInstalled
	}
	return m
}

func (m model) Init() tea.Cmd {
	return refreshAppsCmd("Ready")
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case statusMsg:
		m.status = msg.text
		m.errMsg = ""
		return m, nil
	case errMsg:
		m.errMsg = msg.err.Error()
		m.status = ""
		m.busy = false
		m.checking = false
		m.installing = ""
		return m, nil
	case appsRefreshedMsg:
		m.apps = msg.apps
		m.status = msg.status
		m.errMsg = ""
		m.busy = false
		m.syncSelection()
		return m, scanLocalVersionsCmd(m.apps)
	case localVersionsScannedMsg:
		m.apps = msg.apps
		m.syncSelection()
		return m, nil
	case launchFinishedMsg:
		if msg.err != nil {
			m.errMsg = msg.err.Error()
			m.status = ""
			m.busy = false
			m.checking = false
			m.installing = ""
			return m, nil
		}
		appName := msg.appID
		if state, ok := m.recordLaunch(msg.appID, time.Now()); ok {
			appName = state.Name
		}
		m.status = "Returned from " + appName
		m.errMsg = ""
		m.busy = false
		sortApps(m.apps, m.cfg)
		m.syncSelection()
		return m, saveConfigCmd(m.cfg)
	case versionsCheckedMsg:
		m.apps = msg.apps
		m.status = msg.status
		m.errMsg = ""
		m.busy = false
		m.checking = false
		m.syncSelection()
		return m, nil
	case installFinishedMsg:
		m.apps = msg.apps
		m.busy = false
		m.installing = ""
		if msg.err != nil {
			m.errMsg = msg.err.Error()
			m.status = ""
		} else {
			m.errMsg = ""
			m.status = msg.status
		}
		m.syncSelection()
		return m, scanLocalVersionsCmd(m.apps)
	case tea.KeyMsg:
		return m.updateKey(msg)
	}
	return m, nil
}

func (m model) updateKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "1":
		return m.switchPage(pageInstalled)
	case "2":
		return m.switchPage(pageAvailable)
	case "tab", "right", "l":
		if m.page == pageInstalled {
			return m.switchPage(pageAvailable)
		}
		return m.switchPage(pageInstalled)
	case "shift+tab", "left", "h":
		if m.page == pageAvailable {
			return m.switchPage(pageInstalled)
		}
		return m.switchPage(pageAvailable)
	case "up", "k":
		m.moveSelection(-1)
		return m, nil
	case "down", "j":
		m.moveSelection(1)
		return m, nil
	case "pgdown", "ctrl+d":
		m.moveSelection(m.visibleRowCount())
		return m, nil
	case "pgup", "ctrl+u":
		m.moveSelection(-m.visibleRowCount())
		return m, nil
	case "g", "home":
		m.selected[m.page] = 0
		m.scroll[m.page] = 0
		return m, nil
	case "G", "end":
		items := m.visibleApps()
		if len(items) == 0 {
			return m, nil
		}
		m.selected[m.page] = len(items) - 1
		m.ensureSelectionVisible()
		return m, nil
	case "r":
		if m.busy || m.checking {
			return m, nil
		}
		m.busy = true
		m.checking = true
		m.status = "Checking latest releases..."
		m.errMsg = ""
		return m, checkVersionsCmd(m.apps)
	case "enter":
		if m.busy {
			return m, nil
		}
		app, ok := m.selectedApp()
		if !ok || !app.Installed || m.page != pageInstalled {
			return m, nil
		}
		m.status = "Launching " + app.Name + "..."
		m.errMsg = ""
		return m, launchAppCmd(app)
	case "i":
		if m.busy || m.page != pageAvailable {
			return m, nil
		}
		app, ok := m.selectedApp()
		if !ok || app.Installed {
			return m, nil
		}
		m.busy = true
		m.installing = app.ID
		m.status = "Installing " + app.Name + "..."
		m.errMsg = ""
		return m, installOrUpdateCmd(app, m.cfg)
	case "u":
		if m.busy || m.page != pageInstalled {
			return m, nil
		}
		app, ok := m.selectedApp()
		if !ok || !app.Installed || !app.UpdateAvailable {
			return m, nil
		}
		m.busy = true
		m.installing = app.ID
		m.status = "Updating " + app.Name + "..."
		m.errMsg = ""
		return m, installOrUpdateCmd(app, m.cfg)
	}
	return m, nil
}

func (m *model) moveSelection(delta int) {
	items := m.visibleApps()
	if len(items) == 0 {
		m.selected[m.page] = 0
		m.scroll[m.page] = 0
		return
	}
	idx := m.selected[m.page] + delta
	if idx < 0 {
		idx = 0
	}
	if idx >= len(items) {
		idx = len(items) - 1
	}
	m.selected[m.page] = idx
	m.ensureSelectionVisible()
}

func (m model) switchPage(page string) (tea.Model, tea.Cmd) {
	m.page = page
	m.cfg.LastPage = page
	m.syncSelection()
	return m, saveConfigCmd(m.cfg)
}

func (m *model) syncSelection() {
	for _, page := range []string{pageInstalled, pageAvailable} {
		items := m.appsForPage(page)
		if len(items) == 0 {
			m.selected[page] = 0
			m.scroll[page] = 0
			continue
		}
		if m.selected[page] >= len(items) {
			m.selected[page] = len(items) - 1
		}
		if m.selected[page] < 0 {
			m.selected[page] = 0
		}
		maxScroll := len(items) - m.visibleRowCount()
		if maxScroll < 0 {
			maxScroll = 0
		}
		if m.scroll[page] > maxScroll {
			m.scroll[page] = maxScroll
		}
		if m.scroll[page] < 0 {
			m.scroll[page] = 0
		}
		prevPage := m.page
		m.page = page
		m.ensureSelectionVisible()
		m.page = prevPage
	}
}

func (m *model) ensureSelectionVisible() {
	items := m.visibleApps()
	if len(items) == 0 {
		m.scroll[m.page] = 0
		return
	}
	rows := m.visibleRowCount()
	if rows < 1 {
		rows = 1
	}
	idx := m.selected[m.page]
	top := m.scroll[m.page]
	bottom := top + rows - 1
	if idx < top {
		m.scroll[m.page] = idx
	} else if idx > bottom {
		m.scroll[m.page] = idx - rows + 1
	}
	if m.scroll[m.page] < 0 {
		m.scroll[m.page] = 0
	}
	maxScroll := len(items) - rows
	if maxScroll < 0 {
		maxScroll = 0
	}
	if m.scroll[m.page] > maxScroll {
		m.scroll[m.page] = maxScroll
	}
}

func (m model) visibleRowCount() int {
	if m.height <= 0 {
		return 8
	}
	rows := m.height - 9
	if rows < 4 {
		rows = 4
	}
	if rows > 14 {
		rows = 14
	}
	return rows
}

func (m model) visibleApps() []suiteApp {
	return m.appsForPage(m.page)
}

func (m model) appsForPage(page string) []suiteApp {
	var out []suiteApp
	for _, app := range m.apps {
		if page == pageInstalled && app.Installed {
			out = append(out, app)
		}
		if page == pageAvailable && !app.Installed {
			out = append(out, app)
		}
	}
	return out
}

func (m model) selectedApp() (suiteApp, bool) {
	items := m.visibleApps()
	if len(items) == 0 {
		return suiteApp{}, false
	}
	idx := m.selected[m.page]
	if idx < 0 || idx >= len(items) {
		return suiteApp{}, false
	}
	return items[idx], true
}

func refreshApps(status string) []suiteApp {
	return refreshAppsWithConfig(defaultConfig())
}

func refreshAppsWithConfig(cfg config) []suiteApp {
	apps := make([]suiteApp, 0, len(builtInCatalog()))
	for _, entry := range builtInCatalog() {
		app := suiteApp{appCatalogEntry: entry}
		if path, err := exec.LookPath(entry.Binary); err == nil {
			app.Installed = true
			app.ResolvedPath = path
		}
		apps = append(apps, app)
	}

	sortApps(apps, cfg)
	return apps
}

func refreshAppsCmd(status string) tea.Cmd {
	return func() tea.Msg {
		return appsRefreshedMsg{apps: refreshAppsWithConfig(loadConfig()), status: status}
	}
}

func saveConfigCmd(cfg config) tea.Cmd {
	return func() tea.Msg {
		if err := saveConfig(cfg); err != nil {
			return errMsg{err: err}
		}
		return nil
	}
}

func launchAppCmd(app suiteApp) tea.Cmd {
	cmd := exec.Command(app.ResolvedPath)
	cmd.Env = os.Environ()
	return tea.Sequence(
		tea.ExecProcess(cmd, func(err error) tea.Msg {
			if err != nil {
				return launchFinishedMsg{appID: app.ID, err: fmt.Errorf("launch %s: %w", app.Name, err)}
			}
			return launchFinishedMsg{appID: app.ID}
		}),
	)
}

func sortApps(apps []suiteApp, cfg config) {
	sort.Slice(apps, func(i, j int) bool {
		if apps[i].Installed != apps[j].Installed {
			return apps[i].Installed
		}
		if apps[i].Installed {
			left := frecencyScore(cfg.AppState[apps[i].ID])
			right := frecencyScore(cfg.AppState[apps[j].ID])
			if left != right {
				return left > right
			}
			leftTime := lastLaunchedAt(cfg.AppState[apps[i].ID])
			rightTime := lastLaunchedAt(cfg.AppState[apps[j].ID])
			if !leftTime.Equal(rightTime) {
				return leftTime.After(rightTime)
			}
		}
		return strings.ToLower(apps[i].Name) < strings.ToLower(apps[j].Name)
	})
}

func frecencyScore(state appState) float64 {
	if state.LaunchCount <= 0 {
		return 0
	}
	last := lastLaunchedAt(state)
	if last.IsZero() {
		return float64(state.LaunchCount)
	}
	days := time.Since(last).Hours() / 24
	if days < 0 {
		days = 0
	}
	return float64(state.LaunchCount) + 1/(1+days)
}

func lastLaunchedAt(state appState) time.Time {
	if state.LastLaunched == "" {
		return time.Time{}
	}
	t, err := time.Parse(time.RFC3339, state.LastLaunched)
	if err != nil {
		return time.Time{}
	}
	return t
}

func (m *model) recordLaunch(appID string, when time.Time) (suiteApp, bool) {
	if m.cfg.AppState == nil {
		m.cfg.AppState = map[string]appState{}
	}
	state := m.cfg.AppState[appID]
	state.LaunchCount++
	state.LastLaunched = when.Format(time.RFC3339)
	m.cfg.AppState[appID] = state

	for _, app := range m.apps {
		if app.ID == appID {
			return app, true
		}
	}
	return suiteApp{}, false
}
