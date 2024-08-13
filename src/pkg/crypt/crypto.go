package crypt

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/garry-sharp/Sharder/pkg/crypt/alias"

	"github.com/corvus-ch/shamir"
)

type ShardT struct {
	Alias string
	Id    byte
	Data  []byte
}

// ElevenBitToBytes2 converts a slice of integers representing 11-bit numbers to a byte slice.
// Each 11-bit number is converted to 8 bits and concatenated to form the resulting byte slice.
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
			//fmt.Printf("%s ", concatenatedString[i-7:i+1])
		}
	}
	//fmt.Println()
	// for _, res := range result {
	// 	fmt.Printf("%08b ", res)
	// }
	// fmt.Println()
	return result, nil
}

// GetChecksum calculates the checksum of a byte slice using SHA256 and returns the first byte of the hash.
func GetChecksum(entropy []byte) byte {
	hash := sha256.Sum256(entropy)
	csBits := len(entropy) / 4
	//settings.DebugLog("Checksum", fmt.Sprintf("%0b", hash))
	return hash[0] >> (8 - csBits)
}

// GenerateMnemonic generates a mnemonic phrase of the specified length in bytes.
// It uses a random number generator to generate the bytes.
func GenerateMnemonic(len int, lang string) ([]string, error) {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	res := []string{}

	for i := 0; i < len; i++ {
		randWord := wordMapInverse[lang][r.Intn(2048)]
		if randWord == "" {
			return []string{}, fmt.Errorf("word not found")
		}
		res = append(res, randWord)
	}
	return res, nil
}

// MnemonicFromBytes2 converts a byte slice to a mnemonic phrase using the specified language.
// It uses the BytesToElevenBit2 function to convert the bytes to a slice of integers,
// and then maps each integer to a word in the specified language.
func MnemonicFromBytes2(entropy []byte, lang string) (string, error) {
	words := []string{}
	nums := BytesToElevenBit2(entropy)
	for _, num := range nums {
		words = append(words, wordMapInverse[lang][num])
	}
	return strings.Join(words, " "), nil
}

// BytesToElevenBit2 converts a byte slice to a slice of integers representing 11-bit numbers.
// Each byte is converted to 8 bits and concatenated to form a binary string.
// The binary string is then split into 11-bit segments, which are converted to integers.
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

	// for i, c := range str {
	// 	fmt.Print(string(c))
	// 	if i%11 == 0 && i != 0 {
	// 		fmt.Print(" ")
	// 	}
	// }

	csBits := len(bytes) / 4
	bitcount := 0
	totalLength := len(bytes)*8 + csBits
	cs := GetChecksum(bytes)
	//fmt.Println("Checksum", fmt.Sprintf("%08b", cs))
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

// MnemonicToBytes2 converts a mnemonic phrase to a byte slice using the specified language.
// It uses the parseMnemonic function to split the mnemonic phrase into individual words,
// and then maps each word to an index in the specified language's word map.
// The resulting indices are converted to a byte slice using the ElevenBitToBytes2 function.
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
	//fmt.Println(res)
	return res, nil
}

// Assemble combines the data shards into a single byte slice using Shamir's Secret Sharing algorithm,
// and then converts the byte slice to a mnemonic phrase using the specified language.
func Assemble(shards []ShardT, lang string) (string, error) {
	_s := make(map[byte][]byte)
	for _, shard := range shards {
		_s[shard.Id] = shard.Data
	}
	b, err := shamir.Combine(_s)
	if err != nil {
		return "", err
	} else {
		return MnemonicFromBytes2(b, lang)
	}
}

// Shard splits a mnemonic phrase into data shards using Shamir's Secret Sharing algorithm.
// It converts the mnemonic phrase to a byte slice using the specified language,
// and then splits the byte slice into data shards using the shamir.Split function.
// Each data shard is assigned a unique ID and an alias generated by the babble package.
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
		for i, shard := range ans {
			name := alias.GetAlias(i, shard)
			shards = append(shards, ShardT{Alias: name, Id: i, Data: shard})
		}
		return shards, nil
	}
}
