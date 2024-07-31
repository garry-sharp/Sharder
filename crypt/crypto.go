package crypt

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"strings"
	"time"
	"xxx/settings"

	"github.com/corvus-ch/shamir"
	"github.com/tjarratt/babble"
)

type ShardT struct {
	Alias string
	Id    byte
	Data  []byte
}

func ElevenBitToBytes2(ints []int) ([]byte, error) {
	concatenatedString := ""
	for _, num := range ints {
		concatenatedString += fmt.Sprintf("%011b", num)
	}
	if len(concatenatedString) != len(ints)*11 {
		return []byte{}, fmt.Errorf("Failed to convert to bytes, length mismatch")
	}

	result := []byte{}
	for i := 0; i < len(concatenatedString); i++ {
		if (i+1)%8 == 0 {
			b := byte(0)
			for j, r := range concatenatedString[i-7 : i+1] {
				if r == '1' {
					b = b | 1<<(7-j)
				}
			}
			result = append(result, b)
			fmt.Printf("%s ", concatenatedString[i-7:i+1])
		}
	}
	fmt.Println()
	for _, res := range result {
		fmt.Printf("%08b ", res)
	}
	fmt.Println()
	return result, nil
}

/*
func ElevenBitToBytes(ints []int) []byte {
	var result []byte
	var currentByte byte
	bitIndex := 0

	for _, num := range ints {
		for i := 10; i >= 0; i-- { // Iterate through each bit of the 11-bit number
			bit := ((num) >> i) & 1
			currentByte = (currentByte << 1) | byte(bit)
			settings.VerboseLog("Number", num, "Bitindex", i, "Bit", bit, "ByteIndex", currentByte)
			bitIndex++

			//|| (len(result)*8+8 > len(ints)*11 && i == 0)
			if bitIndex == 8 { // If currentByte is filled
				settings.DebugLog("Appending byte ", fmt.Sprintf("%d (%08b)", currentByte, currentByte))
				// if bitIndex != 8 {
				// 	currentByte = currentByte << (8 - bitIndex)
				// }
				result = append(result, currentByte)
				currentByte = 0
				bitIndex = 0
			}
		}
	}
	return result
}
*/

func GetChecksum(entropy []byte) byte {
	hash := sha256.Sum256(entropy)
	csBits := len(entropy) / 4
	settings.DebugLog("Checksum", fmt.Sprintf("%0b", hash))
	return hash[0] >> (8 - csBits)
}

/** len in bytes */
func GenerateMnemonic(len int) []string {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	res := []byte{}
	for i := 0; i < len; i++ {
		// Append last byte of uint256
		res = append(res, byte(r.Uint64()&255))
	}
	return []string{}
}

func MnemonicFromBytes2(entropy []byte, lang string) (string, error) {
	words := []string{}
	nums := BytesToElevenBit2(entropy)
	for _, num := range nums {
		words = append(words, wordMapInverse[lang][num])
	}
	return strings.Join(words, " "), nil
}

/*
func MnemonicFromBytes(entropy []byte, lang string) (string, error) {
	words := []string{}
	nums := BytesToElevenBit(entropy)
	for _, num := range nums {
		words = append(words, wordMapInverse[lang][num])
	}
	return strings.Join(words, " "), nil
}*/

/*
func BytesToElevenBit(bytes []byte) []int {
	var result []int
	var currentNum int
	bitIndex := 0

	for _, b := range bytes {
		for i := 7; i >= 0; i-- { // Iterate through each bit of the byte
			bit := (int(b) >> i) & 1
			currentNum = (currentNum << 1) | bit
			bitIndex++

			if bitIndex == 11 || (len(result)*11+11 > len(bytes)*8 && i == 0) { // If currentNum is filled
				result = append(result, currentNum)
				currentNum = 0
				bitIndex = 0
			}
		}
	}
	return result
}
*/

func BytesToElevenBit2(bytes []byte) []int {
	str := ""
	for _, b := range bytes {
		for i := 0; i < 8; i++ {
			if (b>>(7-i))&1 == 1 {
				str += "1"
			} else {
				str += "0"
			}
		}
	}

	for i, c := range str {
		fmt.Print(string(c))
		if i%11 == 0 && i != 0 {
			fmt.Print(" ")
		}
	}

	csBits := len(bytes) / 4
	bitcount := 0
	totalLength := len(bytes)*8 + csBits
	cs := GetChecksum(bytes)
	fmt.Println("Checksum", fmt.Sprintf("%08b", cs))
	res := []int{}
	for bitcount < totalLength {
		b := 0
		for i := 0; i < 11; i++ {
			if bitcount+i >= len(str) {
				b += int(cs)
				break
			} else {
				if str[bitcount+i] == '1' {
					b += 1 << (10 - i)
				}
			}
		}
		res = append(res, b)
		bitcount += 11
	}
	return res
}

func MnemonicToBytes2(mnemonic string, lang string) ([]byte, error) {
	_mnemonic := parseMnemonic(mnemonic)

	ints := []int{}

	for _, word := range _mnemonic {
		wordIndex, notFoundError := GetWordIndex(lang, word)
		if notFoundError != nil {
			return []byte{}, notFoundError
		}
		ints = append(ints, wordIndex)
	}

	res, err := ElevenBitToBytes2(ints)
	if err != nil {
		return []byte{}, err
	}
	fmt.Println(res)
	return res, nil
}

func Assemble(shards []ShardT, lang string) (string, error) {
	_s := make(map[byte][]byte)
	for _, shard := range shards {
		_s[shard.Id] = shard.Data
	}
	b, err := shamir.Combine(_s)
	if err != nil {
		return "", err
	} else {
		return MnemonicFromBytes2(b, settings.GetSettings().Lang)
	}
}

func Shard(mnemonic string, k int, n int, lang string) ([]ShardT, error) {
	b, err := MnemonicToBytes2(mnemonic, lang)
	if err != nil {
		return []ShardT{}, err
	}
	ans, err := shamir.Split(b, n, k)
	if err != nil {
		return []ShardT{}, err
	} else {
		shards := []ShardT{}
		b := babble.NewBabbler()
		for i, shard := range ans {
			shards = append(shards, ShardT{Alias: b.Babble(), Id: i, Data: shard})
		}
		return shards, nil
	}
}
