package crypt

import (
	"strings"
	"testing"
)

func TestMnemonicToEthAddress(t *testing.T) {
	mnemonic := "boy tower radio cradle win toast smile milk task require flush danger"
	lang := "en"
	address, err := MnemonicToEthAddress(mnemonic, lang)
	if err != nil {
		t.Error(err)
	} else {
		expected := strings.ToLower("0x4472E997bB90c6d82a134Cd9b1D94d5F75912EB2")
		if address != expected {
			t.Errorf("Expected %s, got %s", expected, address)
		}
	}

}
