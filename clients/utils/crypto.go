package utils

import (
	"crypto/sha256"
	"fmt"

	"github.com/decred/dcrd/dcrec/secp256k1"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

// Creates a mnemonic from randomness
func CreateMnemonic() (string, error) {
	// Generate a mnemonic for memorization or user-friendly seeds
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return "", err
	}

	mnemonic, _ := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", err
	}

	return mnemonic, nil
}

// Creates a master key from mnemonic and passphrase
func CreateMasterKey(mnemonic, passphrase string) (string, error) {
	seed := bip39.NewSeed(mnemonic, passphrase)

	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return "", err
	}

	return masterKey.B58Serialize(), nil
}

// Checks if the public key matches the private key
func MasterKeysMatch(masterKeySerialized, publicKeySerialized string) (bool, error) {
	masterKey, err := bip32.B58Deserialize(masterKeySerialized)
	if err != nil {
		return false, err
	}

	publicKey, err := bip32.B58Deserialize(publicKeySerialized)
	if err != nil {
		return false, err
	}

	return masterKey.PublicKey() == publicKey, nil
}

func PubKeyMatchesMnemonic(mnemonic, passphrase, pubkeySerialized string) (bool, error) {
	seed := bip39.NewSeed(mnemonic, passphrase)
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return false, err
	}
	masterKeySerialized := masterKey.PublicKey().B58Serialize()

	return masterKeySerialized == pubkeySerialized, nil
}

func SignWithPrivateKey(data, mnemonic, passphrase string) ([]byte, error) {
	seed := bip39.NewSeed(mnemonic, passphrase)
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return []byte{}, err
	}

	// Create the hash of the data
	hash := sha256.New()
	hash.Write([]byte(data))
	hashedData := hash.Sum(nil)

	fmt.Println(hashedData)

	privKey, _ := secp256k1.PrivKeyFromBytes(masterKey.Key)

	// Sign the hash of the data
	signature, err := privKey.Sign(hashedData)
	if err != nil {
		return []byte{}, err
	}
	return signature.Serialize(), nil
}
