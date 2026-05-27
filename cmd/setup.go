package cmd

import (
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Deploy OpenCode configuration",
	Long:  "Detects the environment, installs OpenCode if needed, and deploys agents and skills.",
	RunE: func(cmd *cobra.Command, args []string) error {
		ollama, _ := cmd.Flags().GetBool("ollama")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		return runSetup(ollama, dryRun)
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
	setupCmd.Flags().Bool("ollama", true, "Enable Ollama integration")
	setupCmd.Flags().Bool("dry-run", false, "Simulate without writing files")
}
