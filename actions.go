package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

var versionPattern = regexp.MustCompile(`v?\d+\.\d+\.\d+(?:[-+][0-9A-Za-z.-]+)?`)

func detectLocalVersion(binaryPath string) string {
	for _, arg := range []string{"--version", "-version"} {
		ctx, cancel := context.WithTimeout(context.Background(), 700*time.Millisecond)
		out, err := exec.CommandContext(ctx, binaryPath, arg).CombinedOutput()
		cancel()

		if version := extractVersion(string(out)); version != "" {
			return version
		}
		if err == nil {
			continue
		}
	}
	return ""
}

func extractVersion(raw string) string {
	return versionPattern.FindString(raw)
}

func scanLocalVersionsCmd(apps []suiteApp) tea.Cmd {
	return func() tea.Msg {
		updated := make([]suiteApp, len(apps))
		copy(updated, apps)
		for i := range updated {
			if !updated[i].Installed || updated[i].ResolvedPath == "" {
				continue
			}
			updated[i].LocalVersion = detectLocalVersion(updated[i].ResolvedPath)
		}
		return localVersionsScannedMsg{apps: updated}
	}
}

func checkVersionsCmd(apps []suiteApp) tea.Cmd {
	return func() tea.Msg {
		client := &http.Client{Timeout: 8 * time.Second}
		updated := make([]suiteApp, len(apps))
		copy(updated, apps)

		checked := 0
		for i := range updated {
			updated[i].LatestVersion = ""
			updated[i].UpdateAvailable = false
			if !updated[i].Installed {
				continue
			}
			tag, err := fetchLatestRelease(client, updated[i].Repo)
			if err != nil {
				return errMsg{err: fmt.Errorf("check updates for %s: %w", updated[i].Name, err)}
			}
			updated[i].LatestVersion = tag
			updated[i].UpdateAvailable = shouldOfferUpdate(updated[i].LocalVersion, tag)
			checked++
		}

		status := "No installed apps to check."
		if checked > 0 {
			status = fmt.Sprintf("Checked %d installed app", checked)
			if checked != 1 {
				status += "s"
			}
			status += "."
		}
		return versionsCheckedMsg{apps: updated, status: status}
	}
}

func installOrUpdateCmd(app suiteApp, cfg config) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("sh", "-c", fmt.Sprintf("curl -fsSL https://raw.githubusercontent.com/%s/main/install.sh | bash", app.Repo))
		cmd.Env = os.Environ()
		if out, err := cmd.CombinedOutput(); err != nil {
			msg := strings.TrimSpace(string(out))
			if msg == "" {
				msg = err.Error()
			}
			return installFinishedMsg{
				appID: app.ID,
				apps:  refreshAppsWithConfig(cfg),
				err:   fmt.Errorf("%s install failed: %s", app.Name, msg),
			}
		}

		refreshed := refreshAppsWithConfig(cfg)
		action := "Installed "
		if app.Installed {
			action = "Updated "
		}
		return installFinishedMsg{
			appID:  app.ID,
			apps:   refreshed,
			status: action + app.Name + ".",
		}
	}
}

func fetchLatestRelease(client *http.Client, repo string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, "https://api.github.com/repos/"+repo+"/releases/latest", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "tui-hub")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("github returned %s", resp.Status)
	}

	var payload struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", err
	}
	return payload.TagName, nil
}

func shouldOfferUpdate(localVersion, latestVersion string) bool {
	if latestVersion == "" {
		return false
	}
	if localVersion == "" {
		return true
	}
	return normalizeVersion(localVersion) != normalizeVersion(latestVersion)
}

func normalizeVersion(v string) string {
	return strings.TrimPrefix(strings.TrimSpace(v), "v")
}
