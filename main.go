package main

import (
	"xxx/cmd"
	"xxx/crypt"
	"xxx/settings"
)

func main() {
	rootCmd := cmd.SetupCLI()

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		settings.FatalLog(err)
	}

	crypt.LoadWordList(settings.GetSettings().WordListDir)

}
