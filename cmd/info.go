package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/lfa-cli/lfa-cli-ai/internal/detect"
	"github.com/lfa-cli/lfa-cli-ai/internal/ui"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show deployed agents and skills",
	Long:  "Lists deployed agents and skills in the OpenCode config directory.",
	RunE: func(cmd *cobra.Command, args []string) error {
		o := detect.DetectOS()
		configDir := detect.GetOpenCodeConfigDir(o)

		agentsDir := filepath.Join(configDir, "agents")
		skillsDir := filepath.Join(configDir, "skills")

		ui.PrintLogo()
		fmt.Println(ui.TitleStyle.Render(" Deployed Configuration"))
		fmt.Println()

		fmt.Println(ui.AccentStyle.Render(" Agents:"))
		agents, err := os.ReadDir(agentsDir)
		if err != nil {
			fmt.Printf("  %s\n", ui.DimStyle.Render("  (no agents directory)"))
		} else {
			found := false
			for _, e := range agents {
				if !e.IsDir() && filepath.Ext(e.Name()) == ".md" {
					fmt.Printf("  %s %s\n", ui.SuccessStyle.Render("●"), e.Name())
					found = true
				}
			}
			if !found {
				fmt.Printf("  %s\n", ui.DimStyle.Render("  (none deployed)"))
			}
		}
		fmt.Println()

		fmt.Println(ui.AccentStyle.Render(" Skills:"))
		skills, err := os.ReadDir(skillsDir)
		if err != nil {
			fmt.Printf("  %s\n", ui.DimStyle.Render("  (no skills directory)"))
		} else {
			found := false
			for _, e := range skills {
				if e.IsDir() {
					fmt.Printf("  %s %s\n", ui.SuccessStyle.Render("●"), e.Name())
					found = true
				}
			}
			if !found {
				fmt.Printf("  %s\n", ui.DimStyle.Render("  (none deployed)"))
			}
		}
		fmt.Println()

		fmt.Println(ui.DimStyle.Render("  Config directory: " + configDir))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
