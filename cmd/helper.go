package cmd

import "xxx/crypt"

func verifyLang(lang string) bool {
	supportedLangs := crypt.GetSupportedLanguages()
	for _, supportedLang := range supportedLangs {
		if lang == supportedLang {
			return true
		}
	}
	return false
}
