package main

import (
	"os"
	"xxx/cmd"
	"xxx/pkg/crypt"
	"xxx/pkg/settings"

	"github.com/spf13/cobra/doc"
)

func main() {
	rootCmd := cmd.SetupCLI()

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		settings.FatalLog(err)
	}

	f, _ := os.Create("doc")
	defer f.Close()
	doc.GenMarkdown(rootCmd, f)

	crypt.LoadWordList(settings.GetSettings().WordListDir)

}
