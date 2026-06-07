package installer

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/lfa-cli/lfa-cli-ai/internal/config"
	"github.com/lfa-cli/lfa-cli-ai/internal/detect"
)

const DefaultVersion = "latest"
const opencodeBinName = "opencode"

var assetsFS fs.FS

func SetAssetsFS(fs fs.FS) {
	assetsFS = fs
}

func EnsureDirectories(o detect.OS) error {
	configDir := detect.GetOpenCodeConfigDir(o)
	dirs := []string{
		configDir,
		filepath.Join(configDir, "agents"),
		filepath.Join(configDir, "skills"),
		filepath.Join(configDir, "mcp-servers"),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return fmt.Errorf("create dir %s: %w", d, err)
		}
	}
	return nil
}

// GetLatestOpenCodeVersion fetches the latest OpenCode version from GitHub.
func GetLatestOpenCodeVersion() (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get("https://api.github.com/repos/anomalyco/opencode/releases/latest")
	if err != nil {
		return "", fmt.Errorf("fetch latest release: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned %s", resp.Status)
	}
	// Parse JSON to extract tag_name
	var release struct{ TagName string `json:"tag_name"` }
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", fmt.Errorf("parse release: %w", err)
	}
	return strings.TrimPrefix(release.TagName, "v"), nil
}

func InstallOpenCode(o detect.OS, version string) error {
	arch := detect.GetArch()

	// Resolve "latest" to actual version
	if version == "latest" || version == "" {
		latest, err := GetLatestOpenCodeVersion()
		if err != nil {
			return fmt.Errorf("resolve latest version: %w", err)
		}
		version = latest
	}

	url := detect.GetOpenCodeDownloadURL(o, arch, version)

	fmt.Fprintf(os.Stderr, "Downloading OpenCode from %s\n", url)

	// Use appropriate temp file extension
	tmpExt := ".tar.gz"
	if o == detect.Windows || o == detect.Darwin {
		tmpExt = ".zip"
	}
	tmpFile, err := os.CreateTemp("", "opencode-*"+tmpExt)
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: %s", resp.Status)
	}

	progress := &ProgressWriter{
		Total:      resp.ContentLength,
		StatusLine: "Downloading OpenCode",
	}
	if _, err := io.Copy(tmpFile, io.TeeReader(resp.Body, progress)); err != nil {
		return fmt.Errorf("write download: %w", err)
	}
	tmpFile.Close()

	binDir, err := getBinDir(o)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return fmt.Errorf("create bin dir: %w", err)
	}

	if err := extractBinary(tmpFile.Name(), binDir, o); err != nil {
		return fmt.Errorf("extract: %w", err)
	}

	fmt.Fprintf(os.Stderr, "OpenCode installed at %s\n", filepath.Join(binDir, opencodeBinName))
	return nil
}

func getBinDir(o detect.OS) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	switch o {
	case detect.Windows:
		return filepath.Join(home, "AppData", "Local", "opencode", "bin"), nil
	default:
		return filepath.Join(home, ".local", "bin"), nil
	}
}

func extractBinary(archivePath, destDir string, o detect.OS) error {
	// .zip is used for Windows and macOS
	if strings.HasSuffix(archivePath, ".zip") {
		return extractBinaryZip(archivePath, destDir)
	}
	// .tar.gz is used for Linux
	return extractBinaryTarGz(archivePath, destDir)
}

func extractBinaryTarGz(archivePath, destDir string) error {
	f, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer f.Close()

	gzr, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(io.LimitReader(gzr, 500*1024*1024))
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		name := filepath.Base(header.Name)
		if name != opencodeBinName && name != opencodeBinName+".exe" {
			continue
		}
		destPath := filepath.Join(destDir, name)
		outFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			return err
		}
		if _, err := io.Copy(outFile, tr); err != nil {
			outFile.Close()
			return err
		}
		outFile.Close()
		return nil
	}
	return fmt.Errorf("binary not found in archive")
}

func extractBinaryZip(archivePath, destDir string) error {
	zr, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer zr.Close()

	for _, f := range zr.File {
		name := filepath.Base(f.Name)
		if name != opencodeBinName && name != opencodeBinName+".exe" {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		destPath := filepath.Join(destDir, name)
		outFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			rc.Close()
			return err
		}
		if _, err := io.Copy(outFile, rc); err != nil {
			outFile.Close()
			rc.Close()
			return err
		}
		outFile.Close()
		rc.Close()
		return nil
	}
	return fmt.Errorf("binary not found in archive")
}

