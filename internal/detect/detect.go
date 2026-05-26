package detect

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

type OS int

const (
	Linux   OS = 0
	Windows OS = 1
	Darwin  OS = 2
)

func DetectOS() OS {
	switch runtime.GOOS {
	case "linux":
		return Linux
	case "windows":
		return Windows
	case "darwin":
		return Darwin
	default:
		return Linux
	}
}

func (o OS) String() string {
	switch o {
	case Linux:
		return "linux"
	case Windows:
		return "windows"
	case Darwin:
		return "darwin"
	default:
		return "unknown"
	}
}

func DetectOpenCode() (bool, string) {
	path, err := exec.LookPath("opencode")
	if err != nil {
		return false, ""
	}
	return true, path
}

func DetectOllama() (bool, bool) {
	_, err := exec.LookPath("ollama")
	installed := err == nil

	if !installed {
		return false, false
	}

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get("http://localhost:11434/api/tags")
	if err != nil {
		return installed, false
	}
	defer resp.Body.Close()

	return installed, resp.StatusCode == http.StatusOK
}

func GetOpenCodeConfigDir(o OS) string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "~"
	}
	switch o {
	case Windows:
		appData := os.Getenv("APPDATA")
		if appData == "" {
			appData = filepath.Join(home, "AppData", "Roaming")
		}
		return filepath.Join(appData, "opencode")
	case Darwin:
		return filepath.Join(home, "Library", "Application Support", "opencode")
	default:
		return filepath.Join(home, ".config", "opencode")
	}
}

func GetArch() string {
	arch := runtime.GOARCH
	switch arch {
	case "amd64":
		return "amd64"
	case "arm64":
		return "arm64"
	default:
		return "amd64"
	}
}

func GetOpenCodeDownloadURL(o OS, arch, version string) string {
	return fmt.Sprintf(
		"https://github.com/anomalyco/opencode/releases/download/v%s/opencode_%s_%s.tar.gz",
		version, o.String(), arch,
	)
}
