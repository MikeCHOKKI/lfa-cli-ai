package ui

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// Notify envoie une notification de bureau multi-plateforme.
//   - Linux : notify-send (libnotify)
//   - macOS : osascript (script AppleScript)
//   - Windows: PowerShell pop-up (fenêtre MessageBox)
//
// Retourne une erreur si la notification ne peut pas être envoyée.
func Notify(title, message string) error {
	switch runtime.GOOS {
	case "linux":
		return notifyLinux(title, message)
	case "darwin":
		return notifyDarwin(title, message)
	case "windows":
		return notifyWindows(title, message)
	default:
		return fmt.Errorf("notifications not supported on %s", runtime.GOOS)
	}
}

func notifyLinux(title, message string) error {
	return exec.Command("notify-send", title, message).Run()
}

func notifyDarwin(title, message string) error {
	script := fmt.Sprintf(`display notification "%s" with title "%s"`,
		escapeAppleScript(message),
		escapeAppleScript(title),
	)
	return exec.Command("osascript", "-e", script).Run()
}

func notifyWindows(title, message string) error {
	// Use PowerShell to show a modern Windows notification toast.
	// Falls back to msg.exe if PowerShell is unavailable.
	psCmd := fmt.Sprintf(
		`[System.Windows.MessageBox]::Show('%s','%s')`,
		escapePS(message),
		escapePS(title),
	)
	cmd := exec.Command("powershell", "-NoProfile", "-Command", psCmd)
	if err := cmd.Run(); err == nil {
		return nil
	}
	// Fallback: use msg.exe for older Windows
	return exec.Command("msg", "*", fmt.Sprintf("%s: %s", title, message)).Run()
}

func escapeAppleScript(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	return s
}

func escapePS(s string) string {
	s = strings.ReplaceAll(s, `'`, `''`)
	return s
}
