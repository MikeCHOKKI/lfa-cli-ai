package cmd

import (
	"os"

	"github.com/ekimart/lfa-cli-ai/internal/ui"
	"github.com/spf13/cobra"
)

var yesMode bool

var rootCmd = &cobra.Command{
	Use:   "lfa",
	Short: "LFA CLI - OpenCode AI configuration tool",
	Long:  "LFA CLI automates the setup and configuration of OpenCode agents and skills.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if yesMode {
			ui.Logger.Info("Running in non-interactive mode")
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		ui.Logger.Error("Fatal", "error", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&yesMode, "yes", "y", false, "Non-interactive mode (answer yes to all prompts)")
}