func DeployConfig(o detect.OS, ollamaEnabled bool) error {
	if err := EnsureDirectories(o); err != nil {
		return err
	}

	configDir := detect.GetOpenCodeConfigDir(o)
	cfgPath := config.GetConfigPath(o)

	cfg, err := config.LoadExistingConfig(cfgPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("load config: %w", err)
		}
		cfg = config.GenerateDefaultConfig()
	}

	if ollamaEnabled {
		config.LinkOllama(cfg)
	}

	normalizeConfigPaths(cfg)

	if err := config.WriteConfig(cfgPath, cfg); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	if err := deployAgents(configDir); err != nil {
		return fmt.Errorf("agents: %w", err)
	}

	if err := deploySkills(configDir); err != nil {
		return fmt.Errorf("skills: %w", err)
	}

	if err := deployAGENTS(configDir); err != nil {
		return fmt.Errorf("AGENTS.md: %w", err)
	}

	fmt.Fprintf(os.Stderr, "Configuration deployed to %s\n", configDir)
	return nil
}

func normalizeConfigPaths(cfg *config.OpenCodeConfig) {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}
	appData := os.Getenv("APPDATA")   // Windows: C:\Users\...\AppData\Roaming
	localAppData := os.Getenv("LOCALAPPDATA") // Windows: C:\Users\...\AppData\Local
	userProfile := os.Getenv("USERPROFILE")   // Windows: C:\Users\...

	// Helper to expand all platform path variables
	expandPath := func(s string) string {
		// Unix-style: ~/path, ${HOME}/path
		if strings.HasPrefix(s, "~/") {
			return filepath.Join(home, s[2:])
		}
		if strings.HasPrefix(s, "${HOME}/") || strings.HasPrefix(s, "${HOME}\\") {
			return filepath.Join(home, s[8:])
		}
		// Windows-style: %USERPROFILE%\path
		if strings.HasPrefix(s, "%USERPROFILE%\\") || strings.HasPrefix(s, "%USERPROFILE%/") {
			return filepath.Join(userProfile, s[14:])
		}
		if strings.HasPrefix(s, "%APPDATA%\\") || strings.HasPrefix(s, "%APPDATA%/") {
			return filepath.Join(appData, s[10:])
		}
		if strings.HasPrefix(s, "%LOCALAPPDATA%\\") || strings.HasPrefix(s, "%LOCALAPPDATA%/") {
			return filepath.Join(localAppData, s[15:])
		}
		return s
	}

	// Expand paths in external_directory keys
	if extDir, ok := cfg.Permission["external_directory"]; ok {
		if edm, ok := extDir.(map[string]any); ok {
			normalized := make(map[string]any, len(edm))
			for k, v := range edm {
				normalized[expandPath(k)] = v
			}
			cfg.Permission["external_directory"] = normalized
		}
	}
	// Expand paths in filesystem MCP command arguments
	if mcp, ok := cfg.MCP["filesystem"]; ok {
		if mcpMap, ok := mcp.(map[string]any); ok {
			if cmd, ok := mcpMap["command"]; ok {
				if args, ok := cmd.([]any); ok {
					expanded := make([]any, 0, len(args))
					for _, arg := range args {
						s, ok := arg.(string)
						if !ok {
							expanded = append(expanded, arg)
							continue
						}
						expanded = append(expanded, expandPath(s))
					}
					mcpMap["command"] = expanded
				}
			}
		}
	}
	// Expand paths in postgres MCP command arguments
	if mcp, ok := cfg.MCP["postgres"]; ok {
		if mcpMap, ok := mcp.(map[string]any); ok {
			if cmd, ok := mcpMap["command"]; ok {
				if args, ok := cmd.([]any); ok {
					expanded := make([]any, 0, len(args))
					for _, arg := range args {
						s, ok := arg.(string)
						if !ok {
							expanded = append(expanded, arg)
							continue
						}
						expanded = append(expanded, expandPath(s))
					}
					mcpMap["command"] = expanded
				}
			}
		}
	}
}

