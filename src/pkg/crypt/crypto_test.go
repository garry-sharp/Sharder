package crypt

import (
	"bytes"
	"os"
	"reflect"
	"testing"

	"github.com/garry-sharp/Sharder/pkg/settings"
)

type testpair struct {
	mnemonicParsed []string
	mnemonic       string
	mnemonicindex  []int
	bytestring     []byte
	fullchecksum   []byte
	checksum       byte
}

var tests = []testpair{
	testpair{
		mnemonic:       "vague same appear skull sustain focus rally glory tennis april slice blade",
		mnemonicParsed: []string{"vague", "same", "appear", "skull", "sustain", "focus", "rally", "glory", "tennis", "april", "slice", "blade"},
		mnemonicindex:  []int{1925, 1527, 84, 1622, 1751, 721, 1418, 796, 1785, 87, 1627, 185},
		bytestring:     []byte{0b11110000, 0b10110111, 0b11011100, 0b00101010, 0b01100101, 0b01101101, 0b10101110, 0b10110100, 0b01101100, 0b01010011, 0b00011100, 0b11011111, 0b00100001, 0b01011111, 0b00101101, 0b10001011},
		checksum:       0b00001001,
	},
	testpair{
		mnemonic:       "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about",
		mnemonicParsed: []string{"abandon", "abandon", "abandon", "abandon", "abandon", "abandon", "abandon", "abandon", "abandon", "abandon", "abandon", "about"},
		mnemonicindex:  []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3},
		bytestring:     []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		checksum:       0b00000011,
	},
	testpair{
		bytestring: []byte{0x8c, 0x69, 0x73, 0xcf, 0x2d, 0xf7, 0xa6, 0xc2, 0x11, 0x74, 0xb8, 0x5c, 0xd0, 0x59, 0x18, 0x3d, 0x7f, 0xe1, 0x0c, 0x93, 0x64, 0x63, 0x9b, 0x0a},
		mnemonic:   "midnight entire video fossil kidney genre easy novel fresh lizard ecology kit wrap main eternal midnight only fence",
		checksum:   0b00101000,
	},
}

func TestMain(m *testing.M) {
	// Run setup code here
	settings.SetSettings(&settings.Settings{
		Verbose: false,
		Debug:   true,
	})

	// Run all the tests
	exitCode := m.Run()

	// Run teardown code here

	// Exit with the appropriate exit code
	os.Exit(exitCode)
}

func TestE2E1(t *testing.T) {
	r, _ := MnemonicFromBytes2(tests[2].bytestring, "en")
	if !reflect.DeepEqual(r, tests[2].mnemonic) {
		t.Errorf("Expected result: %v, but got: %v", tests[2].mnemonic, r)
	}

}

func TestBytesToElevenBit2(t *testing.T) {
	result := BytesToElevenBit2([]byte{0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f})
	expectedResult := []int{1019, 2015, 1790, 2039, 1983, 1533, 2031, 1919, 1019, 2015, 1790, 2040}
	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Expected result: %v, but got: %v", expectedResult, result)
	}

	//thank year wave sausage worth useful legal winner thank yellow
}

func TestGetWordIndex(t *testing.T) {
	var res []int
	words := parseMnemonic(tests[0].mnemonic)
	for _, word := range words {
		v, _ := GetWordIndex("en", word)
		res = append(res, v)
	}
	reflect.DeepEqual(res, tests[0].mnemonicindex)
}

// TODO more tests
func TestDemodulate(t *testing.T) {

	result, _ := MnemonicToBytes2(tests[0].mnemonic, "en")

	if !bytes.Equal(result, tests[0].bytestring) {
		t.Errorf("Expected result: %v, but got: %v", tests[0].bytestring, result)
	}
}

// func TestChecksum(t *testing.T) {
// 	result := GetChecksum(tests[0].bytestring)
// 	if result != tests[0].checksum {
// 		t.Errorf("Expected result: %08b, but got: %08b", tests[0].checksum, result)
// 	}
// }

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

func TestElevenBitToBytes2(t *testing.T) {
	ElevenBitToBytes2(tests[0].mnemonicindex)
}
