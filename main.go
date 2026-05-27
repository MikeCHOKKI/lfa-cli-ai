package main

import (
	"github.com/lfa-cli/lfa-cli-ai/cmd"
	"github.com/lfa-cli/lfa-cli-ai/internal/installer"
)

func main() {
	installer.SetAssetsFS(getDataFS())
	cmd.Execute()
}
