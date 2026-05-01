package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("230"))

	activeTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("230")).
			Background(lipgloss.Color("57")).
			Padding(0, 1)

	tabStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245")).
			Padding(0, 1)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("230")).
			Background(lipgloss.Color("57")).
			Bold(true)

	nameStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("230")).
			Bold(true)

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245"))

	accentStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("117")).
			Bold(true)

	warnStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("214")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("203")).
			Bold(true)

	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1)
)

func (m model) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("tui-hub"))
	b.WriteString(dimStyle.Render("  launch and manage your tui-suite"))
	b.WriteString("\n")
	b.WriteString(m.renderTabs())
	b.WriteString("\n\n")
	b.WriteString(m.renderList())
	b.WriteString("\n")
	b.WriteString(m.renderStatus())
	b.WriteString("\n")
	b.WriteString(m.renderHelp())

	return b.String()
}

func (m model) renderTabs() string {
	installedLabel := fmt.Sprintf("1 Installed (%d)", len(m.appsForPage(pageInstalled)))
	availableLabel := fmt.Sprintf("2 Available (%d)", len(m.appsForPage(pageAvailable)))

	installed := tabStyle.Render(installedLabel)
	available := tabStyle.Render(availableLabel)
	if m.page == pageInstalled {
		installed = activeTabStyle.Render(installedLabel)
	} else {
		available = activeTabStyle.Render(availableLabel)
	}
	return installed + " " + available
}

func (m model) renderList() string {
	items := m.visibleApps()
	if len(items) == 0 {
		text := "No apps here."
		if m.page == pageInstalled {
			text = "No suite apps found on PATH yet. Switch to Available and press i to install one."
		}
		return panelStyle.Width(max(40, m.width-2)).Render(dimStyle.Render(text))
	}

	contentWidth := 90
	if m.width > 10 {
		contentWidth = m.width - 6
	}
	if contentWidth < 40 {
		contentWidth = 40
	}

	var rows []string
	for i, app := range items {
		line := m.renderAppRow(i == m.selected[m.page], app, contentWidth-4)
		rows = append(rows, line)
	}
	return panelStyle.Width(contentWidth).Render(strings.Join(rows, "\n\n"))
}

func (m model) renderAppRow(selected bool, app suiteApp, width int) string {
	versionBits := []string{}
	if app.LocalVersion != "" {
		versionBits = append(versionBits, "local "+app.LocalVersion)
	}
	if app.LatestVersion != "" {
		versionBits = append(versionBits, "latest "+app.LatestVersion)
	}
	meta := strings.Join(versionBits, "  ")
	if meta == "" {
		if app.Installed {
			meta = "installed"
		} else {
			meta = "not installed"
		}
	}
	if app.UpdateAvailable {
		meta += "  update available"
	}

	line1 := nameStyle.Render(app.Name)
	line2 := dimStyle.Render(app.Description)
	line3 := accentStyle.Render(meta)
	if app.UpdateAvailable {
		line3 = warnStyle.Render(meta)
	}

	row := strings.Join([]string{line1, line2, line3}, "\n")
	row = lipgloss.NewStyle().Width(width).Render(row)
	if selected {
		return selectedStyle.Width(width).Render(row)
	}
	return row
}

func (m model) renderStatus() string {
	if m.errMsg != "" {
		return errorStyle.Render(m.errMsg)
	}
	if m.busy {
		return warnStyle.Render(m.status)
	}
	if m.status != "" {
		return accentStyle.Render(m.status)
	}
	return dimStyle.Render("Ready.")
}

func (m model) renderHelp() string {
	parts := []string{
		accentStyle.Render("tab/1/2") + dimStyle.Render(" switch"),
		accentStyle.Render("j/k") + dimStyle.Render(" move"),
	}
	if m.page == pageInstalled {
		parts = append(parts, accentStyle.Render("enter")+dimStyle.Render(" launch"))
		parts = append(parts, accentStyle.Render("u")+dimStyle.Render(" update"))
	} else {
		parts = append(parts, accentStyle.Render("i")+dimStyle.Render(" install"))
	}
	parts = append(parts, accentStyle.Render("r")+dimStyle.Render(" check versions"))
	parts = append(parts, accentStyle.Render("q")+dimStyle.Render(" quit"))
	return strings.Join(parts, dimStyle.Render("  •  "))
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
