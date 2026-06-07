package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

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
				"command": []string{"npx", "-y", "@modelcontextprotocol/server-memory@0.6.2"},
			},
			"github": map[string]any{
				"type":    "local",
				"command": []string{"npx", "-y", "@modelcontextprotocol/server-github@0.6.2"},
				"environment": map[string]any{
					"GITHUB_PERSONAL_ACCESS_TOKEN": "${GITHUB_TOKEN}",
				},
			},
			"fetch": map[string]any{
				"type":    "local",
				"command": []string{"uvx", "mcp-server-fetch@0.1.4"},
			},
		},
	}
}

func mcpFilesystemCommand(home string) []string {
	dirs := scanExistingProjectDirs(home)
	if len(dirs) == 0 {
		dirs = []string{filepath.Join(home, "Projects")}
	}
	args := []string{"npx", "-y", "@modelcontextprotocol/server-filesystem@0.6.2"}
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
		"command": []string{"npx", "-y", "ollama-mcp"},
	}
}

func LinkPostgreSQL(cfg *OpenCodeConfig, o detect.OS, host, port, user, password, dbname string) {
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", user, password, host, port, dbname)
	// Use the OS-specific config directory for the MCP server path
	pgIndexPath := filepath.Join(detect.GetOpenCodeConfigDir(o), "mcp-servers", "pg-mcp-server", "index.js")
	cfg.MCP["postgres"] = map[string]any{
		"type":    "local",
		"command": []string{"npx", "-y", "tsx", pgIndexPath},
		"environment": map[string]any{
			"DATABASE_URL": connStr,
		},
	}
}

// ─── Tokens ──────────────────────────────────────────────────────────────────

type TokenEntry struct {
	Name        string `json:"name"`
	EnvVar      string `json:"env_var"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
}

var DefaultTokens = []TokenEntry{
	{
		Name:        "GitHub Token",
		EnvVar:      "GITHUB_TOKEN",
		Description: "Personal Access Token for GitHub MCP server (repo scope)",
		Required:    false,
	},
	{
		Name:        "Anthropic API Key",
		EnvVar:      "ANTHROPIC_API_KEY",
		Description: "API key for Anthropic Claude models",
		Required:    false,
	},
}

// LinkToken ajoute ou met à jour une variable d'environnement dans un serveur MCP.
func LinkToken(cfg *OpenCodeConfig, serverKey, envVar, value string) {
	if value == "" {
		return
	}
	mcp, ok := cfg.MCP[serverKey]
	if !ok {
		return
	}
	mcpMap, ok := mcp.(map[string]any)
	if !ok {
		return
	}
	env, ok := mcpMap["environment"].(map[string]any)
	if !ok || env == nil {
		env = make(map[string]any)
		mcpMap["environment"] = env
	}
	// Store the expanded value in the config (resolved from env variable)
	env[envVar] = value
}

// LinkGitHubToken configure le token GitHub dans la config MCP.
func LinkGitHubToken(cfg *OpenCodeConfig, token string) {
	LinkToken(cfg, "github", "GITHUB_PERSONAL_ACCESS_TOKEN", token)
}

type PGConnection struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func DefaultPGConnection() PGConnection {
	return PGConnection{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "postgres",
		DBName:   "opencode_db",
	}
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
