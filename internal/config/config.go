package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/lfa-cli/lfa-cli-ai/internal/detect"
)

type OpenCodeConfig struct {
	Schema       string         `json:"$schema"`
	DefaultAgent string         `json:"default_agent"`
	Skills       map[string]any `json:"skills"`
	Permission   map[string]any `json:"permission"`
	MCP          map[string]any `json:"mcp"`
}

func GetConfigPath(o detect.OS) string {
	return filepath.Join(detect.GetOpenCodeConfigDir(o), "opencode.jsonc")
}

func expandPath(path string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "~"
	}
	if strings.HasPrefix(path, "~/") {
		return filepath.Join(home, path[2:])
	}
	if strings.HasPrefix(path, "${HOME}/") {
		return filepath.Join(home, path[8:])
	}
	return path
}

func GenerateDefaultConfig() *OpenCodeConfig {
	home, _ := os.UserHomeDir()
	user := os.Getenv("USER")
	if user == "" {
		user = os.Getenv("USERNAME")
	}
	if user == "" {
		user = "user"
	}
	if home == "" {
		home = filepath.Join("/home", user)
	}

	return &OpenCodeConfig{
		Schema:       "https://opencode.ai/config.json",
		DefaultAgent: "general",
		Skills: map[string]any{
			"paths": []string{"./skills"},
		},
		Permission: map[string]any{
			"read":      "allow",
			"grep":      "allow",
			"glob":      "allow",
			"list":      "allow",
			"webfetch":  "allow",
			"websearch": "allow",
			"edit":      "ask",
			"bash":      map[string]any{"git *": "allow", "*": "ask"},
			"external_directory": map[string]any{
				filepath.Join(home, "**"): "allow",
				"*":                      "ask",
			},
			"skill": "ask",
		},
		MCP: map[string]any{
			"filesystem": map[string]any{
				"type":    "local",
				"command": mcpFilesystemCommand(home),
			},
			"memory": map[string]any{
				"type":    "local",
				"command": []string{"npx", "-y", "@modelcontextprotocol/server-memory"},
			},
			"github": map[string]any{
				"type":    "local",
				"command": []string{"npx", "-y", "@modelcontextprotocol/server-github"},
				"environment": map[string]any{
					"GITHUB_PERSONAL_ACCESS_TOKEN": "${GITHUB_TOKEN}",
				},
			},
			"fetch": map[string]any{
				"type":    "local",
				"command": []string{"uvx", "mcp-server-fetch"},
			},
		},
	}
}

func mcpFilesystemCommand(home string) []string {
	dirs := scanExistingProjectDirs(home)
	if len(dirs) == 0 {
		dirs = []string{filepath.Join(home, "Projects")}
	}
	args := []string{"npx", "-y", "@modelcontextprotocol/server-filesystem"}
	args = append(args, dirs...)
	return args
}

func scanExistingProjectDirs(home string) []string {
	candidates := []string{
		"Projects", "projects",
		"dev", "Dev",
		"code", "Code",
		"workspace", "Workspace", "Work", "work",
	}
	var found []string
	for _, d := range candidates {
		path := filepath.Join(home, d)
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			found = append(found, path)
		}
	}
	return found
}

func LoadExistingConfig(path string) (*OpenCodeConfig, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cleaned := stripJSONCComments(raw)
	var cfg OpenCodeConfig
	if err := json.Unmarshal(cleaned, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	return &cfg, nil
}

func WriteConfig(path string, cfg *OpenCodeConfig) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func LinkOllama(cfg *OpenCodeConfig) {
	cfg.MCP["ollama"] = map[string]any{
		"type":    "local",
		"command": []string{"ollama", "serve"},
	}
}

func InjectAgents(destDir string, agents []string, readFn func(name string) ([]byte, error)) error {
	agentsDir := filepath.Join(destDir, "agents")
	if err := os.MkdirAll(agentsDir, 0755); err != nil {
		return err
	}
	for _, name := range agents {
		data, err := readFn(name)
		if err != nil {
			return fmt.Errorf("read agent %s: %w", name, err)
		}
		dst := filepath.Join(agentsDir, name)
		if !strings.HasSuffix(name, ".md") {
			dst += ".md"
		}
		if err := os.WriteFile(dst, data, 0644); err != nil {
			return fmt.Errorf("write agent %s: %w", name, err)
		}
	}
	return nil
}

func InjectSkills(destDir string, skills []string, copyFn func(name, dest string) error) error {
	for _, name := range skills {
		if err := copyFn(name, destDir); err != nil {
			return fmt.Errorf("copy skill %s: %w", name, err)
		}
	}
	return nil
}

func stripJSONCComments(raw []byte) []byte {
	result := make([]byte, 0, len(raw))
	inString := false
	for i := 0; i < len(raw); i++ {
		if raw[i] == '"' && (i == 0 || raw[i-1] != '\\') {
			inString = !inString
		}
		if !inString && raw[i] == '/' && i+1 < len(raw) {
			if raw[i+1] == '/' {
				i += 2
				for i < len(raw) && raw[i] != '\n' {
					i++
				}
				if i < len(raw) {
					result = append(result, '\n')
				}
				continue
			}
			if raw[i+1] == '*' {
				i += 2
				for i+1 < len(raw) && !(raw[i] == '*' && raw[i+1] == '/') {
					i++
				}
				i += 2
				continue
			}
		}
		result = append(result, raw[i])
	}
	return result
}
