package crypt

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"xxx/settings"
)

// Map of language + word to int.
var wordMap map[string]map[string]int
var wordMapInverse map[string][]string

func LoadWordList(dir string) error {
	wordMap = make(map[string]map[string]int)
	wordMapInverse = make(map[string][]string)

	langRegexp, _ := regexp.Compile("([a-z]{2}).txt")
	files, err := os.ReadDir(dir)

	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() {
			fileName := file.Name()
			if matches := langRegexp.FindStringSubmatch(fileName); len(matches) > 1 {
				lang := matches[1]
				wordMap[lang] = make(map[string]int)
				words, _ := os.ReadFile(fmt.Sprintf("%s/%s", dir, fileName))
				wordsSplit := strings.Split(string(words), "\n")
				wordMapInverse[lang] = wordsSplit
				for i, word := range wordsSplit {
					wordMap[lang][word] = i
				}
				settings.VerboseLog("Loaded wordlist for", lang)
			}
		}
	}
	return nil
}

func GetWordIndex(lang, word string) (int, error) {

	if _, ok := wordMap[lang]; !ok {
		return 0, fmt.Errorf("language %s not supported", lang)
	}

	if _, ok := wordMap[lang][word]; !ok {
		return 0, fmt.Errorf("word %s not found in language %s", word, lang)
	}

	return wordMap[lang][word], nil
}

func GetSupportedLanguages() []string {
	var langs []string
	for lang := range wordMap {
		langs = append(langs, lang)
	}
	return langs
}
