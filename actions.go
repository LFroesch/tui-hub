package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
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
		client := &http.Client{Timeout: 15 * time.Second}
		version, err := fetchLatestRelease(client, app.Repo)
		if err != nil {
			return installFinishedMsg{
				appID: app.ID,
				apps:  refreshAppsWithConfig(cfg),
				err:   fmt.Errorf("%s install failed: resolve latest release: %w", app.Name, err),
			}
		}

		script, err := fetchInstallScript(client, app.Repo, version)
		if err != nil {
			return installFinishedMsg{
				appID: app.ID,
				apps:  refreshAppsWithConfig(cfg),
				err:   fmt.Errorf("%s install failed: fetch installer for %s: %w", app.Name, version, err),
			}
		}

		cmd := installCommand()
		cmd.Env = append(os.Environ(), "VERSION="+version)
		cmd.Stdin = strings.NewReader(script)
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

func fetchInstallScript(client *http.Client, repo, version string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, installScriptURL(repo, version), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "tui-hub")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("github returned %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func installScriptURL(repo, version string) string {
	name := "install.sh"
	if runtime.GOOS == "windows" {
		name = "install.ps1"
	}
	return fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s", repo, version, name)
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
	normalized := strings.TrimPrefix(strings.TrimSpace(v), "v")
	return strings.TrimSuffix(normalized, "-dirty")
}

func installCommand() *exec.Cmd {
	if runtime.GOOS == "windows" {
		if path, err := exec.LookPath("pwsh"); err == nil {
			return exec.Command(path, "-NoLogo", "-NoProfile", "-Command", "-")
		}
		return exec.Command("powershell", "-NoLogo", "-NoProfile", "-Command", "-")
	}
	return exec.Command("bash")
}
