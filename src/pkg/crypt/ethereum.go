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
	seed := pbkdf2.Key([]byte(mnemonic), []byte("mnemonic"), 2048, 64, sha512.New)
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return "", err
	}

	// Derive the private key from the master key using path m/44'/60'/0'/0/0
	childKey, err := masterKey.Derive(44 + hdkeychain.HardenedKeyStart)
	if err != nil {
		return "", err
	}
	childKey, err = childKey.Derive(60 + hdkeychain.HardenedKeyStart)
	if err != nil {
		return "", err
	}
	childKey, err = childKey.Derive(0 + hdkeychain.HardenedKeyStart)
	if err != nil {
		return "", err
	}
	childKey, err = childKey.Derive(0)
	if err != nil {
		return "", err
	}
	childKey, err = childKey.Derive(0)
	if err != nil {
		return "", err
	}

	privKeyBytes, err := childKey.ECPrivKey()
	if err != nil {
		return "", err
	}

	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(privKeyBytes.PubKey().SerializeUncompressed()[1:])
	hash := hasher.Sum(nil)
	return "0x" + hex.EncodeToString(hash[12:]), nil

}
