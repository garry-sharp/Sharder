package cmd

import (
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
	err := crypt.LoadWordLists()
	if err != nil {
		settings.FatalLog(err)
	}

	if !verifyLang(lang) {
		log.Fatalln("Unsupported language", lang)
		os.Exit(1)
	}
}
