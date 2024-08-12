package cmd

import (
	"embed"
	"fmt"
	"log"
	"os"

	"github.com/garry-sharp/Sharder/pkg/crypt"
	"github.com/garry-sharp/Sharder/pkg/settings"
	"github.com/savioxavier/termlink"
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

var asciiArt string

//go:embed ascii/ascii.txt
var asciiArtDF embed.FS

func init() {
	f, err := asciiArtDF.ReadFile("ascii/ascii.txt")
	if err != nil {
		settings.FatalLog(err)
	}
	asciiArt = string(f)
}

func generateIntroText() string {
	str := ""
	str += fmt.Sprintln(asciiArt)
	str += fmt.Sprintln("This is a free tool to help shard your mnemonics securely")

	eth := termlink.ColorLink("Ethereum", "ethereum:0x61ae64504549432a94D09E0C258c981698253F7A", "blue")
	btc := termlink.ColorLink("Bitcoin", "bitcoin:bc1qvt37xsc3980zk3nvg44dn92vg2whq73xzsxlna", "blue")
	str += fmt.Sprintf("If you find this tool useful, please consider donating %s or %s . Thank you!\n", eth, btc)
	return str
}
