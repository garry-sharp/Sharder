package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/garry-sharp/Sharder/pkg/crypt"
	"github.com/garry-sharp/Sharder/pkg/settings"
)

func verifyLang(lang string) bool {
	supportedLangs := crypt.GetSupportedLanguages()
	for _, supportedLang := range supportedLangs {
		if lang == supportedLang {
			return true
		}
	}
	return false
}

func wordListLoadAndVerify() {
	err := crypt.LoadWordList(settings.GetSettings().WordListDir)
	if err != nil {
		settings.FatalLog(fmt.Sprintf("No wordlists found in dir %s", settings.GetSettings().WordListDir))
	}

	if !verifyLang(lang) {
		log.Fatalln("Unsupported language", lang)
		os.Exit(1)
	}
}
