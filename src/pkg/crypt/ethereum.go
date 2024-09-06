package crypt

import (
	"crypto/sha512"
	"encoding/hex"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/sha3"
)

func MnemonicToEthAddress(mnemonic, lang string) (string, error) {
	var err error
	var key *hdkeychain.ExtendedKey
	seed := pbkdf2.Key([]byte(mnemonic), []byte("mnemonic"), 2048, 64, sha512.New)
	key, _ = hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	key, _ = key.Derive(44 + hdkeychain.HardenedKeyStart)
	key, _ = key.Derive(60 + hdkeychain.HardenedKeyStart)
	key, _ = key.Derive(0 + hdkeychain.HardenedKeyStart)
	key, _ = key.Derive(0)
	key, err = key.Derive(0)
	if err != nil {
		return "", err
	}
	privKeyBytes, err := key.ECPrivKey()
	if err != nil {
		return "", err
	}

	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(privKeyBytes.PubKey().SerializeUncompressed()[1:])
	hash := hasher.Sum(nil)
	return "0x" + hex.EncodeToString(hash[12:]), nil

}
