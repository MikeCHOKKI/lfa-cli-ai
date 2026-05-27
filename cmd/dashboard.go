package cmd

import (
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
	yesMode = true
	return runSetup(true, false)
}

func init() {
	rootCmd.AddCommand(dashboardCmd)
}
