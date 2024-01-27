package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func EncryptKey(dek []byte, kek []byte, filePath string) error {
	encryptedData, err := EncryptData(dek, kek)
	if err != nil {
		return err
	}

	// Extract directory part of filePath
	outputDEKDir := filepath.Dir(filePath)

	// Create encrypted DEK file path by joining the directory and the new filename
	outputDEKFilePath := filepath.Join(outputDEKDir, "DEK.key.enc")

	file, err := os.Create(outputDEKFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(encryptedData)
	if err != nil {
		return err
	}

	return nil
}

func EncryptFile(filePath string, dek []byte) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return err
	}

	srcFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	plaintext, err := io.ReadAll(srcFile)
	if err != nil {
		return err
	}

	encryptedData, err := EncryptData(plaintext, dek)
	if err != nil {
		return err
	}

	// Remove the plaintext from the file
	if err := os.Remove(filePath); err != nil {
		return err
	}

	outputFilePath := fmt.Sprintf("%s.enc", filePath)

	dstFile, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = dstFile.Write(encryptedData)
	if err != nil {
		return err
	}

	return nil
}
func EncryptData(data []byte, key []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	// Concatenate ciphertext and nonce
	result := ciphertext
	return result, nil
}
