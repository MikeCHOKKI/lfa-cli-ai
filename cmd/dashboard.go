package cmd

import (
	"github.com/lfa-cli/lfa-cli-ai/internal/detect"
	"github.com/lfa-cli/lfa-cli-ai/internal/installer"
	"github.com/lfa-cli/lfa-cli-ai/internal/ui"
	"github.com/spf13/cobra"
)

var dashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "Launch interactive setup dashboard",
	Long:  "Opens an interactive terminal UI to walk through the setup process.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if yesMode {
			return runNonInteractive()
		}
		return ui.RunDashboard()
	},
}

func runNonInteractive() error {
	o := detect.DetectOS()
	ui.Logger.Info("Running setup (non-interactive)", "os", o.String())

	ocInstalled, ocPath := detect.DetectOpenCode()
	if ocInstalled {
		ui.Logger.Info("OpenCode found", "path", ocPath)
	} else {
		ui.Logger.Info("Installing OpenCode")
		if err := installer.InstallOpenCode(o, installer.DefaultVersion); err != nil {
			return err
		}
	}

	ollamaInstalled, ollamaReachable := detect.DetectOllama()
	enableOllama := ollamaInstalled || ollamaReachable
	if enableOllama {
		ui.Logger.Info("Ollama detected, linking", "api", ollamaReachable)
	}

	if err := installer.DeployConfig(o, "data", enableOllama); err != nil {
		return err
	}

	ui.Logger.Info("Setup complete")
	return nil
}

func init() {
	rootCmd.AddCommand(dashboardCmd)
}
