package crypt

import (
	"encoding/hex"
	"strings"
	"testing"
)

func hexDec(x string) []byte {
	a, _ := hex.DecodeString(x[2:])
	return a
}

func TestShardToJson(t *testing.T) {
	input := ShardT{
		Alias: "oddball-piano",
		Id:    hexDec("0x03")[0],
		Data:  hexDec("0xd3d5fce5fda6d0a4f482eb0fc2aba67b"),
	}

	expected := `{
		"alias": "oddball-piano",
		"id": "0x03",
		"data": "0xd3d5fce5fda6d0a4f482eb0fc2aba67b"
		}`

	resultB, _ := ShardToJson(input)
	result := strings.ReplaceAll(string(resultB), " ", "")
	result = strings.ReplaceAll(result, "\n", "")
	result = strings.ReplaceAll(result, "\t", "")

	expected = strings.ReplaceAll(expected, " ", "")
	expected = strings.ReplaceAll(expected, "\n", "")
	expected = strings.ReplaceAll(expected, "\t", "")

	if result != expected {
		t.Errorf("Expected result: %v, but got: %v", expected, result)
	}

}
