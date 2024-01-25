package genekeys

import (
	"crypto/rand"
	"crypto/sha256"

	"golang.org/x/crypto/pbkdf2"
)

const (
	keySize   = 32
	iteration = 100000
	saltSize  = 64
)

func GenerateKek(passphrase string) ([]byte, error) {
	// Generate a random salt
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	// Derive the key using PBKDF2
	key := pbkdf2.Key([]byte(passphrase), salt, iteration, keySize, sha256.New)

	return key, nil
}