func deployAgents(configDir string) error {
	agentsDir := filepath.Join(configDir, "agents")
	if err := os.MkdirAll(agentsDir, 0755); err != nil {
		return err
	}

	entries, err := fs.ReadDir(assetsFS, "agents")
	if err != nil {
		return fmt.Errorf("read agents dir: %w", err)
	}

	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".md" {
			continue
		}
		data, err := fs.ReadFile(assetsFS, "agents/"+e.Name())
		if err != nil {
			return fmt.Errorf("read %s: %w", e.Name(), err)
		}
		dst := filepath.Join(agentsDir, e.Name())
		if err := os.WriteFile(dst, data, 0644); err != nil {
			return fmt.Errorf("write %s: %w", e.Name(), err)
		}
	}
	return nil
}

func deploySkills(configDir string) error {
	skillsDst := filepath.Join(configDir, "skills")
	if err := os.MkdirAll(skillsDst, 0755); err != nil {
		return err
	}

	entries, err := fs.ReadDir(assetsFS, "skills")
	if err != nil {
		return fmt.Errorf("read skills dir: %w", err)
	}

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		dstPath := filepath.Join(skillsDst, e.Name())
		if err := copyDirFS("skills/"+e.Name(), dstPath); err != nil {
			return fmt.Errorf("copy skill %s: %w", e.Name(), err)
		}
	}
	return nil
}

func deployAGENTS(configDir string) error {
	data, err := fs.ReadFile(assetsFS, "AGENTS.md")
	if err != nil {
		return err
	}
	dst := filepath.Join(configDir, "AGENTS.md")
	return os.WriteFile(dst, data, 0644)
}

func copyDirFS(src, dst string) error {
	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}
	entries, err := fs.ReadDir(assetsFS, src)
	if err != nil {
		return err
	}
	for _, e := range entries {
		srcPath := src + "/" + e.Name()
		dstPath := filepath.Join(dst, e.Name())
		if e.IsDir() {
			if err := copyDirFS(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			data, err := fs.ReadFile(assetsFS, srcPath)
			if err != nil {
				return err
			}
			if err := os.WriteFile(dstPath, data, 0644); err != nil {
				return err
			}
		}
	}
	return nil
}

type ProgressWriter struct {
	Total      int64
	Current    int64
	StatusLine string
}

func (pw *ProgressWriter) Write(p []byte) (int, error) {
	n := len(p)
	pw.Current += int64(n)
	if pw.Total > 0 {
		pct := int(float64(pw.Current) / float64(pw.Total) * 100)
		if runtime.GOOS != "windows" {
			fmt.Fprintf(os.Stderr, "\r%s: %d%%", pw.StatusLine, pct)
		}
	}
	return n, nil
}

// ─── PostgreSQL ──────────────────────────────────────────────────────────────

func InstallPostgreSQL(o detect.OS) error {
	fmt.Fprintf(os.Stderr, "Installing PostgreSQL...\n")

	switch o {
	case detect.Linux:
		// Detect package manager
		if _, err := exec.LookPath("apt-get"); err == nil {
			cmd := exec.Command("apt-get", "update", "-y")
			cmd.Stdout = os.Stderr
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("apt-get update: %w", err)
			}
			cmd = exec.Command("apt-get", "install", "-y", "postgresql", "postgresql-client")
			cmd.Stdout = os.Stderr
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("install postgresql: %w", err)
			}
			// Start PostgreSQL service
			exec.Command("systemctl", "start", "postgresql").Run()
			exec.Command("systemctl", "enable", "postgresql").Run()
			return nil
		}
		if _, err := exec.LookPath("dnf"); err == nil {
			cmd := exec.Command("dnf", "install", "-y", "postgresql-server", "postgresql-contrib")
			cmd.Stdout = os.Stderr
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("install postgresql: %w", err)
			}
			exec.Command("postgresql-setup", "--initdb").Run()
			exec.Command("systemctl", "start", "postgresql").Run()
			exec.Command("systemctl", "enable", "postgresql").Run()
			return nil
		}
		return fmt.Errorf("unsupported package manager (only apt-get/dnf supported)")

	case detect.Darwin:
		if _, err := exec.LookPath("brew"); err != nil {
			return fmt.Errorf("homebrew not found. Install it first: https://brew.sh")
		}
		cmd := exec.Command("brew", "install", "postgresql@16")
		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("brew install postgresql: %w", err)
		}
		exec.Command("brew", "services", "start", "postgresql@16").Run()
		return nil

	case detect.Windows:
		// Try winget (Windows Package Manager)
		if _, err := exec.LookPath("winget"); err == nil {
			cmd := exec.Command("winget", "install", "--accept-source-agreements", "--accept-package-agreements", "PostgreSQL.PostgreSQL.16")
			cmd.Stdout = os.Stderr
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("winget install postgresql: %w", err)
			}
			fmt.Fprintf(os.Stderr, "PostgreSQL installed via winget. Please restart your terminal.\n")
			return nil
		}
		// Try chocolatey
		if _, err := exec.LookPath("choco"); err == nil {
			cmd := exec.Command("choco", "install", "postgresql16", "-y")
			cmd.Stdout = os.Stderr
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("choco install postgresql: %w", err)
			}
			fmt.Fprintf(os.Stderr, "PostgreSQL installed via chocolatey.\n")
			return nil
		}
		return fmt.Errorf("no package manager found. Install PostgreSQL manually from https://www.postgresql.org/download/windows/")

	default:
		return fmt.Errorf("unsupported OS for PostgreSQL installation")
	}
}

