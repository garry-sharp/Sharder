package crypt

import (
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Map of language + word to int.
var wordMap map[string]map[string]int
var wordMapInverse map[string][]string

//go:embed wordlists/*.txt
var wordlists embed.FS

func init() {
	LoadWordLists()
}

func LoadWordLists() error {
	wordMap = make(map[string]map[string]int)
	wordMapInverse = make(map[string][]string)
	r, err := regexp.Compile("([0-9])*_*([a-z]{2}).txt")
	if r == nil {
		return err
	}

	paths := []string{}
	err = fs.WalkDir(wordlists, "wordlists", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			expRes := r.FindStringSubmatch(d.Name())
			if len(expRes) == 3 {
				paths = append(paths, path)
			}
		}
		return nil
	})

	sort.Strings(paths)
	for _, path := range paths {
		fname := filepath.Base(path)
		expRes := r.FindStringSubmatch(fname)
		lang := string(expRes[2])
		wordMap[lang] = make(map[string]int)
		wordMapInverse[lang] = []string{}
		wordBytes, err := wordlists.ReadFile(path)
		if err != nil {
			return err
		}
		words := strings.Split(string(wordBytes), "\n")

		for i, word := range words {
			if word != "" {
				wordMap[lang][word] = i
				wordMapInverse[lang] = append(wordMapInverse[lang], word)
			}
		}
	}
	return err
}

func GetWordIndex(lang, word string) (int, error) {

	if _, ok := wordMap[lang]; !ok {
		return -1, fmt.Errorf("language %s not supported", lang)
	}

	if _, ok := wordMap[lang][word]; !ok {
		return -1, fmt.Errorf("word %s not found in language %s", word, lang)
	}

	return wordMap[lang][word], nil
}

func GetWordList(lang string) ([]string, error) {
	if _, ok := wordMapInverse[lang]; !ok {
		return nil, fmt.Errorf("language %s not supported", lang)
	}
	return wordMapInverse[lang], nil
}

func GetSupportedLanguages() []string {
	langs := []string{}
	for lang, _ := range wordMap {
		langs = append(langs, lang)
	}
	return langs
}
