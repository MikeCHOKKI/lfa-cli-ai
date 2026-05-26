package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ekimart/lfa-cli-ai/internal/detect"
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

func GenerateDefaultConfig() *OpenCodeConfig {
	return &OpenCodeConfig{
		Schema:       "https://opencode.ai/config.json",
		DefaultAgent: "general",
		Skills: map[string]any{
			"paths": []string{"./skills"},
		},
		Permission: map[string]any{
			"read":               "allow",
			"grep":               "allow",
			"glob":               "allow",
			"list":               "allow",
			"webfetch":           "allow",
			"websearch":          "allow",
			"edit":               "ask",
			"bash":               map[string]any{"git *": "allow", "*": "ask"},
			"external_directory": map[string]any{"~/Projets/**": "allow", "*": "ask"},
			"skill":              "ask",
		},
		MCP: map[string]any{},
	}
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
