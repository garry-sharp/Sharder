package crypt

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"xxx/settings"
)

type ShardJson struct {
	Alias string `json:"alias"`
	Id    string `json:"id"`
	Data  string `json:"data"`
}

// Map of language + word to int.
var wordMap map[string]map[string]int
var wordMapInverse map[string][]string

func init() {
	fmt.Println(os.Getenv("GOPATH"))
	fmt.Println(os.Getenv("GOROOT"))
	wordMap = make(map[string]map[string]int)
	wordMapInverse = make(map[string][]string)
	wd, err := os.Getwd()
	fmt.Println(wd)
	if err != nil {
		log.Fatal(err)
	}
	filepath := fmt.Sprintf("%s/wordlists", wd)
	files, err := os.ReadDir(filepath)
	if err != nil {
		log.Fatal(err)
	}

	langRegexp, _ := regexp.Compile("([a-z]{2}).txt")
	for _, file := range files {
		if !file.IsDir() {
			fileName := file.Name()
			if matches := langRegexp.FindStringSubmatch(fileName); len(matches) > 1 {
				lang := matches[1]
				wordMap[lang] = make(map[string]int)
				words, _ := os.ReadFile(fmt.Sprintf("%s/%s", filepath, fileName))
				wordsSplit := strings.Split(string(words), "\n")
				wordMapInverse[lang] = wordsSplit
				for i, word := range wordsSplit {
					wordMap[lang][word] = i
				}
				settings.VerboseLog("Loaded wordlist for", lang)
			}
		}
	}
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

func ShardToJson(shard ShardT) ([]byte, error) {

	shardJson := ShardJson{
		Alias: shard.Alias,
		Id:    "0x" + hex.EncodeToString([]byte{shard.Id}),
		Data:  "0x" + hex.EncodeToString(shard.Data),
	}

	return json.MarshalIndent(shardJson, "", "\t")
}

func JSONToShard(jsonData []byte) (ShardT, error) {
	var shardJson ShardJson
	err := json.Unmarshal(jsonData, &shardJson)
	if err != nil {
		return ShardT{}, err
	}

	id, err := hex.DecodeString(shardJson.Id[2:])
	if err != nil {
		return ShardT{}, err
	}

	data, err := hex.DecodeString(shardJson.Data[2:])
	if err != nil {
		return ShardT{}, err
	}

	return ShardT{
		Alias: shardJson.Alias,
		Id:    id[0],
		Data:  data,
	}, nil
}

func parseMnemonic(mnemonic string) []string {
	r, _ := regexp.Compile("([a-zA-Z])*")
	matches := r.FindAll([]byte(mnemonic), -1)
	result := []string{}
	for _, match := range matches {
		if string(match) != "" {
			result = append(result, string(match))
		}
	}
	fmt.Println(result)
	return result
}
