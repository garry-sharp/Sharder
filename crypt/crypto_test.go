package crypt

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

// TODO more tests
func TestDemodulate(t *testing.T) {
	// mnemonicA := "ivory maple wage sell gain shop stay praise desk wrist purse road abandon"
	mnemonicB := "jaguar cave"
	lang := "en"

	// expectedResult := []byte{0x01, 0x02, 0x03, 0x04, 0x05}

	expectedResultB := []byte{0b01110111, 0b00100100, 0b10011000}

	result, err := Demodulate(mnemonicB, lang)
	fmt.Println(result)
	fmt.Println(err)

	if !bytes.Equal(result, expectedResultB) {
		t.Errorf("Expected result: %v, but got: %v", expectedResultB, result)
	}
}

func TestBytesToElevenBit(t *testing.T) {
	bytes := []byte{0b01110111, 0b00100100, 0b10011000}
	expectedResult := []int{0b01110111000, 0b10010011000, 0b10011000000}

	result := BytesToElevenBit(bytes)

	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Expected result: %v, but got: %v", expectedResult, result)
	}
}

func TestParseMnemonic(t *testing.T) {
	mnemonic1 := "ivory maple wage sell gain shop stay praise desk wrist purse road abandon"
	mnemonic2 := "ivory maple wage sell gain shop stay praise desk wrist purse road abandon"
	expectedResult := []string{"ivory", "maple", "wage", "sell", "gain", "shop", "stay", "praise", "desk", "wrist", "purse", "road", "abandon"}

	result1 := parseMnemonic(mnemonic1)
	result2 := parseMnemonic(mnemonic2)

	if !reflect.DeepEqual(result1, expectedResult) {
		t.Errorf("Expected result: %v, but got: %v", expectedResult, result1)
	}

	if !reflect.DeepEqual(result2, expectedResult) {
		t.Errorf("Expected result: %v, but got: %v", expectedResult, result2)
	}
}

func TestEntropyToMnemonic(t *testing.T) {
	entropy := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	expectedResult := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}

	result, err := EntropyToInts(entropy)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Expected result: %v, but got: %v", expectedResult, result)
	}
}
