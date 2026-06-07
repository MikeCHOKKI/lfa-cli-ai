package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/lfa-cli/lfa-cli-ai/internal/ui"
	"github.com/spf13/cobra"
)

var notifyCmd = &cobra.Command{
	Use:   "notify [title] [message]",
	Short: "Send a desktop notification (cross-platform)",
	Long: `Sends a desktop notification using the platform-native mechanism:
  Linux   → notify-send
  macOS   → osascript (AppleScript)
  Windows → PowerShell MessageBox

If the message contains spaces, quote it. If message is omitted, it is read from stdin.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		title := "OpenCode"
		message := ""

		switch len(args) {
		case 2:
			message = args[1]
			title = args[0]
		case 1:
			message = args[0]
		case 0:
			// Read message from stdin
			data, err := io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("read stdin: %w", err)
			}
			message = strings.TrimSpace(string(data))
		}

		if message == "" {
			return fmt.Errorf("message is required (provide as argument or pipe to stdin)")
		}

		if err := ui.Notify(title, message); err != nil {
			return fmt.Errorf("notification failed: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(notifyCmd)
}
