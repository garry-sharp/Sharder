package crypt

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/garry-sharp/Sharder/pkg/settings"
)

// Map of language + word to int.
var wordMap map[string]map[string]int
var wordMapInverse map[string][]string

func LoadWordList(dir string) error {
	var err error
	dir, err = folderParser(dir)
	if err != nil {
		return err
	}
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

func GetWordList(lang string) ([]string, error) {
	if _, ok := wordMapInverse[lang]; !ok {
		return nil, fmt.Errorf("language %s not supported", lang)
	}
	return wordMapInverse[lang], nil
}

func GetSupportedLanguages() []string {
	var langs []string
	for lang := range wordMap {
		langs = append(langs, lang)
	}
	return langs
}

func Download(folder string) error {
	var err error
	var filelocations map[string]string
	filelocations = make(map[string]string)
	filelocations["en"] = "https://raw.githubusercontent.com/bitcoin/bips/master/bip-0039/english.txt"
	filelocations["jp"] = "https://raw.githubusercontent.com/bitcoin/bips/master/bip-0039/japanese.txt"
	filelocations["ko"] = "https://raw.githubusercontent.com/bitcoin/bips/master/bip-0039/korean.txt"
	filelocations["es"] = "https://raw.githubusercontent.com/bitcoin/bips/master/bip-0039/spanish.txt"
	filelocations["zh"] = "https://raw.githubusercontent.com/bitcoin/bips/master/bip-0039/chinese_simplified.txt"
	filelocations["fr"] = "https://raw.githubusercontent.com/bitcoin/bips/master/bip-0039/french.txt"
	filelocations["it"] = "https://raw.githubusercontent.com/bitcoin/bips/master/bip-0039/italian.txt"
	filelocations["cz"] = "https://raw.githubusercontent.com/bitcoin/bips/master/bip-0039/czech.txt"
	filelocations["pt"] = "https://raw.githubusercontent.com/bitcoin/bips/master/bip-0039/portuguese.txt"

	folder, err = folderParser(folder)
	if err != nil {
		return err
	}

	for lang, url := range filelocations {
		settings.VerboseLog("Downloading wordlist for", lang)
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return fmt.Errorf("Failed to download wordlist for %s, error code %d", lang, resp.StatusCode)
		}

		f, err := os.Create(fmt.Sprintf("%s/%s.txt", folder, lang))
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.Copy(f, resp.Body)

		if err != nil {
			return err
		}

		settings.StdLog(fmt.Sprintf("Downloaded wordlist for %s to %s/%s.txt", lang, folder, lang))
	}
	return nil
}

func folderParser(folder string) (string, error) {
	folder = filepath.Clean(folder)
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	locs := strings.Split(folder, "/")
	home := usr.HomeDir
	if locs[0] == "~" || locs[0] == "$HOME" {
		locs[0] = home
	}
	folder = filepath.Join(locs...)
	folder, err = filepath.Abs(folder)
	if err != nil {
		return "", err
	}

	fmt.Println(folder)

	_, err = os.Stat(folder)
	if os.IsNotExist(err) {
		e := os.Mkdir(folder, 0755)
		if e != nil {
			return "", e
		}
	}

	return folder, nil
}
