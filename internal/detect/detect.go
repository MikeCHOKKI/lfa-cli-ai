package detect

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
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

	client := &http.Client{Timeout: 1 * time.Second}
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
	// Map Go arch names to OpenCode release naming
	ocArch := arch
	if ocArch == "amd64" {
		ocArch = "x64"
	}
	// Linux uses .tar.gz, Darwin and Windows use .zip
	ext := ".tar.gz"
	if o == Windows || o == Darwin {
		ext = ".zip"
	}
	return fmt.Sprintf(
		"https://github.com/anomalyco/opencode/releases/download/v%s/opencode-%s-%s%s",
		version, o.String(), ocArch, ext,
	)
}

// ─── PostgreSQL ──────────────────────────────────────────────────────────────

type PGStatus struct {
	Installed    bool
	Running      bool
	Version      string
}

func DetectPostgreSQL() PGStatus {
	status := PGStatus{}

	// Check if psql is installed
	psqlPath, err := exec.LookPath("psql")
	if err != nil {
		return status
	}
	status.Installed = true

	// Get version
	out, err := exec.Command(psqlPath, "--version").Output()
	if err == nil {
		status.Version = strings.TrimSpace(string(out))
	}

	// Check if PostgreSQL is running (try to connect to default socket)
	pgIsReady, err := exec.LookPath("pg_isready")
	if err == nil {
		if err := exec.Command(pgIsReady, "-q").Run(); err == nil {
			status.Running = true
		}
	} else {
		// Fallback: try connecting with psql
		err := exec.Command(psqlPath, "-h", "localhost", "-U", "postgres", "-c", "SELECT 1").Run()
		status.Running = err == nil
	}

	return status
}

func DetectPostgreSQLConnectivity(host, port, user, password, dbname string) bool {
	// Use PG environment variables for the check
	cmd := exec.Command("psql", "-h", host, "-p", port, "-U", user, "-d", dbname, "-c", "SELECT 1")
	cmd.Env = append(cmd.Env, "PGPASSWORD="+password)
	return cmd.Run() == nil
}

func IsPostgreSQLSupported(o OS) bool {
	switch o {
	case Linux, Darwin, Windows:
		return true
	default:
		return false
	}
}
