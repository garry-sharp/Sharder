package settings

import (
	"fmt"
	"log"
	"runtime/debug"
	"time"
)

var settings *Settings
var tzFormat string = "yyyy-MM-dd'T'HH:mm:ssZ"

type Settings struct {
	Verbose     bool
	Debug       bool
	Lang        string
	WordListDir string
}

func GetSettings() *Settings {
	if settings == nil {
		settings = &Settings{
			Verbose:     true,
			Debug:       false,
			Lang:        "en",
			WordListDir: "./wordlists",
		}
	}
	return settings
}

func SetSettings(s *Settings) {
	settings = s
}

func VerboseLog(txt ...any) {
	settings := GetSettings()
	if settings.Verbose {
		fmt.Println(time.Now().Format(tzFormat), " - VERBOSE: ", fmt.Sprint(txt...))
	}
}

func DebugLog(txt ...any) {
	settings := GetSettings()
	if settings.Debug {
		fmt.Println(time.Now().Format(tzFormat), " - DEBUG: ", fmt.Sprint(txt...))
	}
}

func StdLog(txt ...any) {
	fmt.Println(time.Now().Format(tzFormat), " - INFO: ", fmt.Sprint(txt...))
}

func ErrLog(txt ...any) {
	fmt.Println(time.Now().Format(tzFormat), " - ERROR: ", fmt.Sprint(txt...))
}

func FatalLog(txt ...interface{}) {
	debug.PrintStack()
	log.Fatalln(append([]interface{}{time.Now().Format(tzFormat), " - FATAL: "}, txt...)...)
}
