package crypt

import (
	"encoding/hex"
	"encoding/json"
	"regexp"
)

type ShardJson struct {
	Alias string `json:"alias"`
	Id    string `json:"id"`
	Data  string `json:"data"`
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
	return result
}

func VerifyMnemonic(mnemonic string, lang string) bool {
	_mnemonic := parseMnemonic(mnemonic)
	for _, word := range _mnemonic {
		_, notFoundError := GetWordIndex(lang, word)
		if notFoundError != nil {
			return false
		}
	}

	_bytes, _ := MnemonicToBytes2(mnemonic, lang)
	newMnemonic, _ := MnemonicFromBytes2(_bytes, lang)
	if newMnemonic != mnemonic {
		return false
	}
	return true
}
