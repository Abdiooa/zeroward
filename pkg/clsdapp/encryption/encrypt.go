package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"os"
)

func Encrypt(input interface{}, key []byte) error {
	switch v := input.(type) {
	case []byte:
		return EncryptKey(v, key)
	case string:
		return EncryptFile(v, key)
	default:
		return errors.New("unsupported input type")
	}
}

func EncryptKey(dek []byte, kek []byte) error {

	encryptedData, err := EncryptData(dek, kek)
	if err != nil {
		return err
	}

	file, err := os.Create("DEK.key.enc")
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
