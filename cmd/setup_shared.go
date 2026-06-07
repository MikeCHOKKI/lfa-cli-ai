package cmd

import (
	"fmt"
	"os"

	"github.com/lfa-cli/lfa-cli-ai/internal/config"
	"github.com/lfa-cli/lfa-cli-ai/internal/detect"
	"github.com/lfa-cli/lfa-cli-ai/internal/installer"
	"github.com/lfa-cli/lfa-cli-ai/internal/ui"
)

func runSetup(ollama bool, dryRun bool) error {
	o := detect.DetectOS()
	ui.Logger.Info("Detected OS", "os", o.String())

	ocInstalled, ocPath := detect.DetectOpenCode()
	if ocInstalled {
		ui.Logger.Info("OpenCode found", "path", ocPath)
	} else {
		ui.Logger.Warn("OpenCode not found")
		proceed, err := promptOrYes("OpenCode is not installed. Install now?", true)
		if err != nil {
			return err
		}
		if !proceed {
			ui.Logger.Info("Setup cancelled")
			return nil
		}
		if !dryRun {
			if err := installer.InstallOpenCode(o, installer.DefaultVersion); err != nil {
				return err
			}
		}
	}

	ollamaInstalled, ollamaReachable := detect.DetectOllama()
	if ollamaInstalled || ollamaReachable {
		ui.Logger.Info("Ollama detected", "api", ollamaReachable)
		if ollamaReachable {
			enableOllama, err := promptOrYes("Link Ollama to OpenCode?", true)
			if err != nil {
				return err
			}
			ollama = enableOllama
		}
	}

	// PostgreSQL setup
	pgConn, pgErr := setupPostgreSQL(o, dryRun)
	if pgErr != nil {
		ui.Logger.Warn("PostgreSQL setup skipped", "reason", pgErr.Error())
	}

	if dryRun {
		ui.Logger.Info("Dry-run mode - no changes written")
		return nil
	}

	if err := installer.DeployConfig(o, ollama); err != nil {
		return err
	}

	// Deploy PostgreSQL MCP server if configured
	if pgErr == nil && pgConn.Host != "" {
		configDir := detect.GetOpenCodeConfigDir(o)

		// Deploy pg-mcp-server files
		if err := installer.DeployPGMcpServer(configDir, pgConn); err != nil {
			ui.Logger.Warn("PG MCP server deployment", "error", err)
		}

		// Initialize PostgreSQL schema (tables, indexes, functions)
		ui.Logger.Info("Initializing PostgreSQL schema...")
		if err := installer.InitPGSchema(pgConn); err != nil {
			ui.Logger.Warn("PG schema init", "error", err)
		} else {
			ui.Logger.Info("PostgreSQL schema ready")
		}

		// Link PostgreSQL in config
		cfgPath := config.GetConfigPath(o)
		cfg, err := config.LoadExistingConfig(cfgPath)
		if err == nil {
			config.LinkPostgreSQL(cfg, o, pgConn.Host, pgConn.Port, pgConn.User, pgConn.Password, pgConn.DBName)
			if err := config.WriteConfig(cfgPath, cfg); err != nil {
				ui.Logger.Warn("Write PostgreSQL config", "error", err)
			}
		}
	}

	// ─── Token configuration ──────────────────────────────────────────────
	cfgPath := config.GetConfigPath(o)
	cfg, cfgErr := config.LoadExistingConfig(cfgPath)
	if cfgErr == nil {
		tokensConfigured := false
		for _, t := range config.DefaultTokens {
			if yesMode {
				continue
			}
			// Check if already set
			if val := os.Getenv(t.EnvVar); val != "" {
				config.LinkToken(cfg, "github", "GITHUB_PERSONAL_ACCESS_TOKEN", val)
				ui.Logger.Info(t.Name, "from_env", t.EnvVar)
				tokensConfigured = true
				continue
			}
			proceed, err := promptOrYes(fmt.Sprintf("Configure %s? (%s)", t.Name, t.Description), false)
			if err != nil || !proceed {
				continue
			}
			val, _ := ui.InputPrompt(
				fmt.Sprintf("Enter your %s", t.Name),
				t.Description,
				"",
			)
			if val != "" {
				config.LinkToken(cfg, "github", "GITHUB_PERSONAL_ACCESS_TOKEN", val)
				// Also export for current session
				os.Setenv(t.EnvVar, val)
				tokensConfigured = true
				ui.Logger.Info(t.Name, "status", "configured")
			}
		}
		if tokensConfigured {
			if err := config.WriteConfig(cfgPath, cfg); err != nil {
				ui.Logger.Warn("Write token config", "error", err)
			}
		}
	}

	ui.Logger.Info("Setup complete")
	fmt.Println(ui.SuccessStyle.Render("✓ OpenCode configured successfully"))
	return nil
}

func setupPostgreSQL(o detect.OS, dryRun bool) (config.PGConnection, error) {
	pgStatus := detect.DetectPostgreSQL()
	conn := config.DefaultPGConnection()

	if !pgStatus.Installed {
		ui.Logger.Warn("PostgreSQL not found")
		if !detect.IsPostgreSQLSupported(o) {
			return conn, fmt.Errorf("PostgreSQL installation not supported on this OS")
		}
		proceed, err := promptOrYes("PostgreSQL is not installed. Install now?", true)
		if err != nil {
			return conn, err
		}
		if !proceed {
			return conn, fmt.Errorf("PostgreSQL installation declined")
		}
		if !dryRun {
			if err := installer.InstallPostgreSQL(o); err != nil {
				return conn, fmt.Errorf("PostgreSQL installation failed: %w", err)
			}
			ui.Logger.Info("PostgreSQL installed successfully")
			// Re-detect after install
			newStatus := detect.DetectPostgreSQL()
			if newStatus.Installed {
				pgStatus = newStatus
			}
		}
	} else {
		if pgStatus.Running {
			ui.Logger.Info("PostgreSQL detected", "version", pgStatus.Version)
		} else {
			ui.Logger.Warn("PostgreSQL installed but not running")
		}
	}

	// Ask user for connection details (skip if non-interactive / yes mode)
	if yesMode {
		return conn, nil
	}

	host, _ := ui.InputPrompt("PostgreSQL host", "localhost", conn.Host)
	port, _ := ui.InputPrompt("PostgreSQL port", "5432", conn.Port)
	user, _ := ui.InputPrompt("PostgreSQL user", "postgres", conn.User)
	password, _ := ui.InputPrompt("PostgreSQL password", "postgres", conn.Password)
	dbname, _ := ui.InputPrompt("PostgreSQL database name", "opencode_db", conn.DBName)

	conn = config.PGConnection{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DBName:   dbname,
	}

	// Verify connectivity if not dry-run
	if !dryRun {
		ui.Logger.Info("Verifying PostgreSQL connection...")
		if detect.DetectPostgreSQLConnectivity(host, port, user, password, dbname) {
			ui.Logger.Info("PostgreSQL connection OK")
		} else {
			ui.Logger.Warn("Could not connect to PostgreSQL with provided credentials")
			proceed, err := promptOrYes("Continue anyway?", false)
			if err != nil || !proceed {
				return config.PGConnection{}, fmt.Errorf("PostgreSQL configuration cancelled")
			}
		}
	}

	return conn, nil
}

func promptOrYes(question string, defaultYes bool) (bool, error) {
	if yesMode {
		return defaultYes, nil
	}
	result, err := ui.ConfirmPrompt(question)
	if err != nil {
		ui.Logger.Error("prompt error", "error", err)
		return defaultYes, nil
	}
	return result, nil
}