func SetupPostgreSQLUserPassword(user, password string) error {
	// Try to set/reset password for the PostgreSQL user
	// First try with peer auth (local socket, no password needed)
	alterCmd := fmt.Sprintf("ALTER USER %s WITH PASSWORD '%s'", user, strings.ReplaceAll(password, "'", "''"))

	cmd := exec.Command("psql", "-U", user, "-c", alterCmd)
	if err := cmd.Run(); err != nil {
		// Fallback: try as postgres OS user
		cmd = exec.Command("sudo", "-u", "postgres", "psql", "-c", alterCmd)
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}
	return nil
}

// InitPGSchema exécute init.sql sur la base PostgreSQL cible.
func InitPGSchema(conn config.PGConnection) error {
	initSQL, err := fs.ReadFile(assetsFS, "pg-mcp-server/init.sql")
	if err != nil {
		return fmt.Errorf("read init.sql: %w", err)
	}

	tmpFile, err := os.CreateTemp("", "pg-init-*.sql")
	if err != nil {
		return fmt.Errorf("create temp sql: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(initSQL); err != nil {
		tmpFile.Close()
		return fmt.Errorf("write tmp sql: %w", err)
	}
	tmpFile.Close()

	cmd := exec.Command("psql",
		"-h", conn.Host,
		"-p", conn.Port,
		"-U", conn.User,
		"-d", conn.DBName,
		"-f", tmpFile.Name(),
		"-v", "ON_ERROR_STOP=1",
	)
	cmd.Env = append(os.Environ(), "PGPASSWORD="+conn.Password)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("run init.sql: %w (verify PostgreSQL is running and credentials are correct)", err)
	}

	fmt.Fprintf(os.Stderr, "PostgreSQL schema initialized (8 tables, indexes, functions)\n")
	return nil
}

func DeployPGMcpServer(configDir string, conn config.PGConnection) error {
	mcpServersDir := filepath.Join(configDir, "mcp-servers", "pg-mcp-server")
	if err := os.MkdirAll(mcpServersDir, 0755); err != nil {
		return fmt.Errorf("create pg-mcp-server dir: %w", err)
	}

	// Deploy package.json
	pkgData, err := fs.ReadFile(assetsFS, "pg-mcp-server/package.json")
	if err != nil {
		return fmt.Errorf("read pg-mcp-server/package.json: %w", err)
	}
	if err := os.WriteFile(filepath.Join(mcpServersDir, "package.json"), pkgData, 0644); err != nil {
		return fmt.Errorf("write package.json: %w", err)
	}

	// Deploy index.js
	jsData, err := fs.ReadFile(assetsFS, "pg-mcp-server/index.js")
	if err != nil {
		return fmt.Errorf("read pg-mcp-server/index.js: %w", err)
	}
	if err := os.WriteFile(filepath.Join(mcpServersDir, "index.js"), jsData, 0644); err != nil {
		return fmt.Errorf("write index.js: %w", err)
	}

	// Run npm install in the deployed directory
	cmd := exec.Command("npm", "install", "--production", "--no-optional")
	cmd.Dir = mcpServersDir
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: npm install failed in %s: %v\n", mcpServersDir, err)
		fmt.Fprintf(os.Stderr, "Run 'cd %s && npm install' manually.\n", mcpServersDir)
	}

	fmt.Fprintf(os.Stderr, "PostgreSQL MCP server deployed to %s\n", mcpServersDir)
	return nil
}
