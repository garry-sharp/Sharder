package alias

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"strings"
)

func GetAlias(shardId byte, shardValue []byte) string {
	pre := append([]byte{shardId}, shardValue...)
	hash1 := sha256.Sum256(pre)
	hash2 := sha256.Sum256(hash1[:])
	bigInt1 := new(big.Int).SetBytes(hash1[:])
	bigInt2 := new(big.Int).SetBytes(hash2[:])

	adjectiveId := big.NewInt(0).Mod(bigInt1, big.NewInt(int64(len(adjectives))))
	nounId := big.NewInt(0).Mod(bigInt2, big.NewInt(int64(len(nouns))))
	return fmt.Sprintf("%s-%s", strings.ToLower(adjectives[adjectiveId.Int64()]), strings.ToLower(nouns[nounId.Int64()]))
}
