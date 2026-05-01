package main

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

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
		return m, nil
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
		return m, nil
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
		return m, installOrUpdateCmd(app)
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
		return m, installOrUpdateCmd(app)
	}
	return m, nil
}

func (m *model) moveSelection(delta int) {
	items := m.visibleApps()
	if len(items) == 0 {
		m.selected[m.page] = 0
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
			continue
		}
		if m.selected[page] >= len(items) {
			m.selected[page] = len(items) - 1
		}
		if m.selected[page] < 0 {
			m.selected[page] = 0
		}
	}
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
	apps := make([]suiteApp, 0, len(builtInCatalog()))
	for _, entry := range builtInCatalog() {
		app := suiteApp{appCatalogEntry: entry}
		if path, err := exec.LookPath(entry.Binary); err == nil {
			app.Installed = true
			app.ResolvedPath = path
			app.LocalVersion = detectLocalVersion(path)
		}
		apps = append(apps, app)
	}

	sort.Slice(apps, func(i, j int) bool {
		if apps[i].Installed != apps[j].Installed {
			return apps[i].Installed
		}
		return strings.ToLower(apps[i].Name) < strings.ToLower(apps[j].Name)
	})
	return apps
}

func refreshAppsCmd(status string) tea.Cmd {
	return func() tea.Msg {
		return appsRefreshedMsg{apps: refreshApps(status), status: status}
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
				return errMsg{err: fmt.Errorf("launch %s: %w", app.Name, err)}
			}
			return statusMsg{text: "Returned from " + app.Name}
		}),
	)
}
