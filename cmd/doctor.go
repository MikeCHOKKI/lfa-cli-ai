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

		var rows []struct {
			check string
			label string
			value string
		}

		rows = append(rows, struct {
			check string
			label string
			value string
		}{check, "OS", o.String()})

		if ocInstalled {
			rows = append(rows, struct {
				check string
				label string
				value string
			}{ok, "OpenCode", ocPath})
		} else {
			rows = append(rows, struct {
				check string
				label string
				value string
			}{fail, "OpenCode", "not installed"})
		}

		if ollamaInstalled {
			rows = append(rows, struct {
				check string
				label string
				value string
			}{ok, "Ollama", "installed"})
		} else {
			rows = append(rows, struct {
				check string
				label string
				value string
			}{warn, "Ollama", "not found"})
		}

		if ollamaReachable {
			rows = append(rows, struct {
				check string
				label string
				value string
			}{ok, "Ollama API", "reachable"})
		} else {
			rows = append(rows, struct {
				check string
				label string
				value string
			}{fail, "Ollama API", "unreachable"})
		}

		rows = append(rows, struct {
			check string
			label string
			value string
		}{check, "Config Dir", configDir})

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
