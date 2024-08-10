package main

import (
	"os"

	"github.com/garry-sharp/Sharder/cmd"
	"github.com/garry-sharp/Sharder/pkg/crypt"
	"github.com/garry-sharp/Sharder/pkg/settings"

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
	doc.GenMarkdownTree(rootCmd, "../docs")
	doc.GenMarkdown(rootCmd, f)

	crypt.LoadWordList(settings.GetSettings().WordListDir)

}
