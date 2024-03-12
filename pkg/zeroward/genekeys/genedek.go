package genekeys

import (
	"crypto/rand"
	"fmt"
)

const (
	dekSize = 32
)

func GenerateDek() ([]byte, error) {
	dek := make([]byte, dekSize)
	_, err := rand.Read(dek)
	if err != nil {
		return nil, fmt.Errorf("error generating DEK: %v", err)
	}

	return dek, nil
}
