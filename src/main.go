package main

import (
	"github.com/garry-sharp/Sharder/cmd"
	"github.com/garry-sharp/Sharder/pkg/settings"
)

func main() {
	rootCmd := cmd.SetupCLI()

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		settings.FatalLog(err)
	}

	// f, _ := os.Create("doc")
	// defer f.Close()

	// doc.GenMarkdownTree(rootCmd, "../docs")
	// doc.GenMarkdown(rootCmd, f)

}
