package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "0.2.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("lfa-cli-ai v%s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
