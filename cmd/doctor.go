package cmd

import (
	"fmt"

	"github.com/lfa-cli/lfa-cli-ai/internal/detect"
	"github.com/lfa-cli/lfa-cli-ai/internal/ui"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Run system diagnostics",
	Long:  "Checks your system for OpenCode, Ollama, and configuration paths.",
	RunE: func(cmd *cobra.Command, args []string) error {
		o := detect.DetectOS()
		ocInstalled, ocPath := detect.DetectOpenCode()
		ollamaInstalled, ollamaReachable := detect.DetectOllama()
		configDir := detect.GetOpenCodeConfigDir(o)

		check := ui.AccentStyle.Render("●")
		ok := ui.SuccessStyle.Render("✓")
		warn := ui.WarningStyle.Render("~")
		fail := ui.ErrorStyle.Render("✗")

		type checkRow struct {
			check string
			label string
			value string
		}
		var rows []checkRow

		rows = append(rows, checkRow{check, "OS", o.String()})

		if ocInstalled {
			rows = append(rows, checkRow{ok, "OpenCode", ocPath})
		} else {
			rows = append(rows, checkRow{fail, "OpenCode", "not installed"})
		}

		if ollamaInstalled {
			rows = append(rows, checkRow{ok, "Ollama", "installed"})
		} else {
			rows = append(rows, checkRow{warn, "Ollama", "not found"})
		}

		if ollamaReachable {
			rows = append(rows, checkRow{ok, "Ollama API", "reachable"})
		} else {
			rows = append(rows, checkRow{fail, "Ollama API", "unreachable"})
		}

		rows = append(rows, checkRow{check, "Config Dir", configDir})

		ui.PrintLogo()
		fmt.Println(ui.TitleStyle.Render(" System Diagnostics"))
		fmt.Println()

		for _, r := range rows {
			fmt.Printf("  %s  %-14s %s\n", r.check, ui.AccentStyle.Render(r.label+":"), r.value)
		}

		fmt.Println()
		fmt.Println(ui.DimStyle.Render("  Use 'lfa setup' to deploy configuration"))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
