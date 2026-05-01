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
			Underline(true)

	tabStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245"))

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("230")).
			Background(lipgloss.Color("57")).
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

	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("117")).
			Bold(true)
)

func (m model) View() string {
	if m.width == 0 || m.height == 0 {
		return "loading..."
	}

	header := m.renderHeader()
	sep := dimStyle.Render(strings.Repeat("─", max(20, m.width)))
	content := m.renderTable()
	status := m.renderStatus()
	footer := m.renderHelp()

	return lipgloss.JoinVertical(lipgloss.Left, header, sep, content, sep, footer, status)
}

func (m model) renderHeader() string {
	title := titleStyle.Render("tui-hub") + " " + dimStyle.Render(m.version)
	tabs := m.renderTabs()
	right := dimStyle.Render(fmt.Sprintf("%d apps", len(m.visibleApps())))

	line := title + "  " + tabs
	gap := m.width - lipgloss.Width(line) - lipgloss.Width(right)
	if gap < 2 {
		return line
	}
	return line + strings.Repeat(" ", gap) + right
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
	return installed + dimStyle.Render("  │  ") + available
}

func (m model) renderTable() string {
	items := m.visibleApps()
	panelWidth := m.width - 2
	if panelWidth < 50 {
		panelWidth = 50
	}

	if len(items) == 0 {
		text := "No apps here."
		if m.page == pageInstalled {
			text = "No suite apps found on PATH yet. Switch to Available and press i to install one."
		}
		return panelStyle.Width(panelWidth).Render(dimStyle.Render(text))
	}

	rowsVisible := m.visibleRowCount()
	start := m.scroll[m.page]
	if start < 0 {
		start = 0
	}
	if start > len(items) {
		start = len(items)
	}
	end := start + rowsVisible
	if end > len(items) {
		end = len(items)
	}

	innerWidth := panelWidth - 4
	if innerWidth < 46 {
		innerWidth = 46
	}
	nameW, versionW, statusW := 16, 12, 16
	descW := innerWidth - nameW - versionW - statusW - 6
	if descW < 16 {
		descW = 16
	}

	lines := []string{
		m.renderHeaderRow(nameW, descW, versionW, statusW),
		dimStyle.Render(strings.Repeat("─", innerWidth)),
	}
	for i := start; i < end; i++ {
		lines = append(lines, m.renderDataRow(i == m.selected[m.page], items[i], nameW, descW, versionW, statusW))
	}
	for len(lines) < rowsVisible+2 {
		lines = append(lines, strings.Repeat(" ", innerWidth))
	}

	meta := dimStyle.Render(fmt.Sprintf("rows %d-%d of %d", start+1, end, len(items)))
	if len(items) <= rowsVisible {
		meta = dimStyle.Render(fmt.Sprintf("%d rows", len(items)))
	}
	lines = append(lines, "", meta)

	return panelStyle.Width(panelWidth).Render(strings.Join(lines, "\n"))
}

func (m model) renderHeaderRow(nameW, descW, versionW, statusW int) string {
	return headerStyle.Render(
		padRight("Name", nameW) + "  " +
			padRight("Description", descW) + "  " +
			padRight("Version", versionW) + "  " +
			padRight("Status", statusW),
	)
}

func (m model) renderDataRow(selected bool, app suiteApp, nameW, descW, versionW, statusW int) string {
	version := app.LocalVersion
	if version == "" {
		version = "-"
	}

	status := "available"
	if app.Installed {
		status = "installed"
	}
	if app.UpdateAvailable {
		status = "update -> " + app.LatestVersion
	}

	nameCell := m.renderNameCell(app, nameW)
	line := nameCell + "  " +
		padRight(truncate(app.Description, descW), descW) + "  " +
		padRight(version, versionW) + "  " +
		padRight(status, statusW)

	if selected {
		return selectedStyle.Width(nameW + descW + versionW + statusW + 6).Render(line)
	}
	if app.UpdateAvailable {
		return warnStyle.Render(line)
	}
	return line
}

func (m model) renderNameCell(app suiteApp, width int) string {
	label := strings.TrimSpace(app.Icon + " " + app.Name)
	styled := lipgloss.NewStyle().
		Foreground(lipgloss.Color(app.Color)).
		Bold(true).
		Render(truncate(label, width))
	return padRight(styled, width)
}

func (m model) renderStatus() string {
	if m.errMsg != "" {
		return errorStyle.Render("  " + m.errMsg)
	}
	if m.busy {
		return warnStyle.Render("  " + m.status)
	}
	if m.status != "" {
		return accentStyle.Render("  " + m.status)
	}
	return dimStyle.Render("  Ready.")
}

func (m model) renderHelp() string {
	parts := []string{
		accentStyle.Render("j/k") + " " + dimStyle.Render("move"),
		accentStyle.Render("ctrl+u/d") + " " + dimStyle.Render("page"),
		accentStyle.Render("g/G") + " " + dimStyle.Render("top/bottom"),
		accentStyle.Render("tab/1/2") + " " + dimStyle.Render("switch"),
	}
	if m.page == pageInstalled {
		parts = append(parts, accentStyle.Render("enter")+" "+dimStyle.Render("launch"))
		parts = append(parts, accentStyle.Render("u")+" "+dimStyle.Render("update"))
	} else {
		parts = append(parts, accentStyle.Render("i")+" "+dimStyle.Render("install"))
	}
	parts = append(parts, accentStyle.Render("r")+" "+dimStyle.Render("check versions"))
	parts = append(parts, accentStyle.Render("q")+" "+dimStyle.Render("quit"))
	return "  " + strings.Join(parts, dimStyle.Render("  ·  "))
}

func padRight(s string, width int) string {
	w := lipgloss.Width(s)
	if w >= width {
		return truncate(s, width)
	}
	return s + strings.Repeat(" ", width-w)
}

func truncate(s string, width int) string {
	if width <= 0 {
		return ""
	}
	if lipgloss.Width(s) <= width {
		return s
	}
	if width == 1 {
		return "…"
	}
	runes := []rune(s)
	out := ""
	for _, r := range runes {
		if lipgloss.Width(out+string(r)+"…") > width {
			break
		}
		out += string(r)
	}
	return out + "…"
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
