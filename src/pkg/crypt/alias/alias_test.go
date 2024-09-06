package alias

import (
	"encoding/hex"
	"testing"
)

func TestGetAlias(t *testing.T) {
	id, _ := hex.DecodeString("03")
	data, _ := hex.DecodeString("d3d5fce5fda6d0a4f482eb0fc2aba67b")
	expected := "oddball-piano"
	result := GetAlias(id[0], data)
	if result != expected {
		t.Errorf("Expected result: %v, but got: %v", expected, result)
	}
}
