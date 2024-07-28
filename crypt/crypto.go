package crypt

import (
	"crypto/sha256"
	"strings"
	"xxx/settings"

	"github.com/corvus-ch/shamir"
	"github.com/tjarratt/babble"
)

type ShardT struct {
	Alias string
	Id    byte
	Data  []byte
}

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

			if bitIndex == 8 || (len(result)*8+8 > len(ints)*11 && i == 0) { // If currentByte is filled
				settings.DebugLog("Appending byte", currentByte)
				if bitIndex != 8 {
					currentByte = currentByte << (8 - bitIndex)
				}
				result = append(result, currentByte)
				currentByte = 0
				bitIndex = 0
			}
		}
	}
	return result
}

func GetChecksum(entropy []byte, csBits int) byte {
	hash := sha256.Sum256(entropy)
	// Get the first 4 bits (1 byte) of the hash
	return hash[0] >> csBits
}

func EntropyToInts(entropy []byte) ([]int, error) {
	entropyBits := len(entropy) * 8
	checksumBits := entropyBits / 32 // checksum is 1/32 of entropy

	// Step 1: Compute checksum and append it to entropy
	checksum := GetChecksum(entropy, checksumBits)
	entropyWithChecksum := append(entropy, checksum)

	// Step 2: Collect the indices for the mnemonic
	numWords := (entropyBits + checksumBits) / 11
	wordIndexes := make([]int, numWords)

	// Convert to a slice of bits
	bitSlice := make([]byte, (entropyBits+checksumBits+7)/8)
	copy(bitSlice, entropyWithChecksum)

	for i := 0; i < numWords; i++ {
		startBit := i * 11
		startByte := startBit / 8
		startOffset := startBit % 8

		// Extract the bits for the 11-bit index
		var bitSegment uint16 // Use uint16 to avoid overflow
		for j := 0; j < 11; j++ {
			byteIndex := startByte + (j / 8)
			bitIndex := (startOffset + j) % 8
			if byteIndex < len(bitSlice) {
				if (bitSlice[byteIndex] & (1 << (7 - bitIndex))) != 0 {
					bitSegment |= (1 << (10 - j)) // Set the bit in the segment
				}
			}
		}

		wordIndexes[i] = int(bitSegment) // Convert to int after extracting the bits
	}

	return wordIndexes, nil
}

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

func Demodulate(mnemonic string, lang string) ([]byte, error) {
	_mnemonic := parseMnemonic(mnemonic)

	ints := []int{}

	for _, word := range _mnemonic {
		wordIndex, notFoundError := GetWordIndex(lang, word)
		if notFoundError != nil {
			return []byte{}, notFoundError
		}
		ints = append(ints, wordIndex)
	}

	return ElevenBitToBytes(ints), nil
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
		words := []string{}
		nums := BytesToElevenBit(b)
		for _, num := range nums {
			words = append(words, wordMapInverse[lang][num])
		}
		return strings.Join(words, " "), nil
	}
}

func Shard(mnemonic string, k int, n int, lang string) ([]ShardT, error) {
	b, err := Demodulate(mnemonic, lang)
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
