package cmd

import (
	"fmt"

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

	if dryRun {
		ui.Logger.Info("Dry-run mode - no changes written")
		return nil
	}

	if err := installer.DeployConfig(o, ollama); err != nil {
		return err
	}

	ui.Logger.Info("Setup complete")
	fmt.Println(ui.SuccessStyle.Render("✓ OpenCode configured successfully"))
	return nil
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
