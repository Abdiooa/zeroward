package decryption

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"os"
)

func DecryptKey(encryptedKeyFile string, kekKey []byte) ([]byte, error) {
	// Read the encrypted DEK key from the file
	encryptedKey, err := os.ReadFile(encryptedKeyFile)
	if err != nil {
		return nil, err
	}

	// Use the KEK key to decrypt the DEK key
	decryptedKey, err := DecryptData(encryptedKey, kekKey)
	if err != nil {
		return nil, err
	}

	return decryptedKey, nil
}
func DecryptFile(cipherFile string, dekKey []byte) error {
	ciphertext, err := os.ReadFile(cipherFile)
	if err != nil {
		return fmt.Errorf("error reading ciphertext file: %v", err)
	}

	// Use the KEK key to decrypt the DEK key
	decryptedData, err := DecryptData(ciphertext, dekKey)
	if err != nil {
		return err
	}

	if err := os.Remove(cipherFile); err != nil {
		return err
	}

	// Remove the ".enc" extension from the file name
	decryptedFilePath := cipherFile[:len(cipherFile)-4]

	dstFile, err := os.Create(decryptedFilePath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = dstFile.Write(decryptedData)
	if err != nil {
		return err
	}

	return nil
}

func DecryptData(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("error creating AES cipher: %v", err)
	}

	// Create a GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Extract nonce and actual ciphertext
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext is too short")
	}

	nonce := ciphertext[:nonceSize]
	actualCiphertext := ciphertext[nonceSize:]

	// Decrypt the ciphertext
	plaintext, err := gcm.Open(nil, nonce, actualCiphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
