package installer

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/ekimart/lfa-cli-ai/internal/config"
	"github.com/ekimart/lfa-cli-ai/internal/detect"
)

const DefaultVersion = "0.1.0"
const opencodeBinName = "opencode"

func EnsureDirectories(o detect.OS) error {
	configDir := detect.GetOpenCodeConfigDir(o)
	dirs := []string{
		configDir,
		filepath.Join(configDir, "agents"),
		filepath.Join(configDir, "skills"),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return fmt.Errorf("create dir %s: %w", d, err)
		}
	}
	return nil
}

func InstallOpenCode(o detect.OS, version string) error {
	arch := detect.GetArch()
	url := detect.GetOpenCodeDownloadURL(o, arch, version)

	fmt.Fprintf(os.Stderr, "Downloading OpenCode from %s\n", url)

	tmpFile, err := os.CreateTemp("", "opencode-*.tar.gz")
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

	if err := extractBinary(tmpFile.Name(), binDir); err != nil {
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

func extractBinary(archivePath, destDir string) error {
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

	tr := tar.NewReader(gzr)
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

func DeployConfig(o detect.OS, dataDir string, ollamaEnabled bool) error {
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

	if err := config.WriteConfig(cfgPath, cfg); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	if err := deployAgents(configDir); err != nil {
		return fmt.Errorf("agents: %w", err)
	}

	if err := deploySkills(dataDir, configDir); err != nil {
		return fmt.Errorf("skills: %w", err)
	}

	if err := deployAGENTS(configDir); err != nil {
		return fmt.Errorf("AGENTS.md: %w", err)
	}

	fmt.Fprintf(os.Stderr, "Configuration deployed to %s\n", configDir)
	return nil
}

func deployAgents(configDir string) error {
	agentsDir := filepath.Join(configDir, "agents")
	if err := os.MkdirAll(agentsDir, 0755); err != nil {
		return err
	}

	entries, err := os.ReadDir(filepath.Join("data", "agents"))
	if err != nil {
		return fmt.Errorf("read agents dir: %w", err)
	}

	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".md" {
			continue
		}
		src := filepath.Join("data", "agents", e.Name())
		dst := filepath.Join(agentsDir, e.Name())
		data, err := os.ReadFile(src)
		if err != nil {
			return fmt.Errorf("read %s: %w", e.Name(), err)
		}
		if err := os.WriteFile(dst, data, 0644); err != nil {
			return fmt.Errorf("write %s: %w", e.Name(), err)
		}
	}
	return nil
}

func deploySkills(dataDir, configDir string) error {
	skillsSrc := filepath.Join(dataDir, "skills")
	skillsDst := filepath.Join(configDir, "skills")
	if err := os.MkdirAll(skillsDst, 0755); err != nil {
		return err
	}

	entries, err := os.ReadDir(skillsSrc)
	if err != nil {
		return fmt.Errorf("read skills dir: %w", err)
	}

	for _, e := range entries {
		srcPath := filepath.Join(skillsSrc, e.Name())
		info, err := os.Stat(srcPath)
		if err != nil || !info.IsDir() {
			continue
		}
		if err := copyDir(srcPath, filepath.Join(skillsDst, e.Name())); err != nil {
			return fmt.Errorf("copy skill %s: %w", e.Name(), err)
		}
	}
	return nil
}

func deployAGENTS(configDir string) error {
	src := filepath.Join("data", "AGENTS.md")
	dst := filepath.Join(configDir, "AGENTS.md")
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}

func copyDir(src, dst string) error {
	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, e := range entries {
		srcPath := filepath.Join(src, e.Name())
		dstPath := filepath.Join(dst, e.Name())
		info, err := os.Stat(srcPath)
		if err != nil {
			return err
		}
		if info.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			data, err := os.ReadFile(srcPath)
			if err != nil {
				return err
			}
			if err := os.WriteFile(dstPath, data, info.Mode()); err != nil {
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
